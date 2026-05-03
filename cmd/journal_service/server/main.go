package main

import (
	"github.com/EliasBlind/EduFlow/internal/journal_service/app"
	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	envutil "github.com/EliasBlind/EduFlow/pkg/env"
	"github.com/EliasBlind/EduFlow/pkg/logger"
	"github.com/EliasBlind/EduFlow/pkg/validator"
)

func main() {
	cfg := config.MustLoad()
	val := validator.New()

	log := logger.MustLoad(envutil.EnvDev)

	log.Info("start application")

	application := app.New(log, cfg)
	application.GRPCService.MustRun()

	_ = val
}
