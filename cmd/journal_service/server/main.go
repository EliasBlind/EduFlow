package main

import (
	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	"github.com/EliasBlind/EduFlow/internal/journal_service/logger"
	"github.com/EliasBlind/EduFlow/internal/journal_service/validator"
	"github.com/EliasBlind/EduFlow/pkg/interceptors"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	val := validator.New()

	log := logger.MustLoad(config.EnvDev)

	log.Info("start application")

	// TODO: init app

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptor(log)),
	)
	_ = cfg
	_ = val
	_ = gRPCServer
}
