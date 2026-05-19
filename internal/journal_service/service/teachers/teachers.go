package teacherssvc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	CreateTeacher(ctx context.Context, params *domain.Teacher) (*domain.Teacher, error)

	ListTeachers(ctx context.Context, params *domain.ListTeachersRequest) ([]domain.Teacher, error)

	UpdateTeacher(ctx context.Context, params *domain.UpdateTeacher) (*domain.Teacher, error)

	DeleteTeacher(ctx context.Context, id uuid.UUID) error
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

func (a *Auth) CreateTeacher(
	ctx context.Context,
	params *domain.Teacher,
) (*domain.Teacher, error) {
	const op = "auth.CreateTeacher"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "teacher_name", params.FullName)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	teacher, err := a.sql.CreateTeacher(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Debug("teacher already exists")
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to create teacher", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teacher created successfully", "teacher_id", teacher.ID)
	return teacher, nil
}

func (a *Auth) ListTeachers(
	ctx context.Context,
	params *domain.ListTeachersRequest,
) ([]domain.Teacher, error) {
	const op = "auth.ListTeachers"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID)

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	teachers, err := a.sql.ListTeachers(ctx, params)
	if err != nil {
		log.Error("failed to list teachers", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teachers listed successfully", "count", len(teachers))
	return teachers, nil
}

func (a *Auth) UpdateTeacher(
	ctx context.Context,
	params *domain.UpdateTeacher,
) (*domain.Teacher, error) {
	const op = "auth.UpdateTeacher"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "teacher_id", params.ID)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	teacher, err := a.sql.UpdateTeacher(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("teacher not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update teacher", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teacher updated successfully")
	return teacher, nil
}

func (a *Auth) DeleteTeacher(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteTeacher"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "teacher_id", id)

	if !claims.Role.IsAdmin() {
		return domain.ErrForbidden
	}

	err = a.sql.DeleteTeacher(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("teacher not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete teacher", "error", err)
		return domain.ErrInternal
	}

	log.Info("teacher deleted successfully")
	return nil
}
