package students

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

type Students interface {
	CreateStudent(ctx context.Context, req *domain.CreateStudent) (*domain.Student, error)

	GetStudent(ctx context.Context, studentID uuid.UUID) (*domain.Student, error)

	ListStudents(ctx context.Context, req *domain.ListStudentsRequest) (*domain.ListStudentsResponse, error)

	UpdateStudent(ctx context.Context, req *domain.UpdateStudent) (*domain.Student, error)

	DeleteStudent(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedHomeworkServiceServer
	students Students
	log      *slog.Logger
}

func Register(gRPC *grpc.Server, logger *slog.Logger, students Students) {
	journalv1.RegisterHomeworkServiceServer(
		gRPC,
		&serverAPI{
			students: students,
			log:      logger,
		},
	)
}

func (s *serverAPI) CreateStudent(ctx context.Context, req *journalv1.CreateStudentRequest) (*journalv1.Student, error) {
	params, err := mapper.CreateStudentToDomain(req)
	if err != nil {
		return nil, err
	}

	student, err := s.students.CreateStudent(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.StudentToProto(student), nil
}

func (s *serverAPI) GetStudent(ctx context.Context, req *journalv1.StatusCode) (*journalv1.Student, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	student, err := s.students.GetStudent(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.StudentToProto(student), nil
}

func (s *serverAPI) ListStudents(ctx context.Context, req *journalv1.ListStudentsRequest) (*journalv1.ListStudentsResponse, error) {
	params, err := mapper.ListStudentsToDomain(req)
	if err != nil {
		return nil, err
	}

	listStudent, err := s.students.ListStudents(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.ListStudentsToProto(listStudent), nil
}

func (s *serverAPI) UpdateStudent(ctx context.Context, req *journalv1.UpdateStudentRequest) (*journalv1.Student, error) {
	params, err := mapper.UpdateStudentToDomain(req)
	if err != nil {
		return nil, err
	}

	student, err := s.students.UpdateStudent(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.StudentToProto(student), nil
}

func (s *serverAPI) DeleteStudent(ctx context.Context, req *journalv1.DeleteStudentRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.students.DeleteStudent(ctx, id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
