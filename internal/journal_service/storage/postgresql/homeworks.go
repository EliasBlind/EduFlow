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

func (s *Storage) RecordHomework(
	ctx context.Context,
	params *domain.RecordHomework,
) (*domain.Homework, error) {
	const op = "storage.postgresql.RecordHomework"
	log := s.log.With(
		"op", op,
		"teacher_id", params.TeacherID.String(),
		"class_id", params.ClassID.String(),
		"subject_id", params.SubjectID,
	)

	argForTS := sqlgen.GetTeachingIDParams{
		TeacherID: mapper.ToPgUUID(params.TeacherID),
		SubjectID: mapper.ToPgUUID(params.SubjectID),
		ClassID:   mapper.ToPgUUID(params.ClassID),
	}

	tsId, err := s.queries.GetTeachingID(ctx, argForTS)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	arg := sqlgen.RecordHomeworkParams{
		TsID:            tsId,
		DescriptionTask: params.DescriptionTask,
		AssignedAt:      mapper.ToDate(params.Start),
		DeadlineAt:      mapper.ToDate(params.End),
	}

	homework, err := s.queries.RecordHomework(ctx, arg)
	if err != nil {
		log.Error("failed to record homework", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ts, err := s.queries.GetTeachingLoad(ctx, homework.TsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update homework", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Homework{
		ID:              homework.ID.Bytes,
		ClassID:         ts.ClassID.Bytes,
		TeacherID:       ts.TeacherID.Bytes,
		SubjectID:       ts.SubjectID.Bytes,
		DescriptionTask: homework.DescriptionTask,
		Start:           mapper.FromTimestamp(homework.AssignedAt),

		End: mapper.FromTimestamp(homework.DeadlineAt),
	}, nil
}

func (s *Storage) UpdateHomework(
	ctx context.Context,
	params *domain.UpdateHomework,
) (*domain.Homework, error) {
	const op = "storage.postgresql.UpdateHomework"
	log := s.log.With(
		"op", op,
		"id", params.ID.String(),
	)

	arg := sqlgen.UpdateHomeworkParams{
		ID:              mapper.ToPgUUID(params.ID),
		DescriptionTask: params.DescriptionTask,
		AssignedAt:      mapper.ToDate(params.Start),
		DeadlineAt:      mapper.ToDate(params.End),
	}

	homework, err := s.queries.UpdateHomework(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("homework not found or teacher mismatch", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update homework", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ts, err := s.queries.GetTeachingLoad(ctx, homework.TsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("homework not found or teacher mismatch", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update homework", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Homework{
		ID:              homework.ID.Bytes,
		ClassID:         ts.ClassID.Bytes,
		TeacherID:       ts.TeacherID.Bytes,
		SubjectID:       ts.SubjectID.Bytes,
		DescriptionTask: homework.DescriptionTask,
		Start:           mapper.FromTimestamp(homework.AssignedAt),
		End:             mapper.FromTimestamp(homework.DeadlineAt),
	}, nil
}

func (s *Storage) ListHomeworks(
	ctx context.Context,
	classId,
	subjectId uuid.UUID,
) ([]domain.Homework, error) {
	const op = "storage.postgresql.ListHomeworks"
	log := s.log.With(
		"op", op,
		"class_id", classId.String(),
		"subject_id", subjectId.String(),
	)

	arg := sqlgen.ListHomeworksParams{
		ClassID:   mapper.ToPgUUID(classId),
		SubjectID: mapper.ToPgUUID(subjectId),
	}

	homeworks, err := s.queries.ListHomeworks(ctx, arg)
	if err != nil {
		log.Error("failed to list homeworks", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mapper.MapSlice(
		homeworks,
		func(f *sqlgen.ListHomeworksRow) domain.Homework {
			return domain.Homework{
				ID:              f.ID.Bytes,
				ClassID:         f.ClassID.Bytes,
				TeacherID:       f.TeacherID.Bytes,
				SubjectID:       f.SubjectID.Bytes,
				DescriptionTask: f.DescriptionTask,

				Start: mapper.FromTimestamp(f.AssignedAt),
				End:   mapper.FromTimestamp(f.DeadlineAt),
			}
		},
	), nil
}

func (s *Storage) DeleteHomeworks(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteHomeworks"
	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	err := s.queries.DeleteHomeworks(ctx, mapper.ToPgUUID(id))
	if err != nil {
		log.Error("failed to delete homework", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
