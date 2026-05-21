package classessvc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	CreateClass(ctx context.Context, params *domain.CreateClass) (*domain.Class, error)
	GetClass(ctx context.Context, classId uuid.UUID) (*domain.Class, error)

	ListClasses(ctx context.Context) ([]domain.Class, error)
	ListTeacherClasses(ctx context.Context, params *domain.ListTeacherClasses) ([]domain.Class, error)
	UpdateClass(ctx context.Context, param *domain.UpdateClass) (*domain.Class, error)
	DeleteClass(ctx context.Context, classId uuid.UUID) error
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

func (a *Auth) CreateClass(
	ctx context.Context,
	params *domain.CreateClass,
) (*domain.Class, error) {
	const op = "auth.CreateClass"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"role", claims.Role,
	)

	log.Info("attempting to create class", "class_name", params.ClassName)

	if !claims.Role.IsAdmin() {
		log.Warn("permission denied", "required_role", "admin")
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	class, err := a.sql.CreateClass(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to create class", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("class created successfully", "class_id", class.ID)
	return class, nil
}

func (a *Auth) GetClass(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Class, error) {
	const op = "auth.GetClass"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("failed to get user claims", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"role", claims.Role,
		"class_id", id,
	)

	class, err := a.sql.GetClass(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("class not found in database")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get class from database", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("successfully retrieved class")
	return class, nil
}

func (a *Auth) ListClasses(ctx context.Context) ([]domain.Class, error) {
	const op = "auth.ListClasses"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("failed to get user claims", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"role", claims.Role,
	)

	if !claims.Role.IsAdmin() {
		log.Warn("permission denied: listing all classes requires admin role")
		return nil, domain.ErrForbidden
	}

	classes, err := a.sql.ListClasses(ctx)
	if err != nil {
		log.Error("failed to list classes from database", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("successfully retrieved all system classes", "count", len(classes))
	return classes, nil
}

func (a *Auth) ListTeacherClasses(
	ctx context.Context,
	params *domain.ListTeacherClasses,
) ([]domain.Class, error) {
	const op = "auth.ListTeacherClasses"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("failed to get user claims", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"role", claims.Role,
		"target_teacher_id", params.TeacherID,
	)

	log.Debug("attempting to list teacher classes")

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	classes, err := a.sql.ListTeacherClasses(ctx, params)
	if err != nil {
		log.Error("failed to list classes from database", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("successfully retrieved teacher classes")
	return classes, nil
}

func (a *Auth) UpdateClass(
	ctx context.Context,
	params *domain.UpdateClass,
) (*domain.Class, error) {
	const op = "auth.UpdateClass"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("failed to get user claims", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"role", claims.Role,
		"class_id", params.ID,
	)

	log.Debug("attempting to update class")

	if !claims.Role.IsAdmin() {
		log.Warn("permission denied: update requires admin role")
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid update params", "error", err)
		return nil, domain.ErrInvalidData
	}

	class, err := a.sql.UpdateClass(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("class to update not found")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update class in database", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("class updated successfully")
	return class, nil
}

func (a *Auth) DeleteClass(
	ctx context.Context,
	classId uuid.UUID,
) error {
	const op = "auth.DeleteClass"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("failed to get user claims", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"role", claims.Role,
		"class_id", classId,
	)

	log.Debug("attempting to delete class")

	if !claims.Role.IsAdmin() {
		log.Warn("permission denied: delete requires admin role")
		return domain.ErrForbidden
	}

	err = a.sql.DeleteClass(ctx, classId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("class to delete not found")
			return domain.ErrNotFound
		}
		log.Error("failed to delete class from database", "error", err)
		return domain.ErrInternal
	}

	log.Info("class deleted successfully")
	return nil
}
