package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/sqlgen"
	mapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	pkgmapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateStudent(
	ctx context.Context,
	params *domain.Student,
) (*domain.Student, error) {
	const op = "storage.postgresql.CreateStudent"

	log := s.log.With(
		"op", op,
		"class_id", params.ClassID.String(),
		"full_name", params.FullName,
	)

	arg := sqlgen.CreateStudentParams{
		ClassID:  mapper.PtrToPgUUID(params.ClassID),
		FullName: params.FullName,
	}

	if params.ID != nil {
		arg.ID = mapper.PtrToPgUUID(params.ID)
	}

	class, err := s.queries.CreateStudent(ctx, arg)
	if err != nil {
		log.Error("failed to create student", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Student{
		ID:       pkgmapper.ToPtrUUID(class.ID.Bytes),
		ClassID:  mapper.Ptr(uuid.UUID(class.ClassID.Bytes)),
		FullName: class.FullName,
	}, nil
}

func (s *Storage) GetStudent(
	ctx context.Context,
	studentID uuid.UUID,
) (*domain.Student, error) {
	const op = "storage.postgresql.GetStudent"

	log := s.log.With(
		"op", op,
		"student_id", studentID.String(),
	)

	student, err := s.queries.GetStudent(ctx, mapper.ToPgUUID(studentID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get student", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Student{
		ID:       pkgmapper.ToPtrUUID(student.ID.Bytes),
		ClassID:  mapper.Ptr(uuid.UUID(student.ClassID.Bytes)),
		FullName: student.FullName,
	}, nil
}

func (s *Storage) ListStudents(
	ctx context.Context,
	params *domain.ListStudentsRequest,
) ([]domain.Student, error) {
	const op = "storage.postgresql.ListStudents"

	log := s.log.With(
		"op", op,
		"class_id", params.ClassID.String(),
	)

	arg := sqlgen.ListStudentsParams{
		ClassID: mapper.ToPgUUID(params.ClassID),
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	}

	students, err := s.queries.ListStudents(ctx, arg)
	if err != nil {
		log.Error("failed to list students", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mapper.MapSlice(
		students,
		func(f *sqlgen.Student) domain.Student {
			return domain.Student{
				ID:       pkgmapper.ToPtrUUID(f.ID.Bytes),
				ClassID:  mapper.Ptr(uuid.UUID(f.ClassID.Bytes)),
				FullName: f.FullName,
			}
		},
	), nil
}

func (s *Storage) UpdateStudent(
	ctx context.Context,
	params *domain.UpdateStudent,
) (*domain.Student, error) {
	const op = "storage.postgresql.UpdateStudent"

	var classIDStr, fullNameStr string
	if params.ClassID != nil {
		classIDStr = params.ClassID.String()
	}
	if params.FullName != nil {
		fullNameStr = *params.FullName
	}

	log := s.log.With(
		"op", op,
		"id", params.ID.String(),
		"class_id", classIDStr,
		"full_name", fullNameStr,
	)

	arg := sqlgen.UpdateStudentParams{
		ID:       mapper.ToPgUUID(params.ID),
		ClassID:  mapper.ToPgUUID(*params.ClassID),
		FullName: params.FullName,
	}

	res, err := s.queries.UpdateStudent(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("student not found for update", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update student", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Student{
		ID:       pkgmapper.ToPtrUUID(res.ID.Bytes),
		ClassID:  mapper.Ptr(uuid.UUID(res.ClassID.Bytes)),
		FullName: res.FullName,
	}, nil
}

func (s *Storage) DeleteStudent(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteStudent"

	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	err := s.queries.DeleteStudent(ctx, mapper.ToPgUUID(id))
	if err != nil {
		log.Error("failed to delete student", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
