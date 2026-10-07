package app

import (
	"log/slog"

	grpcapp "github.com/EliasBlind/EduFlow/internal/journal_service/app/grpc"
	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/postgresql"
	"github.com/EliasBlind/EduFlow/pkg/validator"
)

type App struct {
	GRPCService *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	val := validator.New()

	sql, err := postgresql.New(log, &cfg.Postgresql)
	if err != nil {
		panic(err)
	}

	gRPCSrv := grpcapp.New(
		log,
		cfg,
		val,
		sql,
	)

	return &App{
		GRPCService: gRPCSrv,
	}
}

func (a *App) Stop() {
	a.GRPCService.Stop()
}
