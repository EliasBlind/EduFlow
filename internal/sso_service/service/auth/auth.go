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
	"github.com/akara-io/zxcvbn"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrExists = errors.New("the record already exists")

type PostgresSql interface {
	CreateUser(ctx context.Context, params *domain.User) (uuid.UUID, error)

	GetPersonByLogin(ctx context.Context, login string) (*domain.User, error)

	GetPersonById(ctx context.Context, userId uuid.UUID) (*domain.User, error)

	UserExist(ctx context.Context, login string) (bool, error)

	// Сохранение новой сессии при Login или Refresh
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, appID int, expiresAt time.Time) (*uuid.UUID, error)

	// Получение сессии для проверки при RefreshToken
	GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*domain.RefreshSession, error)

	// Удаление старой сессии (при Refresh Token Rotation или Logout)
	DeleteSessionByTokenID(ctx context.Context, tokenID uuid.UUID) error

	// (Опционально) Удаление всех сессий пользователя (например, при смене пароля)
	DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error
}

var ErrNotFound = errors.New("key not found in cache")

type Redis interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

type Mailer interface {
	SendVerificationCode(ctx context.Context, email, code string) error
}

type pendingUser struct {
	Params *domain.User
	Code   string
}

type Auth struct {
	log    *slog.Logger
	cfg    *config.TokenConfig
	val    *validator.Validate
	redis  Redis
	sql    PostgresSql
	mailer Mailer
}

func New(
	log *slog.Logger,
	cfg *config.TokenConfig,
	val *validator.Validate,
	redis Redis,
	sql PostgresSql,
	mailer Mailer,
) *Auth {
	gob.Register(domain.RegisterRequest{})
	gob.Register(pendingUser{})
	return &Auth{
		log:    log,
		cfg:    cfg,
		val:    val,
		redis:  redis,
		mailer: mailer,
		sql:    sql,
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

	if err := a.validateRequest(params); err != nil {
		return err
	}

	userExist := a.checkUniqueness(ctx, params.Login)

	hash, code, err := a.preparePendingUser(params)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("err", err))
		return status.Error(codes.Internal, "failed to process security data")
	}

	pending := pendingUser{
		Params: &domain.User{
			Email:        params.Email,
			Login:        params.Login,
			PasswordHash: hash,
		},
		Code: code,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(pending); err != nil {
		log.Error("failed to encode pending user", slog.Any("err", err))
		return status.Error(codes.Internal, "internal error")
	}

	if err := <-userExist; err != nil {
		log.Warn("registration blocked", slog.Any("err", err))
		return err
	}

	err = a.redis.Set(ctx, params.Email, buf.Bytes(), a.cfg.VerificationTTL)
	if err != nil {
		log.Error("failed to save to redis", slog.Any("err", err))
		return status.Error(codes.Internal, "storage error")
	}

	go a.sendEmail(params.Email, code)

	log.Info("registration initiated, waiting for confirmation")
	return nil
}

func (a *Auth) VerifyEmail(
	ctx context.Context,
	params *domain.VerifyRequest,
) (*uuid.UUID, error) {
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
		return nil, status.Error(codes.Internal, "internal error during verification")

	}

	if user.Code != params.Code {
		log.Warn("invalid verification code",
			slog.String("email", params.Email),
			slog.String("op", "verify_code"),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid verification code")
	}

	id, err := a.sql.CreateUser(ctx, user.Params)
	if err = a.handleDbError(log, err); err != nil {
		return nil, err
	}

	return &id, nil
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
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	person, err := a.sql.GetPersonByLogin(ctx, params.Login)
	if err != nil {
		log.Warn("failed to get person")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = bcrypt.CompareHashAndPassword(person.PasswordHash, []byte(params.Password))
	if err != nil {
		log.Warn("invalid credentials", slog.Any("err", err))
		return nil, status.Error(codes.Unauthenticated, "invalid login or password")
	}

	accessToken, err := a.generateToken(person)
	if err != nil {
		log.Error("failed to generate access token", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "failed to generate tokens")
	}

	expiresAt := time.Now().Add(a.cfg.RefreshTokenTTL)
	refreshTokenID, err := a.sql.CreateRefreshToken(ctx, person.Id, params.AppId, expiresAt)
	if err != nil {
		log.Error("failed to save refresh session", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "failed to save session")
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
		return false, status.Error(codes.InvalidArgument, "invalid refresh token format")
	}

	session, err := a.sql.GetSessionByTokenID(ctx, tokenId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("refresh token not found, already logged out", slog.String("token_id", tokenId.String()))
			return false, nil 
		}
		log.Error("failed to get session", slog.Any("err", err))
		return false, status.Error(codes.Internal, "internal error")
	}

	err = a.sql.DeleteSessionByTokenID(ctx, tokenId)
	if err != nil {
		log.Error("failed to delete session", slog.String("token_id", tokenId.String()), slog.Any("err", err))
		return false, status.Error(codes.Internal, "failed to logout")
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
		return nil, status.Error(codes.InvalidArgument, "The refresh token id was not found")
	}

	session, err := a.sql.GetSessionByTokenID(ctx, tokenId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("refresh token not found in database", slog.Any("err", err))
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}

		log.Error("failed to get session", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	if params.AppId != session.AppID {
		log.Warn("app_id mismatch during refresh",
			slog.Int("request_app_id", params.AppId),
			slog.Int("session_app_id", session.AppID),
			slog.String("user_id", session.UserID.String()),
			slog.String("token_id", tokenId.String()),
		)
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token for this application")
	}

	person, err := a.sql.GetPersonById(ctx, session.UserID)
	if err != nil {
		log.Warn("user associated with session not found", slog.Any("err", err))
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	accessToken, err := a.generateToken(person)
	if err != nil {
		log.Error("failed to generate access token", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "failed to generate tokens")
	}

	expiresAt := time.Now().Add(a.cfg.RefreshTokenTTL)
	refreshTokenID, err := a.sql.CreateRefreshToken(ctx, person.Id, params.AppId, expiresAt)
	if err != nil {
		log.Error("failed to save refresh session", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "failed to save session")
	}

	err = a.sql.DeleteSessionByTokenID(ctx, tokenId)
	if err != nil {
		log.Error("failed to revoke session",
			slog.String("token_id", tokenId.String()),
			slog.Any("err", err),
		)

		return nil, status.Error(codes.Internal, "failed to terminate session")
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID.String(),
	}, nil
}

func (a *Auth) validateRequest(p *domain.RegisterRequest) error {
	if err := a.val.Struct(p); err != nil {
		return status.Error(codes.InvalidArgument, "invalid form")
	}
	if err := a.validatePasswordStrength(p.Password, p.Login); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
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
		return nil, "", status.Error(codes.Internal, "security error")
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
			ch <- status.Error(codes.AlreadyExists, "user already exists")
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

	if errors.Is(err, ErrNotFound) {
		log.Warn("cache miss")
		return status.Error(codes.FailedPrecondition, "verification code expired or not found")
	}

	log.Error("redis failure", "err", err)
	return status.Error(codes.Internal, "internal services error")
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
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pending, nil
}

func (a *Auth) handleDbError(log *slog.Logger, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrExists) {
		log.Warn("user already exists", "email", "check logs for context")
		return status.Error(codes.AlreadyExists, "user with this email or username already exists")
	}

	log.Error("database failure", "err", err)
	return status.Error(codes.Internal, "internal database error")
}

func (a *Auth) generateToken(person *domain.User) (string, error) {

	claims := domain.UserClaims{
		Id:    person.Id,
		Login: person.Login,
		Role:  *person.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
