package journal

import (
	"context"
	"fmt"
	"log/slog"

	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	"github.com/EliasBlind/EduFlow/internal/sso_service/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	log           *slog.Logger
	cfg           *config.JournalConfig
	conn          *grpc.ClientConn
	clientTeacher journalv1.TeacherServiceClient
	clientStudent journalv1.StudentServiceClient
}

func New(
	log *slog.Logger,
	cfg *config.JournalConfig,
) *Auth {
	host := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to grpc server", slog.Any("err", err))
		panic("Error connect")
	}

	clientTeacher := journalv1.NewTeacherServiceClient(conn)
	clientStudent := journalv1.NewStudentServiceClient(conn)

	return &Auth{
		log:           log,
		cfg:           cfg,
		conn:          conn,
		clientTeacher: clientTeacher,
		clientStudent: clientStudent,
	}
}

func (a *Auth) CreateStudent(
	ctx context.Context,
	jwt string,
	id uuid.UUID,
	name string,
) error {
	md := metadata.Pairs("authorization", "Bearer "+jwt)
	ctx = metadata.NewOutgoingContext(ctx, md)

	strId := id.String()
	_, err := a.clientStudent.CreateStudent(
		ctx,
		&journalv1.CreateStudentRequest{
			Id:       &strId,
			FullName: name,
		},
	)
	if err != nil {
		a.log.Error("journalv1.CreateStudent grpc call failed", slog.Any("err", err))
		return domain.ErrInternal
	}
	return nil
}

func (a *Auth) CreateTeacher(
	ctx context.Context,
	jwt string,
	id uuid.UUID,
	name string,
) error {
	md := metadata.Pairs("authorization", "Bearer "+jwt)
	ctx = metadata.NewOutgoingContext(ctx, md)

	strId := id.String()
	_, err := a.clientTeacher.CreateTeacher(
		ctx,
		&journalv1.CreateTeacherRequest{
			Id:       &strId,
			FullName: name,
		},
	)
	if err != nil {
		a.log.Error("journalv1.CreateTeacher grpc call failed", slog.Any("err", err))
		return domain.ErrInternal
	}
	return nil
}

func (a *Auth) Stop() {
	if err := a.conn.Close(); err != nil {
		a.log.Error("failed to close grpc connection", slog.Any("err", err))
	}
}
