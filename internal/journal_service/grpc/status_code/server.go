package statuscode

import (
	"context"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/mapper"
	journalv1 "github.com/EliasBlind/EduFlow/pkg/protos/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StatusCodes interface {
	CreateStatusCode(ctx context.Context, statusCodeName string) (*domain.StatusCode, error)

	UpdateStatusCode(ctx context.Context, req *domain.UpdateStatusCode) (*domain.StatusCode, error)

	ListStatusCode(ctx context.Context) ([]domain.StatusCode, error)

	DeleteStatusCode(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedStatusCodeServiceServer
	statusCodes StatusCodes
}

func Register(gRPC *grpc.Server, statusCodes StatusCodes) {
	journalv1.RegisterStatusCodeServiceServer(
		gRPC,
		&serverAPI{
			statusCodes: statusCodes,
		},
	)
}

func (s *serverAPI) CreateStatusCode(ctx context.Context, req *journalv1.CreateStatusCodeRequest) (*journalv1.StatusCode, error) {
	statusCode, err := s.statusCodes.CreateStatusCode(ctx, req.FullName)
	if err != nil {
		return nil, err
	}
	return mapper.StatusCodeToProto(statusCode), nil
}

func (s *serverAPI) UpdateStatusCode(ctx context.Context, req *journalv1.UpdateStatusCodeRequest) (*journalv1.StatusCode, error) {
	params, err := mapper.UpdateStatusCodeToDomain(req)
	if err != nil {
		return nil, err
	}

	statusCode, err := s.statusCodes.UpdateStatusCode(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.StatusCodeToProto(statusCode), nil
}

func (s *serverAPI) ListStatusCode(ctx context.Context, req *journalv1.ListStatusCodeRequest) (*journalv1.ListStatusCodeResponse, error) {
	statusCodes, err := s.statusCodes.ListStatusCode(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ListStatusCodeToProto(statusCodes), nil
}

func (s *serverAPI) DeleteStatusCode(ctx context.Context, req *journalv1.DeleteStatusCodeRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.statusCodes.DeleteStatusCode(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
