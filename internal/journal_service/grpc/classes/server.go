package classes

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

type Classes interface {
	CreateClass(ctx context.Context, params *domain.CreateClass) (*domain.Class, error)

	GetClass(ctx context.Context, class_id uuid.UUID) (*domain.Class, error)

	ListTeacherClasses(ctx context.Context, params *domain.ListTeacherClasses) (*domain.ListClasses, error)

	UpdateClass(ctx context.Context, params *domain.UpdateClass) (*domain.Class, error)

	DeleteClassRequest(ctx context.Context, class_id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedClassesServiceServer
	classes Classes
	log     *slog.Logger
}

func Register(gRPC *grpc.Server, logger *slog.Logger, classes Classes) {
	journalv1.RegisterClassesServiceServer(
		gRPC,
		&serverAPI{
			classes: classes,
			log:     logger,
		},
	)
}

func (s *serverAPI) CreateClass(ctx context.Context, req *journalv1.CreateClassRequest) (*journalv1.Class, error) {
	params := mapper.CreateClassToDomain(req)
	class, err := s.classes.CreateClass(ctx, params)

	if err != nil {
		return nil, err
	}
	return mapper.ClassToProto(class), nil
}

func (s *serverAPI) GetClass(ctx context.Context, req *journalv1.GetClassRequest) (*journalv1.Class, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	class, err := s.classes.GetClass(ctx, id)

	if err != nil {
		return nil, err
	}
	return mapper.ClassToProto(class), err
}

func (s *serverAPI) ListTeacherClasses(ctx context.Context, req *journalv1.ListTeacherClassesRequest) (*journalv1.ListClassesResponse, error) {
	param, err := mapper.ListTeacherClassesToDomain(req)
	if err != nil {
		return nil, err
	}

	list_teacher_classes, err := s.classes.ListTeacherClasses(ctx, param)
	if err != nil {
		return nil, err
	}
	return mapper.ListTeacherClassesToProto(list_teacher_classes), err
}

func (s *serverAPI) UpdateClass(ctx context.Context, req *journalv1.UpdateClassRequest) (*journalv1.Class, error) {
	param, err := mapper.UpdateClassToDomain(req)
	if err != nil {
		return nil, err
	}

	class, err := s.classes.UpdateClass(ctx, param)

	if err != nil {
		return nil, err
	}
	return mapper.ClassToProto(class), nil
}

func (s *serverAPI) DeleteClass(ctx context.Context, req *journalv1.DeleteClassRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.classes.DeleteClassRequest(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
