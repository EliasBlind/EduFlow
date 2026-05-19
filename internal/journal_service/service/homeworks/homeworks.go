package homeworkssvc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	RecordHomework(ctx context.Context, homeworkData *domain.RecordHomework) (*domain.Homework, error)
	UpdateHomework(ctx context.Context, newHomeworkData *domain.UpdateHomework) (*domain.Homework, error)
	ListHomeworks(ctx context.Context, classId, subjectId uuid.UUID) ([]domain.Homework, error)
	DeleteHomeworks(ctx context.Context, id uuid.UUID) error
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
func (a *Auth) RecordHomework(
	ctx context.Context,
	params *domain.RecordHomework,
) (*domain.Homework, error) {
	const op = "auth.RecordHomework"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"class_id", params.ClassID,
		"subject_id", params.SubjectID,
	)

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	homework, err := a.sql.RecordHomework(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to record homework", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("homework recorded successfully", "homework_id", homework.ID)
	return homework, nil
}

func (a *Auth) UpdateHomework(
	ctx context.Context,
	params *domain.UpdateHomework,
) (*domain.Homework, error) {
	const op = "auth.UpdateHomework"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"homework_id", params.ID,
	)

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	homework, err := a.sql.UpdateHomework(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("homework not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update homework", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("homework updated successfully")
	return homework, nil
}

func (a *Auth) ListHomeworks(
	ctx context.Context,
	params *domain.ListHomeworkRequest,
) ([]domain.Homework, error) {
	const op = "auth.ListHomeworks"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"class_id", params.ClassID,
		"subject_id", params.SubjectID,
	)

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	homeworks, err := a.sql.ListHomeworks(ctx, params.ClassID, params.SubjectID)
	if err != nil {
		log.Error("failed to list homeworks", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("homeworks listed successfully", "count", len(homeworks))
	return homeworks, nil
}

func (a *Auth) DeleteHomework(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteHomework"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"homework_id", id,
	)

	err = a.sql.DeleteHomeworks(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("homework not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete homework", "error", err)
		return domain.ErrInternal
	}

	log.Info("homework deleted successfully")
	return nil
}
