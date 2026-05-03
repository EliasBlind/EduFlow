package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EliasBlind/EduFlow/internal/sso_service/app"
	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	envutil "github.com/EliasBlind/EduFlow/pkg/env"
	"github.com/EliasBlind/EduFlow/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.MustLoad(envutil.EnvDev)

	log.Info("start application")

	application := app.New(log, cfg)
	go application.GRPCService.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	application.Stop()
	log.Info("application stopped")
}
