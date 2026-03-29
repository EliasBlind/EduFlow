package main

import (
	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	"github.com/EliasBlind/EduFlow/internal/journal_service/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.MustLoad(config.EnvDev)

	log.Info("start application")
	
	// TODO: init app

	// TODO: started gRPC server
}
