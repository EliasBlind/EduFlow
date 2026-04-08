package grades

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

type Grades interface {
	RecordGrade(ctx context.Context, req *domain.RecordGrade) (*domain.Grade, error)

	ListGrades(ctx context.Context, req *domain.ListGradesRequest) (*domain.ListGradesResponse, error)

	UpdateGrade(ctx context.Context, req *domain.UpdateGrade) (*domain.Grade, error)

	DeleteGrade(ctx context.Context, id uuid.UUID) error
}

type serverAPI struct {
	journalv1.UnimplementedGradesServiceServer
	grades Grades
	log    *slog.Logger
}

func Register(gRPC *grpc.Server, logger *slog.Logger, grades Grades) {
	journalv1.RegisterGradesServiceServer(
		gRPC,
		&serverAPI{
			grades: grades,
			log:    logger,
		},
	)
}

func (s *serverAPI) RecordGrade(ctx context.Context, req *journalv1.RecordGradeRequest) (*journalv1.Grade, error) {
	params, err := mapper.RecordGradeToDomain(req)
	if err != nil {
		return nil, err
	}

	grade, err := s.grades.RecordGrade(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.GradeToProto(grade), nil
}

func (s *serverAPI) ListGrades(ctx context.Context, req *journalv1.ListGradesRequest) (*journalv1.ListGradesResponse, error) {
	params, err := mapper.ListGradesToDomain(req)
	if err != nil {
		return nil, err
	}

	listGrades, err := s.grades.ListGrades(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.ListGradesToProto(listGrades), nil
}

func (s *serverAPI) UpdateGrade(ctx context.Context, req *journalv1.UpdateGradeRequest) (*journalv1.Grade, error) {
	params, err := mapper.UpdateGradeToDomain(req)
	if err != nil {
		return nil, err
	}

	grade, err := s.grades.UpdateGrade(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapper.GradeToProto(grade), nil
}

func (s *serverAPI) DeleteGrade(ctx context.Context, req *journalv1.DeleteGradeRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	
	err = s.grades.DeleteGrade(ctx, id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
