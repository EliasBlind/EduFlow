package statuscodesvc

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
	CreateStatusCode(ctx context.Context, statusCodeName string) (*domain.StatusCode, error)

	UpdateStatusCode(ctx context.Context, params *domain.UpdateStatusCode) (*domain.StatusCode, error)

	ListStatusCode(ctx context.Context) ([]domain.StatusCode, error)

	DeleteStatusCode(ctx context.Context, id uuid.UUID) error
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

func (a *Auth) CreateStatusCode(
	ctx context.Context,
	statusCodeName string,
) (*domain.StatusCode, error) {
	const op = "auth.CreateStatusCode"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"status_name", statusCodeName,
	)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	statusCode, err := a.sql.CreateStatusCode(ctx, statusCodeName)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Debug("status code already exists")
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to create status code", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("status code created successfully", "id", statusCode.ID)
	return statusCode, nil
}

func (a *Auth) UpdateStatusCode(
	ctx context.Context,
	params *domain.UpdateStatusCode,
) (*domain.StatusCode, error) {
	const op = "auth.UpdateStatusCode"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"status_id", params.ID,
	)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	statusCode, err := a.sql.UpdateStatusCode(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("status code not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update status code", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("status code updated successfully")
	return statusCode, nil
}

func (a *Auth) ListStatusCode(
	ctx context.Context,
) ([]domain.StatusCode, error) {
	const op = "auth.ListStatusCode"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
	)

	statusCodes, err := a.sql.ListStatusCode(ctx)
	if err != nil {
		log.Error("failed to list status codes", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("status codes listed successfully", "count", len(statusCodes))
	return statusCodes, nil
}

func (a *Auth) DeleteStatusCode(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteStatusCode"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With(
		"op", op,
		"user_id", claims.ID,
		"status_code_id", id,
	)

	if !claims.Role.IsAdmin() {
		return domain.ErrForbidden
	}

	err = a.sql.DeleteStatusCode(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("status code not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete status code", "error", err)
		return domain.ErrInternal
	}

	log.Info("status code deleted successfully")
	return nil
}
