package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/sqlgen"
	mapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateClass(
	ctx context.Context,
	params *domain.CreateClass,
) (*domain.Class, error) {
	const op = "storage.postgresql.CreateClass"

	log := s.log.With(
		"op", op,
		"class_name", params.ClassName,
		"years_of_study", params.YearOfStudy,
		"graduation_year", params.GraduationYear,
	)

	if params.GraduationYear == nil {
		log.Error("graduation year is required")
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidData)
	}

	arg := sqlgen.CreateClassParams{
		ClassName:      params.ClassName,
		YearOfStudy:    int32(params.YearOfStudy),
		GraduationYear: int32(*params.GraduationYear),
	}

	class, err := s.queries.CreateClass(
		ctx,
		arg,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Warn("class already exists", "err", err)
			return nil, fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}

		log.Error("failed to create class", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Class{
		ID:             uuid.UUID(class.ID.Bytes),
		ClassName:      class.ClassName,
		YearOfStudy:    uint32(class.YearOfStudy),
		GraduationYear: uint32(class.GraduationYear),
	}, nil
}

func (s *Storage) GetClass(
	ctx context.Context,
	classId uuid.UUID,
) (*domain.Class, error) {
	const op = "storage.postgresql.GetClass"

	log := s.log.With(
		"op", op,
		"class_id", classId.String(),
	)

	class, err := s.queries.GetClass(ctx, mapper.ToPgUUID(classId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("class not found", "id", classId)

			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}

		log.Error("failed to get class", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Class{
		ID:             class.ID.Bytes,
		ClassName:      class.ClassName,
		YearOfStudy:    uint32(class.YearOfStudy),
		GraduationYear: uint32(class.GraduationYear),
	}, nil
}

func (s *Storage) ListClasses(ctx context.Context) ([]domain.Class, error) {
	const op = "storage.postgresql.ListClasses"

	log := s.log.With(
		"op", op,
	)

	classes, err := s.queries.ListClasses(ctx)
	if err != nil {
		log.Error("failed to list classes", "error", err)
		return nil, err
	}

	return mapper.MapSliceRef(
		classes,
		func(f *sqlgen.Class) domain.Class {
			return domain.Class{
				ID:             uuid.UUID(f.ID.Bytes),
				ClassName:      f.ClassName,
				YearOfStudy:    uint32(f.YearOfStudy),
				GraduationYear: uint32(f.GraduationYear),
			}
		},
	), nil
}

func (s *Storage) ListTeacherClasses(
	ctx context.Context,
	params *domain.ListTeacherClasses,
) ([]domain.Class, error) {
	const op = "storage.postgresql.ListTeacherClasses"

	log := s.log.With(
		"op", op,
		"teacher_id", params.TeacherID,
	)

	arg := sqlgen.ListTeacherClassesParams{
		TeacherID: mapper.ToPgUUID(params.TeacherID),
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
	}

	classes, err := s.queries.ListTeacherClasses(ctx, arg)
	if err != nil {
		log.Error("failed to list classes", "error", err)
		return nil, err
	}

	return mapper.MapSliceRef(
		classes,
		func(f *sqlgen.Class) domain.Class {
			return domain.Class{
				ID:             uuid.UUID(f.ID.Bytes),
				ClassName:      f.ClassName,
				YearOfStudy:    uint32(f.YearOfStudy),
				GraduationYear: uint32(f.GraduationYear),
			}
		},
	), nil
}

func (s *Storage) UpdateClass(
	ctx context.Context,
	param *domain.UpdateClass,
) (*domain.Class, error) {
	const op = "storage.postgresql.UpdateClass"

	log := s.log.With(
		"op", op,
		"class_id", param.ID.String(),
	)

	arg := sqlgen.UpdateClassParams{
		ID:             mapper.ToPgUUID(param.ID),
		ClassName:      param.ClassName,
		YearOfStudy:    mapper.Ptr(int32(*param.YearOfStudy)),
		GraduationYear: mapper.Ptr(int32(*param.GraduationYear)),
	}

	class, err := s.queries.UpdateClass(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("class not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update class", "error", err)
		return nil, err
	}
	return &domain.Class{
		ID:             class.ID.Bytes,
		ClassName:      class.ClassName,
		YearOfStudy:    uint32(class.YearOfStudy),
		GraduationYear: uint32(class.GraduationYear),
	}, nil
}

func (s *Storage) DeleteClass(
	ctx context.Context,
	classId uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteClass"

	log := s.log.With(
		"op", op,
		"class_id", classId.String(),
	)
	err := s.queries.DeleteClass(ctx, mapper.ToPgUUID(classId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("class not found for update")
			return domain.ErrNotFound
		}
		log.Error("failed to update class", "error", err)
		return err
	}
	return nil
}
