package gradessvc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	RecordGrade(ctx context.Context, params *domain.RecordGrade) (*domain.Grade, error)

	GetGrades(ctx context.Context, teacherId uuid.UUID, params *domain.ListGradesRequest) ([]domain.Grade, error)

	UpdateGrades(ctx context.Context, grade *domain.UpdateGrade) (*domain.Grade, error)

	DeleteGrade(ctx context.Context, id uuid.UUID) error
}

type Auth struct {
	log *slog.Logger
	val *validator.Validate
	sql PostgresSql
}

func New(
	log *slog.Logger,
	val *validator.Validate,
	sql PostgresSql,
) *Auth {
	return &Auth{
		log: log,
		val: val,
		sql: sql,
	}
}

func (a *Auth) RecordGrade(ctx context.Context, params *domain.RecordGrade) (*domain.Grade, error) {
	const op = "auth.RecordGrade"
	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "student_id", params.StudentID)

	if err := a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	if params.DateOfGrade == nil {
		now := time.Now()
		params.DateOfGrade = &now
	}

	grade, err := a.sql.RecordGrade(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to save record", "error", err)
		return nil, domain.ErrInternal
	}

	return grade, nil
}
func (a *Auth) ListGrades(
	ctx context.Context,
	params *domain.ListGradesRequest,
) ([]domain.Grade, error) {
	const op = "auth.ListGrades"
	log := a.log.With("op", op, "student_id", params.StudentID, "class_id", params.ClassID)

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		log.Warn("unauthorized attempt")
		return nil, domain.ErrUnauthorized
	}

	grades, err := a.sql.GetGrades(ctx, claims.ID, params)

	if err != nil {
		log.Error("failed to get grades", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("grades listed successfully", "count", len(grades))
	return grades, nil
}

func (a *Auth) UpdateGrade(
	ctx context.Context,
	params *domain.UpdateGrade,
) (*domain.Grade, error) {
	const op = "auth.UpdateGrade"
	log := a.log.With(
		"op", op,
		"Note", params.Note)

	_, err := domain.GetUserClaims(ctx)
	if err != nil {
		log.Warn("unauthorized attempt")
		return nil, domain.ErrUnauthorized
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	grade, err := a.sql.UpdateGrades(ctx, params)

	if err != nil {
		log.Error("failed to update grade", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("grade updated successfully")
	return grade, nil
}

func (a *Auth) DeleteGrade(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteGrade"
	log := a.log.With("op", op, "id", id)

	_, err := domain.GetUserClaims(ctx)
	if err != nil {
		log.Warn("unauthorized attempt")
		return domain.ErrUnauthorized
	}

	err = a.sql.DeleteGrade(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("grade not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete grade", "error", err)
		return domain.ErrInternal
	}

	log.Info("grade deleted successfully")
	return nil
}
