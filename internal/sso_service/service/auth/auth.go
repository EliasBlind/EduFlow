package auth

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/akara-io/zxcvbn"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type PostgresSql interface {
	CreateUser(ctx context.Context, params *domain.User) (uuid.UUID, error)

	GetPersonByLogin(ctx context.Context, login string) (*domain.User, error)

	GetPersonById(ctx context.Context, userId uuid.UUID) (*domain.User, error)

	UserExist(ctx context.Context, login string) (bool, error)

	// Сохранение новой сессии при Login или Refresh
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, appID int, expiresAt time.Time) (*uuid.UUID, error)

	// Получение сессии для проверки при RefreshToken
	GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*domain.RefreshSession, error)

	ListUsers(ctx context.Context) ([]domain.User, error)

	// Удаление старой сессии (при Refresh Token Rotation или Logout)
	DeleteSessionByTokenID(ctx context.Context, tokenID uuid.UUID) error

	// (Опционально) Удаление всех сессий пользователя (например, при смене пароля)
	DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error

	SetRole(ctx context.Context, user *domain.User) error
}

type Redis interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

type Mailer interface {
	SendVerificationCode(ctx context.Context, email, code string) error
}

type JournalService interface {
	CreateStudent(
		ctx context.Context,
		jwt string,
		id uuid.UUID,
		name string,
	) error

	CreateTeacher(
		ctx context.Context,
		jwt string,
		id uuid.UUID,
		name string,
	) error
}

type pendingUser struct {
	Params *domain.User
	AppID  int
	Code   string
}

type Auth struct {
	log     *slog.Logger
	cfg     *config.TokenConfig
	val     *validator.Validate
	redis   Redis
	sql     PostgresSql
	mailer  Mailer
	journal JournalService
}

func New(
	log *slog.Logger,
	cfg *config.TokenConfig,
	val *validator.Validate,
	redis Redis,
	sql PostgresSql,
	mailer Mailer,
	journal JournalService,
) *Auth {
	gob.Register(domain.RegisterRequest{})
	gob.Register(pendingUser{})
	return &Auth{
		log:     log,
		cfg:     cfg,
		val:     val,
		redis:   redis,
		mailer:  mailer,
		sql:     sql,
		journal: journal,
	}
}

func (a *Auth) Register(
	ctx context.Context,
	params *domain.RegisterRequest,
) error {
	const op = "auth.Register"
	log := a.log.With(
		"op", op,
		"email", params.Email,
	)

	if err := a.val.Struct(params); err != nil {
		return domain.ErrInvalidData
	}

	if err := a.validatePasswordStrength(params.Password, params.Login); err != nil {
		return domain.ErrWeakPassword
	}

	userExist := a.checkUniqueness(ctx, params.Login)

	hash, code, err := a.preparePendingUser(params)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("err", err))
		return domain.ErrInternal
	}

	pending := pendingUser{
		Params: &domain.User{
			Email:        params.Email,
			Login:        params.Login,
			PasswordHash: hash,
		},
		AppID: params.AppId,
		Code:  code,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(pending); err != nil {
		log.Error("failed to encode pending user", slog.Any("err", err))
		return domain.ErrInternal
	}

	if err := <-userExist; err != nil {
		log.Warn("registration blocked", slog.Any("err", err))
		return err
	}

	err = a.redis.Set(ctx, params.Email, buf.Bytes(), a.cfg.VerificationTTL)
	if err != nil {
		log.Error("failed to save to redis", slog.Any("err", err))
		return domain.ErrInternal
	}

	go a.sendEmail(params.Email, code)

	log.Info("registration initiated, waiting for confirmation")
	return nil
}

