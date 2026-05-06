package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	log    *slog.Logger
	cfg    *config.RedisConfig
	client *redis.Client
}

func New(
	log *slog.Logger,
	cfg *config.RedisConfig,
) (*Storage, error) {
	const op = "storage.redis.New"
	log = log.With(
		"op", op,
		"host", cfg.Host,
		"port", cfg.Port,
	)

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   cfg.Db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error("failed to connect to redis", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("redis initialized successfully")

	return &Storage{
		log:    log,
		cfg:    cfg,
		client: client,
	}, nil
}

func (s *Storage) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	const op = "storage.redis.Set"
	log := s.log.With(
		"op", op,
		"key", key,
	)

	err := s.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Error("failed to set value in redis", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, key string) ([]byte, error) {
	const op = "storage.redis.Get"
	log := s.log.With(
		"op", op,
		"key", key,
	)

	value, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrCodeNotFound
		}
		log.Error("failed to get value from redis", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return value, nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	const op = "storage.redis.Set"
	log := s.log.With(
		"op", op,
		"key", key,
	)

	err := s.client.Del(ctx, key).Err()
	if err != nil {
		log.Error("failed to delete key from redis", slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Stop() {
	const op = "storage.redis.Stop"
	log := s.log.With("op", op)

	log.Info("closing redis connection pool")

	if err := s.client.Close(); err != nil {
		log.Error("failed to close redis connection", "error", err)
	}
}
