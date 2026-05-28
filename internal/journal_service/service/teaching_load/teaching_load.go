package teachingloadsvc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	usercalimas "github.com/EliasBlind/EduFlow/pkg/user_calimas"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	CreateTeachingLoad(ctx context.Context, params *domain.CreateTeachingLoad) (*domain.TeachingLoad, error)

	GetTeachingLoad(ctx context.Context, id uuid.UUID) (*domain.TeachingLoad, error)

	ListTeachingLoad(ctx context.Context, params *domain.ListTeachingLoadRequest) ([]domain.TeachingLoad, error)

	UpdateTeachingLoad(ctx context.Context, params *domain.UpdateTeachingLoad) (*domain.TeachingLoad, error)

	DeleteTeachingLoad(ctx context.Context, id uuid.UUID) error
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

func (a *Auth) CreateTeachingLoad(
	ctx context.Context,
	params *domain.CreateTeachingLoad,
) (*domain.TeachingLoad, error) {
	const op = "auth.CreateTeachingLoad"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"teacher_id", params.TeacherID,
		"class_id", params.ClassID,
	)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	teachingLoad, err := a.sql.CreateTeachingLoad(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Debug("teaching load already exists")
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to create teaching load", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teaching load created successfully", "id", teachingLoad.ID)
	return teachingLoad, nil
}

func (a *Auth) GetTeachingLoad(
	ctx context.Context,
	id uuid.UUID,
) (*domain.TeachingLoad, error) {
	const op = "auth.GetTeachingLoad"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "load_id", id)

	teachingLoad, err := a.sql.GetTeachingLoad(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("teaching load not found")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get teaching load", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teaching load retrieved successfully")
	return teachingLoad, nil
}

func (a *Auth) ListTeachingLoad(
	ctx context.Context,
	params *domain.ListTeachingLoadRequest,
) ([]domain.TeachingLoad, error) {
	const op = "auth.ListTeachingLoad"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID)

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	teachingLoads, err := a.sql.ListTeachingLoad(ctx, params)
	if err != nil {
		log.Error("failed to list teaching loads", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teaching loads listed successfully", "count", len(teachingLoads))
	return teachingLoads, nil
}

func (a *Auth) UpdateTeachingLoad(
	ctx context.Context,
	params *domain.UpdateTeachingLoad,
) (*domain.TeachingLoad, error) {
	const op = "auth.UpdateTeachingLoad"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	log := a.log.With("op", op, "user_id", claims.ID, "load_id", params.ID)

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	teachingLoad, err := a.sql.UpdateTeachingLoad(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("teaching load not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update teaching load", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("teaching load updated successfully")
	return teachingLoad, nil
}

func (a *Auth) DeleteTeachingLoad(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteTeachingLoad"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "load_id", id)

	if !claims.Role.IsAdmin() {
		return domain.ErrForbidden
	}

	err = a.sql.DeleteTeachingLoad(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("teaching load not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete teaching load", "error", err)
		return domain.ErrInternal
	}

	log.Info("teaching load deleted successfully")
	return nil
}
