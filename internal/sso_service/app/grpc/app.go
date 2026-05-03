package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	authGrpc "github.com/EliasBlind/EduFlow/internal/sso_service/grpc/auth"
	"github.com/EliasBlind/EduFlow/internal/sso_service/mailer"
	authService "github.com/EliasBlind/EduFlow/internal/sso_service/service/auth"
	postgres "github.com/EliasBlind/EduFlow/internal/sso_service/storage/postgresql"
	"github.com/EliasBlind/EduFlow/internal/sso_service/storage/redis"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	val *validator.Validate,
	redis *redis.Storage,
	sql *postgres.Storage,
	mailer *mailer.Mailer,
) *App {

	gRPCServer := grpc.NewServer()

	authServ := authService.New(
		log,
		&cfg.Token,
		val,
		redis,
		sql,
		mailer,
	)
	authGrpc.Register(gRPCServer, authServ)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.GRPC.Port,
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.port),
	)

	l, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", app.port),
	)

	if err != nil {
		log.Error("%s %w", op, err)
		return nil
	}

	log.Info("GRPC server is running",
		slog.String("addr", l.Addr().String()),
	)

	if err := app.gRPCServer.Serve(l); err != nil {
		log.Error("%s %w", op, err)
		return err
	}
	return nil
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("Stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()
}
