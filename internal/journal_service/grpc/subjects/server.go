package subjects

import (
	"context"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/mapper"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Subjects interface {
	CreateSubject(ctx context.Context, subjectName string) (*domain.Subject, error)

	UpdateSubject(ctx context.Context, req *domain.UpdateSubject) (*domain.Subject, error)

	DeleteSubject(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedHomeworkServiceServer
	subjects Subjects
	log      *slog.Logger
}

func Register(gRPC *grpc.Server, logger *slog.Logger, subjects Subjects) {
	journalv1.RegisterHomeworkServiceServer(
		gRPC,
		&serverAPI{
			subjects: subjects,
			log:      logger,
		},
	)
}

func (s *serverAPI) CreateSubject(ctx context.Context, req *journalv1.CreateSubjectRequest) (*journalv1.Subject, error) {
	subject, err := s.subjects.CreateSubject(ctx, req.FullName)
	if err != nil {
		return nil, err
	}

	return mapper.SubjectToProto(subject), nil
}

func (s *serverAPI) UpdateSubject(ctx context.Context, req *journalv1.UpdateSubjectRequest) (*journalv1.Subject, error) {
	params, err := mapper.UpdateSubjectToDomain(req)
	if err != nil {
		return nil, err
	}

	subject, err := s.subjects.UpdateSubject(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.SubjectToProto(subject), nil
}

func (s *serverAPI) DeleteSubject(ctx context.Context, req *journalv1.DeleteSubjectRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.subjects.DeleteSubject(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
