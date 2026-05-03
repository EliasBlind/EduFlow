package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/grpc/classes"
	"github.com/EliasBlind/EduFlow/internal/journal_service/grpc/grades"
	"github.com/EliasBlind/EduFlow/internal/journal_service/grpc/homeworks"
	statuscode "github.com/EliasBlind/EduFlow/internal/journal_service/grpc/status_code"
	"github.com/EliasBlind/EduFlow/internal/journal_service/grpc/students"
	"github.com/EliasBlind/EduFlow/internal/journal_service/grpc/subjects"
	"github.com/EliasBlind/EduFlow/internal/journal_service/grpc/teachers"
	teachingload "github.com/EliasBlind/EduFlow/internal/journal_service/grpc/teaching_load"
	"github.com/EliasBlind/EduFlow/pkg/interceptors"
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
) *App {

	token := domain.NewToken(cfg.GRPC.SecretKey)

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryAuthInterceptor(token, log)),
	)

	classes.Register(gRPCServer, nil)
	grades.Register(gRPCServer, nil)
	homeworks.Register(gRPCServer, nil)
	statuscode.Register(gRPCServer, nil)
	students.Register(gRPCServer, nil)
	subjects.Register(gRPCServer, nil)
	teachers.Register(gRPCServer, nil)
	teachingload.Register(gRPCServer, nil)

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
