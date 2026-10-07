package app

import (
	"log/slog"

	grpcapp "github.com/EliasBlind/EduFlow/internal/sso_service/app/grpc"
	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	"github.com/EliasBlind/EduFlow/internal/sso_service/mailer"
	postgres "github.com/EliasBlind/EduFlow/internal/sso_service/storage/postgresql"
	"github.com/EliasBlind/EduFlow/internal/sso_service/storage/redis"
	"github.com/EliasBlind/EduFlow/pkg/validator"
)

type App struct {
	log         *slog.Logger
	GRPCService *grpcapp.App
	redis       *redis.Storage
	sql         *postgres.Storage
	mailer      *mailer.Mailer
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	val := validator.New()

	redis, err := redis.New(log, &cfg.Redis)
	if err != nil {
		panic(err)
	}

	sql, err := postgres.New(log, &cfg.PostgresSql)
	if err != nil {
		panic(err)
	}

	mailerClient, err := mailer.New(log, &cfg.Mailtrap)
	if err != nil {
		panic(err)
	}

	gRPCSrv := grpcapp.New(log, cfg, val, redis, sql, mailerClient)

	return &App{
		log:         log,
		GRPCService: gRPCSrv,
		redis:       redis,
		sql:         sql,
		mailer:      mailerClient,
	}
}

func (a *App) Stop() {
	const op = "app.Stop"
	log := a.log.With("op", op)

	log.Info("stopping application")

	a.GRPCService.Stop()
	a.redis.Stop()
	a.sql.Stop()

	log.Info("application stopped successfully")
}
