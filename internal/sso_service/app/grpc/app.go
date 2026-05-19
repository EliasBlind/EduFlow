package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	authGrpc "github.com/EliasBlind/EduFlow/internal/sso_service/grpc/auth"
	"github.com/EliasBlind/EduFlow/internal/sso_service/mailer"
	"github.com/EliasBlind/EduFlow/internal/sso_service/mapper"
	authService "github.com/EliasBlind/EduFlow/internal/sso_service/service/auth"
	"github.com/EliasBlind/EduFlow/internal/sso_service/service/journal"
	postgres "github.com/EliasBlind/EduFlow/internal/sso_service/storage/postgresql"
	"github.com/EliasBlind/EduFlow/internal/sso_service/storage/redis"
	"github.com/EliasBlind/EduFlow/pkg/i18n"
	"github.com/EliasBlind/EduFlow/pkg/interceptors"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log            *slog.Logger
	gRPCServer     *grpc.Server
	journalService *journal.Auth
	port           int
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	val *validator.Validate,
	redis *redis.Storage,
	sql *postgres.Storage,
	mailer *mailer.Mailer,
) *App {

	fds := os.DirFS(".")
	trans, err := i18n.NewTranslator(fds, cfg.Locale.Path, cfg.Locale.DefaultLang)
	if err != nil {
		panic(err)
	}

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptors.ErrorInterceptor(trans, mapper.MapToRPCError),
		),
	)
	reflection.Register(gRPCServer)

	journalService := journal.New(
		log,
		&cfg.Journal,
	)

	authServ := authService.New(
		log,
		&cfg.Token,
		val,
		redis,
		sql,
		mailer,
		journalService,
	)
	authGrpc.Register(gRPCServer, authServ)

	return &App{
		log:            log,
		gRPCServer:     gRPCServer,
		port:           cfg.GRPC.Port,
		journalService: journalService,
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
	app.journalService.Stop()
}
