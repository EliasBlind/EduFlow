package studentssvc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostgresSql interface {
	CreateStudent(ctx context.Context, params *domain.Student) (*domain.Student, error)

	GetStudent(ctx context.Context, studentID uuid.UUID) (*domain.Student, error)

	ListStudents(ctx context.Context, params *domain.ListStudentsRequest) ([]domain.Student, error)

	UpdateStudent(ctx context.Context, params *domain.UpdateStudent) (*domain.Student, error)

	DeleteStudent(ctx context.Context, id uuid.UUID) error
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

func (a *Auth) CreateStudent(
	ctx context.Context,
	params *domain.Student,
) (*domain.Student, error) {
	const op = "auth.CreateStudent"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	student, err := a.sql.CreateStudent(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Debug("student already exists")
			return nil, domain.ErrAlreadyExists
		}
		log.Error("failed to create student", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("student created successfully", "student_id", student.ID)
	return student, nil
}

func (a *Auth) GetStudent(
	ctx context.Context,
	studentID uuid.UUID,
) (*domain.Student, error) {
	const op = "auth.GetStudent"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "student_id", studentID)

	student, err := a.sql.GetStudent(ctx, studentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("student not found")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get student", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("student retrieved successfully")
	return student, nil
}

func (a *Auth) ListStudents(
	ctx context.Context,
	params *domain.ListStudentsRequest,
) ([]domain.Student, error) {
	const op = "auth.ListStudents"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "class_id", params.ClassID)

	students, err := a.sql.ListStudents(ctx, params)
	if err != nil {
		log.Error("failed to list students", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("students listed successfully", "count", len(students))
	return students, nil
}

func (a *Auth) UpdateStudent(
	ctx context.Context,
	params *domain.UpdateStudent,
) (*domain.Student, error) {
	const op = "auth.UpdateStudent"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return nil, domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "student_id", params.ID)

	if !claims.Role.IsAdmin() {
		return nil, domain.ErrForbidden
	}

	if err = a.val.Struct(params); err != nil {
		log.Warn("invalid params", "error", err)
		return nil, domain.ErrInvalidData
	}

	student, err := a.sql.UpdateStudent(ctx, params)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("student not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update student", "error", err)
		return nil, domain.ErrInternal
	}

	log.Info("student updated successfully")
	return student, nil
}

func (a *Auth) DeleteStudent(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "auth.DeleteStudent"

	claims, err := domain.GetUserClaims(ctx)
	if err != nil {
		a.log.Warn("unauthorized attempt", "op", op, "error", err)
		return domain.ErrUnauthorized
	}

	log := a.log.With("op", op, "user_id", claims.ID, "student_id", id)

	if !claims.Role.IsAdmin() {
		return domain.ErrForbidden
	}

	err = a.sql.DeleteStudent(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Debug("student not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to delete student", "error", err)
		return domain.ErrInternal
	}

	log.Info("student deleted successfully")
	return nil
}
