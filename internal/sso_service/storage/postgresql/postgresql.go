package postgresql

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"github.com/EliasBlind/EduFlow/internal/sso_service/storage/sqlgen"
	pkgmapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log     *slog.Logger
	cfg     *config.PostgresqlConfig
	queries *sqlgen.Queries
	pool    *pgxpool.Pool
}

func New(
	log *slog.Logger,
	cfg *config.PostgresqlConfig,
) (*Storage, error) {
	const op = "postgres.New"
	log = log.With(
		"op", op,
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
	)

	dbString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Sslmode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Error("failed to create pgxpool", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to ping postgres", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("postgres connection established successfully")

	return &Storage{
		log:     log,
		queries: sqlgen.New(pool),
		pool:    pool,
	}, nil
}

func (s *Storage) CreateUser(ctx context.Context, params *domain.User) (uuid.UUID, error) {
	const op = "storage.postgresql.CreateUser"
	log := s.log.With(
		"op", op,
		"email", params.Email,
		"username", params.Login,
	)

	arg := sqlgen.CreatePersonParams{
		Email:        params.Email,
		Username:     params.Login,
		PasswordHash: params.PasswordHash,
	}

	person, err := s.queries.CreatePerson(ctx, arg)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" { // Unique violation
				log.Warn("user already exists", "error", err)
				return uuid.Nil, domain.ErrUserAlreadyExists
			}
		}
		log.Error("failed to create person in database", "error", err)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := uuid.FromBytes(person.ID.Bytes[:])
	if err != nil {
		log.Error("failed to parse uuid from database result", "error", err)
		return uuid.Nil, fmt.Errorf("%s: uuid parse error: %w", op, err)
	}

	log.Info("user created successfully", "user_id", res.String())

	return res, nil
}

func (s *Storage) GetPersonByLogin(ctx context.Context, login string) (*domain.User, error) {
	const op = "storage.postgresql.GetPersonByLogin"
	log := s.log.With(
		"op", op,
		"login", login,
	)

	person, err := s.queries.GetPersonByLogin(ctx, login)
	if err != nil {
		log.Error("failed to get person from database", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := uuid.FromBytes(person.ID.Bytes[:])
	if err != nil {
		log.Error("failed to parse uuid from database", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("person retrieved successfully")

	role := roles.GetRole(person.UserRole)
	return &domain.User{
		Id:           uid,
		Email:        person.Email,
		Login:        person.Username,
		PasswordHash: person.PasswordHash,
		Role:         &role,
	}, nil
}

func (s *Storage) GetPersonById(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	const op = "storage.postgresql.GetPersonById"
	log := s.log.With(
		"op", op,
	)

	pgId := pkgmapper.ToPgUUID(uid)
	person, err := s.queries.GetPersonById(ctx, pgId)
	if err != nil {
		log.Error("failed to get person from database", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("person retrieved successfully")

	role := roles.GetRole(person.UserRole)
	return &domain.User{
		Id:    uid,
		Email: person.Email,
		Login: person.Username,
		Role:  &role,
	}, nil
}

func (s *Storage) UserExist(ctx context.Context, login string) (bool, error) {
	const op = "storage.UserExist"
	log := s.log.With(
		"op", op,
		"login", login,
	)
	exists, err := s.queries.UserExist(ctx, login)
	if err != nil {
		log.Error("failed to check user existence",
			slog.String("op", op),
			slog.String("login", login),
			slog.Any("err", err),
		)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (s *Storage) CreateRefreshToken(
	ctx context.Context,
	userID uuid.UUID,
	appID int,
	expiresAt time.Time,
) (*uuid.UUID, error) {
	const op = "storage.postgresql.SaveRefreshToken"

	tokenID := uuid.New()

	hash := sha256.Sum256([]byte(tokenID.String()))
	arg := sqlgen.SaveRefreshTokenParams{
		UserID:    pkgmapper.ToPgUUID(userID),
		AppID:     int32(appID),
		TokenID:   hash[:],
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	}

	_, err := s.queries.SaveRefreshToken(ctx, arg)
	if err != nil {
		s.log.Error("failed to save refresh token", "op", op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &tokenID, nil
}

func (s *Storage) GetSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (*domain.RefreshSession, error) {
	const op = "storage.postgresql.GetSessionByTokenID"
	hash := sha256.Sum256([]byte(tokenID.String()))
	session, err := s.queries.GetSessionByTokenID(ctx, hash[:])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.RefreshSession{
		ID:        uuid.UUID(session.ID.Bytes),
		UserID:    uuid.UUID(session.UserID.Bytes),
		AppID:     int(session.AppID),
		TokenID:   tokenID,
		ExpiresAt: session.ExpiresAt.Time,
	}, nil
}

func (s *Storage) DeleteSessionByTokenID(ctx context.Context, tokenID uuid.UUID) error {
	const op = "storage.postgresql.DeleteSessionByTokenID"
	hash := sha256.Sum256([]byte(tokenID.String()))
	err := s.queries.DeleteSessionByTokenID(ctx, hash[:])
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	const op = "storage.postgresql.DeleteAllUserSessions"

	err := s.queries.DeleteAllUserSessions(ctx, pkgmapper.ToPgUUID(userID))
	if err != nil {
		s.log.Error("failed to delete all user sessions", "op", op, "user_id", userID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.queries.ListUsers(ctx)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return pkgmapper.MapSlice(
		users,
		func(user *sqlgen.ListUsersRow) domain.User {
			role := roles.GetRole(user.UserRole)
			return domain.User{
				Id:    user.ID.Bytes,
				Email: user.Email,
				Login: user.Username,
				Role:  &role,
			}
		},
	), nil
}

func (s *Storage) SetRole(ctx context.Context, user *domain.User) error {
	arg := sqlgen.UpdateRoleParams{
		ID:       pkgmapper.ToPgUUID(user.Id),
		UserRole: *user.Role,
	}

	_, err := s.queries.UpdateRole(ctx, arg)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

func (s *Storage) Stop() {
	const op = "storage.postgresql.Stop"

	s.log.With("op", op).Info("closing postgres connection pool")

	s.pool.Close()
}
