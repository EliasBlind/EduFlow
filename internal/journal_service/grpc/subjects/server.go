package subjects

import (
	"context"

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

	ListSubjects(ctx context.Context) ( []domain.Subject, error)
}

type serverAPI struct {
	journalv1.UnimplementedSubjectServiceServer
	subjects Subjects
}

func Register(gRPC *grpc.Server, subjects Subjects) {
	journalv1.RegisterSubjectServiceServer(
		gRPC,
		&serverAPI{
			subjects: subjects,
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

func (s *serverAPI) ListSubjects(
	ctx context.Context,
	req *journalv1.ListSubjectsRequest,
) (*journalv1.ListSubjectsResponse, error) {
	subjects, err := s.subjects.ListSubjects(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.SubjectsToProto(subjects), nil
}
