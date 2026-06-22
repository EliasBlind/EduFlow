package subjectssvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	usercalimas "github.com/EliasBlind/EduFlow/pkg/user_calimas"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	CreateSubject(ctx context.Context, subjectName string) (*domain.Subject, error)

	UpdateSubject(ctx context.Context, params *domain.UpdateSubject) (*domain.Subject, error)

	DeleteSubject(ctx context.Context, id uuid.UUID) error

	ListSubjects(ctx context.Context) ([]domain.Subject, error)
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
func (a *Auth) CreateSubject(
	ctx context.Context,
	subjectName string,
) (*domain.Subject, error) {
	const op = "auth.CreateSubject"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "subject_name", subjectName)

	if !claims.Role.IsAdmin() {
		log.Warn("permission denied: operational access requires admin role",
			slog.String("user_login", claims.Login),
			slog.String("rejected_role", fmt.Sprintf("%q", claims.Role)),
		)
		return nil, domain.ErrForbidden
	}

	newSubject, err := a.sql.CreateSubject(ctx, subjectName)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Debug("subject already exists")
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to create subject", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("subject created successfully", "subject_id", newSubject.ID)
	return newSubject, nil
}

func (a *Auth) UpdateSubject(
	ctx context.Context,
	params *domain.UpdateSubject,
) (*domain.Subject, error) {
	const op = "auth.UpdateSubject"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "subject_id", params.ID)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	subject, err := a.sql.UpdateSubject(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("subject not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update subject", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("subject updated successfully")
	return subject, nil
}

func (a *Auth) DeleteSubject(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteSubject"

	claims, err := usercalimas.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "subject_id", id)

	if !claims.Role.IsAdmin() {
		return domain.ErrForbidden
	}

	err = a.sql.DeleteSubject(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("subject not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete subject", "error", err)
		return domain.ErrInternal
	}

	log.Info("subject deleted successfully")
	return nil
}

func (a *Auth) ListSubjects(ctx context.Context) ([]domain.Subject, error) {
	subjects, err := a.sql.ListSubjects(ctx)
	if err != nil {
		return nil, domain.ErrInternal
	}
	return subjects, nil
}
