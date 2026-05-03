package teachingload

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/mapper"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TeachingLoad interface {
	CreateTeachingLoad(ctx context.Context, req *domain.CreateTeachingLoad) (*domain.TeachingLoad, error)

	GetTeachingLoad(ctx context.Context, id uuid.UUID) (*domain.TeachingLoad, error)

	ListTeachingLoad(ctx context.Context, req *domain.ListTeachingLoadRequest) (*domain.ListTeachingLoadResponse, error)

	UpdateTeachingLoad(ctx context.Context, req *domain.UpdateTeachingLoad) (*domain.TeachingLoad, error)

	DeleteTeachingLoad(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedTeachingLoadServiceServer
	teachingLoad TeachingLoad
}

func Register(gRPC *grpc.Server, teachingLoad TeachingLoad) {
	journalv1.RegisterTeachingLoadServiceServer(
		gRPC,
		&serverAPI{
			teachingLoad: teachingLoad,
		},
	)
}

func (s *serverAPI) CreateTeachingLoad(ctx context.Context, req *journalv1.CreateTeachingLoadRequest) (*journalv1.TeachingLoad, error) {
	params, err := mapper.CreateTeachingLoadToDomain(req)
	if err != nil {
		return nil, err
	}

	teachingLoad, err := s.teachingLoad.CreateTeachingLoad(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.TeachingLoadToProto(teachingLoad), nil
}

func (s *serverAPI) GetTeachingLoad(ctx context.Context, req *journalv1.GetTeachingLoadRequest) (*journalv1.TeachingLoad, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	teachingLoad, err := s.teachingLoad.GetTeachingLoad(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.TeachingLoadToProto(teachingLoad), nil
}

func (s *serverAPI) ListTeachingLoad(ctx context.Context, req *journalv1.ListTeachingLoadRequest) (*journalv1.ListTeachingLoadResponse, error) {
	params, err := mapper.ListTeachingLoadToDomain(req)
	if err != nil {
		return nil, err
	}

	listTeachingLoad, err := s.teachingLoad.ListTeachingLoad(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.ListTeachingLoadToProto(listTeachingLoad), nil
}

func (s *serverAPI) UpdateTeachingLoad(ctx context.Context, req *journalv1.UpdateTeachingLoadRequest) (*journalv1.TeachingLoad, error) {
	params, err := mapper.UpdateTeachingLoadToDomain(req)
	if err != nil {
		return nil, err
	}

	teachingLoad, err := s.teachingLoad.UpdateTeachingLoad(ctx, params)
	if err != err {
		return nil, err
	}
	return mapper.TeachingLoadToProto(teachingLoad), nil
}

func (s *serverAPI) DeleteTeachingLoad(ctx context.Context, req *journalv1.DeleteTeachingLoadRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.teachingLoad.DeleteTeachingLoad(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
