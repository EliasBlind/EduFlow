package app

import (
	"log/slog"

	grpcapp "github.com/EliasBlind/EduFlow/internal/journal_service/app/grpc"
	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
)

type App struct {
	GRPCService *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	gRPCSrv := grpcapp.New(log, cfg)

	return &App{
		GRPCService: gRPCSrv,
	}
}