func (a *Auth) VerifyEmail(
	ctx context.Context,
	params *domain.VerifyRequest,
) (*domain.TokenPair, error) {
	const op = "auth.VerifyEmail"
	log := a.log.With(
		"op", op,
		"email", params.Email,
	)

	data, err := a.redis.Get(ctx, params.Email)
	if err = a.handleRedisError(log, err); err != nil {
		return nil, err
	}

	user, err := a.decodePendingUser(data)
	if err != nil {
		return nil, err
	}

	err = a.redis.Del(ctx, params.Email)
	if err != nil {
		log.Error("failed to delete verification code from redis", slog.Any("err", err))
		return nil, domain.ErrInternal

	}

	if user.Code != params.Code {
		log.Warn("invalid verification code",
			slog.String("email", params.Email),
			slog.String("op", "verify_code"),
		)
		return nil, domain.ErrInvalidCode
	}

	user.Params.Id, err = a.sql.CreateUser(ctx, user.Params)

	accessToken, err := a.generateToken(user.Params)
	if err != nil {
		log.Error("failed to generate access token", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	expiresAt := time.Now().Add(a.cfg.RefreshTokenTTL)
	refreshTokenID, err := a.sql.CreateRefreshToken(ctx, user.Params.Id, user.AppID, expiresAt)
	if err != nil {
		log.Error("failed to save refresh session", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID.String(),
	}, nil
}

func (a *Auth) Login(
	ctx context.Context,
	params *domain.LoginRequest,
) (*domain.TokenPair, error) {
	const op = "auth.Login"
	log := a.log.With(
		"op", op,
		"login", params.Login,
	)

	if err := a.val.Struct(params); err != nil {
		log.Warn("invalid request", slog.Any("err", err))
		return nil, domain.ErrInvalidData
	}

	person, err := a.sql.GetPersonByLogin(ctx, params.Login)
	if err != nil {
		log.Warn("failed to get person")
		return nil, domain.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword(person.PasswordHash, []byte(params.Password))
	if err != nil {
		log.Warn("invalid credentials", slog.Any("err", err))
		return nil, domain.ErrInvalidData
	}

	accessToken, err := a.generateToken(person)
	if err != nil {
		log.Error("failed to generate access token", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	expiresAt := time.Now().Add(a.cfg.RefreshTokenTTL)
	refreshTokenID, err := a.sql.CreateRefreshToken(ctx, person.Id, params.AppId, expiresAt)
	if err != nil {
		log.Error("failed to save refresh session", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID.String(),
	}, nil
}

func (a *Auth) Logout(ctx context.Context, refreshToken string) (bool, error) {
	const op = "auth.Logout"
	log := a.log.With("op", op)

	tokenId, err := uuid.Parse(refreshToken)
	if err != nil {
		log.Warn("invalid refresh token format", slog.String("token", refreshToken), slog.Any("err", err))
		return false, domain.ErrInternal
	}

	session, err := a.sql.GetSessionByTokenID(ctx, tokenId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("refresh token not found, already logged out", slog.String("token_id", tokenId.String()))
			return false, nil
		}
		log.Error("failed to get session", slog.Any("err", err))
		return false, domain.ErrInternal
	}

	err = a.sql.DeleteSessionByTokenID(ctx, tokenId)
	if err != nil {
		log.Error("failed to delete session", slog.String("token_id", tokenId.String()), slog.Any("err", err))
		return false, domain.ErrInternal
	}

	log.Info("user logged out", slog.String("user_id", session.UserID.String()), slog.Int("app_id", session.AppID))
	return true, nil
}

func (a *Auth) RefreshToken(ctx context.Context, params *domain.RefreshRequest) (*domain.TokenPair, error) {
	const op = "auth.RefreshToken"
	log := a.log.With(
		"op", op,
		"app_id", params.AppId,
	)

	tokenId, err := uuid.Parse(params.RefreshToken)
	if err != nil {
		log.Error("Invalid uuid received in jwt")
		return nil, domain.ErrUnauthenticated
	}

	session, err := a.sql.GetSessionByTokenID(ctx, tokenId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("refresh token not found in database", slog.Any("err", err))
			return nil, domain.ErrUnauthenticated
		}

		log.Error("failed to get session", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	if params.AppId != session.AppID {
		log.Warn("app_id mismatch during refresh",
			slog.Int("request_app_id", params.AppId),
			slog.Int("session_app_id", session.AppID),
			slog.String("user_id", session.UserID.String()),
			slog.String("token_id", tokenId.String()),
		)
		return nil, domain.ErrTokenAudienceMismatch
	}

	person, err := a.sql.GetPersonById(ctx, session.UserID)
	if err != nil {
		log.Warn("user associated with session not found", slog.Any("err", err))
		return nil, domain.ErrUserNotFound
	}

	accessToken, err := a.generateToken(person)
	if err != nil {
		log.Error("failed to generate access token", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	expiresAt := time.Now().Add(a.cfg.RefreshTokenTTL)
	refreshTokenID, err := a.sql.CreateRefreshToken(ctx, person.Id, params.AppId, expiresAt)
	if err != nil {
		log.Error("failed to save refresh session", slog.Any("err", err))
		return nil, domain.ErrInternal
	}

	err = a.sql.DeleteSessionByTokenID(ctx, tokenId)
	if err != nil {
		log.Error("failed to revoke session",
			slog.String("token_id", tokenId.String()),
			slog.Any("err", err),
		)

		return nil, domain.ErrInternal
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID.String(),
	}, nil
}

func (a *Auth) ListUsers(ctx context.Context, token string) ([]domain.User, error) {
	user, err := a.validateToken(token)
	if err != nil {
		return nil, err
	}

	if user.Role != "admin" {
		return nil, domain.ErrAccessDenied
	}

	res, err := a.sql.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *Auth) SetRole(ctx context.Context, token string, user *domain.User) error {
	const op = "statuscodesvc.SetRole"

	// Первичный логгер
	log := a.log.With(slog.String("op", op))
	log.Info("attempting to set user role", slog.String("target_user_id", user.Id.String()))

	userReq, err := a.validateToken(token)
	if err != nil {
		log.Error("token validation failed", slog.Any("err", err))
		return err
	}

	if userReq.Role != "admin" {
		log.Warn("access denied: user is not an admin", slog.String("admin_role", userReq.Role))
		return domain.ErrAccessDenied
	}

	if user.Role == nil || *user.Role == roles.RoleUnknown {
		log.Warn("invalid data: target user role is unknown or nil")
		return domain.ErrInvalidData
	}

	targetRole := user.Role

	dbUser, err := a.sql.GetPersonById(ctx, user.Id)
	if err != nil {
		log.Error("failed to get person by id from db", slog.Any("err", err))
		return domain.ErrInternal
	}

	user = dbUser
	user.Role = targetRole

	if user.Login == "" {
		log.Warn("user login from DB is empty! Using fallback name to prevent grpc crash")
		return domain.ErrInvalidData
	}

	bgCtx := context.WithoutCancel(ctx)
	go func() {
		log.Info("starting asynchronous profile creation in journal service",
			slog.String("target_user_id", user.Id.String()),
			slog.String("target_user_login", user.Login),
			slog.String("assigned_role", user.Role.String()),
			slog.String("admin_user_id", userReq.ID),
		)

		if user.Role.IsStudent() {
			err := a.journal.CreateStudent(
				bgCtx,
				token,
				user.Id,
				user.Login,
			)
			if err != nil {
				log.Error("failed to async create student in journal", slog.Any("err", err))
			} else {
				log.Info("successfully created student in journal service")
			}
		} else if user.Role.IsTeacher() {
			err := a.journal.CreateTeacher(
				bgCtx,
				token,
				user.Id,
				user.Login,
			)
			if err != nil {
				log.Error("failed to async create teacher in journal", slog.Any("err", err))
			} else {
				log.Info("successfully created teacher in journal service")
			}
		}
	}()

	log.Info("updating user role in database",
		slog.String("target_user_id", user.Id.String()),
		slog.String("assigned_role", user.Role.String()),
	)

	err = a.sql.SetRole(ctx, user)
	if err != nil {
		log.Error("failed to update user role in database", slog.Any("err", err))
		return err
	}

	log.Info("user role successfully updated")
	return nil
}

func (a *Auth) validatePasswordStrength(password, login string) error {
	const op = "auth.validatePasswordStrength"
	log := a.log.With(
		"op", op,
		"login", login,
	)

	result := zxcvbn.PasswordStrength(password, []string{login})

	if result.Score < 3 {
		log.Warn("weak password detected",
			"score", result.Score,
			"warning", result.Feedback.Warning,
			"suggestions", result.Feedback.Suggestions,
		)

		errMsg := "the password is too weak"
		if result.Feedback.Warning != "" {
			errMsg += ": " + result.Feedback.Warning
		}
		if len(result.Feedback.Suggestions) > 0 {
			errMsg += " (advice: " + strings.Join(result.Feedback.Suggestions, ", ") + ")"
		}
		return fmt.Errorf("%s: %w", op, errors.New(errMsg))
	}
	return nil
}

func (a *Auth) preparePendingUser(p *domain.RegisterRequest) ([]byte, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", domain.ErrInternal
	}

	return hash, fmt.Sprintf("%06d", rand.Intn(1000000)), nil
}

func (a *Auth) checkUniqueness(ctx context.Context, login string) <-chan error {
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		exists, err := a.sql.UserExist(ctx, login)
		if err != nil {
			ch <- fmt.Errorf("db check failed: %w", err)
			return
		}
		if exists {
			ch <- domain.ErrUserAlreadyExists
			return
		}
		ch <- nil
	}()

	return ch
}

func (a *Auth) sendEmail(email, code string) error {
	const op = "auth.sendEmail"
	log := a.log.With(
		"op", op,
		"email", email,
	)

	sendCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.mailer.SendVerificationCode(sendCtx, email, code)
	if err != nil {
		log.Error("failed to send email", slog.Any("err", err))
		return err
	}
	return nil
}

func (a *Auth) handleRedisError(log *slog.Logger, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, domain.ErrCodeNotFound) {
		log.Warn("cache miss")
		return domain.ErrCodeNotFound
	}

	log.Error("redis failure", "err", err)
	return domain.ErrInternal
}

func (a *Auth) decodePendingUser(data []byte) (*pendingUser, error) {
	const op = "auth.decodePendingUser"
	log := a.log.With(
		"op", op,
	)

	var pending pendingUser
	reader := bytes.NewReader(data)
	if err := gob.NewDecoder(reader).Decode(&pending); err != nil {
		log.Error("failed to decode pending user", "err", err)
		return nil, domain.ErrInternal
	}
	return &pending, nil
}

func (a *Auth) generateToken(person *domain.User) (string, error) {

	claims := domain.UserClaims{
		Id:    person.Id,
		Login: person.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	if person.Role != nil {
		role := person.Role.String()
		claims.Role = role
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Auth) validateToken(tokenStr string) (*domain.UserClaims, error) {
	var claims domain.UserClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return []byte(a.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthenticated
	}

	return &claims, nil
}
