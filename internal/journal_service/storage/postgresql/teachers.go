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

func (s *Storage) CreateTeacher(
	ctx context.Context,
	params *domain.Teacher,
) (*domain.Teacher, error) {
	const op = "storage.postgresql.CreateTeacher"

	log := s.log.With(
		"op", op,
		"teacher_name", params.FullName,
	)

	arg := sqlgen.CreateTeacherParams{
		FullName: params.FullName,
	}

	if params.ID != nil {
		arg.ID = mapper.ToPgUUID(*params.ID)
	}

	res, err := s.queries.CreateTeacher(ctx, arg)
	if err != nil {
		log.Error("failed to create teacher", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Teacher{
		ID:       pkgmapper.ToPtrUUID(res.ID.Bytes),
		FullName: res.FullName,
	}, nil
}

func (s *Storage) ListTeachers(
	ctx context.Context,
	params *domain.ListTeachersRequest,
) ([]domain.Teacher, error) {
	const op = "storage.postgresql.ListTeachers"

	log := s.log.With(
		"op", op,
	)

	arg := sqlgen.ListTeachersParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	}

	teachers, err := s.queries.ListTeachers(ctx, arg)
	if err != nil {
		log.Error("failed to list teachers", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mapper.MapSliceRef(
		teachers,
		func(f *sqlgen.Teacher) domain.Teacher {
			return domain.Teacher{
				ID:       pkgmapper.ToPtrUUID(f.ID.Bytes),
				FullName: f.FullName,
			}
		},
	), nil
}

func (s *Storage) UpdateTeacher(
	ctx context.Context,
	params *domain.UpdateTeacher,
) (*domain.Teacher, error) {
	const op = "storage.postgresql.UpdateTeacher"

	log := s.log.With(
		"op", op,
		"id", params.ID.String(),
		"full_name", params.FullName,
	)

	arg := sqlgen.UpdateTeacherParams{
		ID:       mapper.ToPgUUID(params.ID),
		FullName: params.FullName,
	}

	res, err := s.queries.UpdateTeacher(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("teacher not found for update", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update teacher", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Teacher{
		ID:       pkgmapper.ToPtrUUID(res.ID.Bytes),
		FullName: res.FullName,
	}, nil
}

func (s *Storage) DeleteTeacher(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteTeacher"

	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	err := s.queries.DeleteTeacher(ctx, mapper.ToPgUUID(id))
	if err != nil {
		log.Error("failed to delete teacher", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
