package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/sqlgen"
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

func (s *Storage) Stop() {
	const op = "storage.postgresql.Stop"

	s.log.With("op", op).Info("closing postgres connection pool")

	s.pool.Close()
}
