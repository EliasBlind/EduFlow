package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"os"

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
	"github.com/EliasBlind/EduFlow/internal/journal_service/mapper"
	classessvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/classes"
	gradessvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/grades"
	homeworkssvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/homeworks"
	statuscodesvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/status_code"
	studentssvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/students"
	subjectssvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/subjects"
	teacherssvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/teachers"
	teachingloadsvc "github.com/EliasBlind/EduFlow/internal/journal_service/service/teaching_load"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/postgresql"
	"github.com/EliasBlind/EduFlow/pkg/i18n"
	"github.com/EliasBlind/EduFlow/pkg/interceptors"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	sql *postgresql.Storage,
) *App {

	fds := os.DirFS(".")
	trans, err := i18n.NewTranslator(fds, cfg.Locale.Path, cfg.Locale.DefaultLang)
	if err != nil {
		panic(err)
	}

	token := domain.NewToken(cfg.GRPC.SecretKey)

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryAuthInterceptor(token, log),
			interceptors.ErrorInterceptor(trans, mapper.MapToRPCError),
		),
	)
	reflection.Register(gRPCServer)

	cls := classessvc.New(
		log,
		val,
		sql,
	)
	classes.Register(gRPCServer, cls)

	grad := gradessvc.New(
		log,
		val,
		sql,
	)
	grades.Register(gRPCServer, grad)

	hmwrk := homeworkssvc.New(
		log,
		val,
		sql,
	)
	homeworks.Register(gRPCServer, hmwrk)

	stscode := statuscodesvc.New(
		log,
		val,
		sql,
	)
	statuscode.Register(gRPCServer, stscode)

	stdnt := studentssvc.New(
		log,
		val,
		sql,
	)
	students.Register(gRPCServer, stdnt)

	sbjct := subjectssvc.New(
		log,
		val,
		sql,
	)
	subjects.Register(gRPCServer, sbjct)

	tchr := teacherssvc.New(
		log,
		val,
		sql,
	)
	teachers.Register(gRPCServer, tchr)

	tl := teachingloadsvc.New(
		log,
		val,
		sql,
	)
	teachingload.Register(gRPCServer, tl)

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
