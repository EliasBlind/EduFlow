package teachers

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/mapper"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Teachers interface {
	CreateTeacher(ctx context.Context, teacherName string) (*domain.Teacher, error)

	ListTeachers(ctx context.Context, req *domain.ListTeachersRequest) (*domain.ListTeachersResponse, error)

	UpdateTeacher(ctx context.Context, req *domain.UpdateTeacher) (*domain.Teacher, error)

	DeleteTeacher(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedTeacherServiceServer
	teachers Teachers
}

func Register(gRPC *grpc.Server, teachers Teachers) {
	journalv1.RegisterTeacherServiceServer(
		gRPC,
		&serverAPI{
			teachers: teachers,
		},
	)
}

func (s *serverAPI) CreateTeacher(ctx context.Context, req *journalv1.CreateTeacherRequest) (*journalv1.Teacher, error) {
	teacher, err := s.teachers.CreateTeacher(ctx, req.FullName)
	if err != nil {
		return nil, err
	}
	return mapper.TeacherToProto(teacher), nil
}

func (s *serverAPI) ListTeachers(ctx context.Context, req *journalv1.ListTeachersRequest) (*journalv1.ListTeachersResponse, error) {
	params := mapper.ListTeachersToDomain(req)

	listTeachers, err := s.teachers.ListTeachers(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.ListTeachersToProto(listTeachers), nil
}

func (s *serverAPI) UpdateTeacher(ctx context.Context, req *journalv1.UpdateTeacherRequest) (*journalv1.Teacher, error) {
	params, err := mapper.UpdateTeacherToDomain(req)
	if err != nil {
		return nil, err
	}

	teacher, err := s.teachers.UpdateTeacher(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.TeacherToProto(teacher), nil
}

func (s *serverAPI) DeleteTeacher(ctx context.Context, req *journalv1.DeleteTeacherRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	err = s.teachers.DeleteTeacher(ctx, id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
