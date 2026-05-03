package homeworks

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/mapper"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Homeworks interface {
	RecordHomework(ctx context.Context, req *domain.RecordHomework) (*domain.Homework, error)

	UpdateHomework(ctx context.Context, req *domain.UpdateHomework) (*domain.Homework, error)

	ListHomework(ctx context.Context, req *domain.ListHomeworkRequest) (*domain.ListHomeworkResponse, error)

	DeleteHomework(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedHomeworkServiceServer
	homeworks Homeworks
}

func Register(gRPC *grpc.Server, homeworks Homeworks) {
	journalv1.RegisterHomeworkServiceServer(
		gRPC,
		&serverAPI{
			homeworks: homeworks,
		},
	)
}

func (s *serverAPI) RecordHomework(ctx context.Context, req *journalv1.RecordHomeworkRequest) (*journalv1.Homework, error) {
	params, err := mapper.RecordHomeworkToDomain(req)
	if err != nil {
		return nil, err
	}

	homework, err := s.homeworks.RecordHomework(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.HomeworkToProto(homework), nil
}

func (s *serverAPI) UpdateHomework(ctx context.Context, req *journalv1.UpdateHomeworkRequest) (*journalv1.Homework, error) {
	params, err := mapper.UpdateHomeworkToDomain(req)
	if err != nil {
		return nil, err
	}

	homework, err := s.homeworks.UpdateHomework(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.HomeworkToProto(homework), nil
}

func (s *serverAPI) ListHomework(ctx context.Context, req *journalv1.ListHomeworkRequest) (*journalv1.ListHomeworkResponse, error) {
	params, err := mapper.ListHomeworkToDomain(req)
	if err != nil {
		return nil, err
	}

	listHomework, err := s.homeworks.ListHomework(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.ListHomeworkToProto(listHomework), nil
}

func (s *serverAPI) DeleteHomework(ctx context.Context, req *journalv1.DeleteHomeworkRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.homeworks.DeleteHomework(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
