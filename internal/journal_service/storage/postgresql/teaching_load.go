package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/sqlgen"
	mapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateTeachingLoad(
	ctx context.Context,
	params *domain.CreateTeachingLoad,
) (*domain.TeachingLoad, error) {
	const op = "storage.postgresql.CreateTeachingLoad"

	log := s.log.With(
		"op", op,
		"teacher_id", params.TeacherID.String(),
		"subject_id", params.SubjectID.String(),
		"class_id", params.ClassID.String(),
	)

	arg := sqlgen.CreateTeachingLoadParams{
		TeacherID: mapper.ToPgUUID(params.TeacherID),
		SubjectID: mapper.ToPgUUID(params.SubjectID),
		ClassID:   mapper.ToPgUUID(params.ClassID),
	}

	ts, err := s.queries.CreateTeachingLoad(ctx, arg)
	if err != nil {
		log.Error("failed to create teaching load", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.TeachingLoad{
		ID:        ts.ID.Bytes,
		TeacherID: ts.TeacherID.Bytes,
		SubjectID: ts.SubjectID.Bytes,
		ClassID:   ts.ClassID.Bytes,
	}, nil
}

func (s *Storage) GetTeachingLoad(
	ctx context.Context,
	id uuid.UUID,
) (*domain.TeachingLoad, error) {
	const op = "storage.postgresql.GetTeachingLoad"

	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	res, err := s.queries.GetTeachingLoad(ctx, mapper.ToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get teaching load", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.TeachingLoad{
		ID:        res.ID.Bytes,
		TeacherID: res.TeacherID.Bytes,
		SubjectID: res.SubjectID.Bytes,
		ClassID:   res.ClassID.Bytes,
	}, nil
}

func (s *Storage) ListTeachingLoad(
	ctx context.Context,
	params *domain.ListTeachingLoadRequest,
) ([]domain.TeachingLoad, error) {
	const op = "storage.postgresql.ListTeachingLoad"

	// Guard against nil-pointer logs if filters are optional
	var teacherIDStr, classIDStr string
	if params.TeacherID != nil {
		teacherIDStr = params.TeacherID.String()
	}
	if params.ClassID != nil {
		classIDStr = params.ClassID.String()
	}

	log := s.log.With(
		"op", op,
		"teacher_id", teacherIDStr,
		"class_id", classIDStr,
	)

	arg := sqlgen.ListTeachingLoadParams{
		TeacherID: mapper.PtrToPgUUID(params.TeacherID),
		ClassID:   mapper.PtrToPgUUID(params.ClassID),
	}

	loads, err := s.queries.ListTeachingLoad(ctx, arg)
	if err != nil {
		log.Error("failed to list teaching loads", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mapper.MapSliceRef(
		loads,
		func(f *sqlgen.TeacherSubject) domain.TeachingLoad {
			return domain.TeachingLoad{
				ID:        f.ID.Bytes,
				TeacherID: f.TeacherID.Bytes,
				SubjectID: f.SubjectID.Bytes,
				ClassID:   f.ClassID.Bytes,
			}
		},
	), nil
}

func (s *Storage) UpdateTeachingLoad(
	ctx context.Context,
	params *domain.UpdateTeachingLoad,
) (*domain.TeachingLoad, error) {
	const op = "storage.postgresql.UpdateTeachingLoad"

	log := s.log.With(
		"op", op,
		"id", params.ID.String(),
	)

	arg := sqlgen.UpdateTeachingLoadParams{
		ID:        mapper.ToPgUUID(params.ID),
		TeacherID: mapper.ToPgUUID(params.TeacherID),
	}

	res, err := s.queries.UpdateTeachingLoad(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("teaching load not found for update", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update teaching load", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.TeachingLoad{
		ID:        res.ID.Bytes,
		TeacherID: res.TeacherID.Bytes,
		SubjectID: res.SubjectID.Bytes,
		ClassID:   res.ClassID.Bytes,
	}, nil
}

func (s *Storage) DeleteTeachingLoad(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteTeachingLoad"

	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	err := s.queries.DeleteTeachingLoad(ctx, mapper.ToPgUUID(id))
	if err != nil {
		log.Error("failed to delete teaching load", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
