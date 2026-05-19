package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
	"github.com/EliasBlind/EduFlow/internal/journal_service/storage/sqlgen"
	mapper "github.com/EliasBlind/EduFlow/pkg/mappers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Storage) RecordGrade(
	ctx context.Context,
	params *domain.RecordGrade,
) (*domain.Grade, error) {
	const op = "storage.postgresql.RecordGrade"

	var statusCodeStr string
	if params.StatusCodeID != nil {
		statusCodeStr = params.StatusCodeID.String()
	}

	log := s.log.With(
		slog.String("op", op),
		slog.String("student_id", params.StudentID.String()),
		slog.String("ts_id", params.TsID.String()),
		slog.Any("grade", params.Grade),
		slog.String("status_code", statusCodeStr),
		slog.Uint64("lesson_number", uint64(params.LessonNumber)),
	)

	log.Debug("attempting to record new grade")

	var statusCodePg pgtype.UUID
	if params.StatusCodeID != nil {
		statusCodePg = mapper.ToPgUUID(*params.StatusCodeID)
	}

	// Безопасное извлечение score для параметров запроса
	var scorePtr *int16
	if params.Grade != nil {
		scorePtr = mapper.Ptr(int16(*params.Grade))
	}

	arg := sqlgen.RecordGradeParams{
		StudentID:    mapper.ToPgUUID(params.StudentID),
		TsID:         mapper.ToPgUUID(params.TsID),
		StatusCodeID: statusCodePg,
		Score:        scorePtr,
		LessonNumber: int16(params.LessonNumber),
		LessonDate:   mapper.ToDate(params.DateOfGrade),
	}

	grade, err := s.queries.RecordGrade(ctx, arg)
	if err != nil {
		log.Error("failed to execute record grade query", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("grade recorded successfully, fetching teaching load context", slog.Any("ts_id", grade.TsID))

	ts, err := s.queries.GetTeachingLoad(ctx, grade.TsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("teaching load relation not found for recorded grade", slog.Any("ts_id", grade.TsID))
			return nil, domain.ErrNotFound
		}
		log.Error("failed to fetch teaching load for recorded grade", slog.Any("error", err), slog.Any("ts_id", grade.TsID))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grade successfully recorded and linked with teaching load")

	var finalScorePtr *uint32
	if grade.Score != nil {
		finalScorePtr = mapper.Ptr(uint32(*grade.Score))
	}

	return &domain.Grade{
		ID:           grade.ID.Bytes,
		SubjectID:    ts.SubjectID.Bytes,
		StudentID:    grade.StudentID.Bytes,
		ClassID:      ts.ClassID.Bytes,
		DateOfGrade:  &grade.LessonDate.Time,
		Grade:        finalScorePtr,
		StatusCodeID: grade.StatusCodeID.Bytes,
	}, nil
}

func (s *Storage) GetGrades(
	ctx context.Context,
	teacherId uuid.UUID,
	params *domain.ListGradesRequest,
) ([]domain.Grade, error) {
	const op = "storage.postgresql.GetClassGrades"

	log := s.log.With(
		slog.String("op", op),
		slog.String("teacher_id", teacherId.String()),
		slog.String("subject_id", params.SubjectID.String()),
		slog.String("class_id", params.ClassID.String()),
	)

	log.Debug("fetching grades list from database")

	arg := sqlgen.GetGradesParams{
		SubjectID: mapper.PtrToPgUUID(params.SubjectID),
		ClassID:   mapper.PtrToPgUUID(params.ClassID),
		StudentID: mapper.PtrToPgUUID(params.StudentID),
	}

	grades, err := s.queries.GetGrades(ctx, arg)
	if err != nil {
		log.Error("failed to get grades list", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully fetched grades list", slog.Int("count", len(grades)))

	return mapper.MapSlice(
		grades,
		func(f *sqlgen.GetGradesRow) domain.Grade {
			var scorePtr *uint32
			if f.Score != nil {
				scorePtr = mapper.Ptr(uint32(*f.Score))
			}

			return domain.Grade{
				ID:           f.ID.Bytes,
				SubjectID:    f.SubjectID.Bytes,
				StudentID:    f.StudentID.Bytes,
				ClassID:      f.ClassID.Bytes,
				DateOfGrade:  &f.LessonDate.Time,
				Grade:        scorePtr,
				StatusCodeID: f.StatusCodeID.Bytes,
			}
		},
	), nil
}

func (s *Storage) UpdateGrades(
	ctx context.Context,
	params *domain.UpdateGrade,
) (*domain.Grade, error) {
	const op = "storage.postgresql.UpdateGrades"

	log := s.log.With(
		slog.String("op", op),
		slog.String("grade_id", params.GradeID.String()),
		slog.Any("grade", params.Grade),
		slog.Any("status_code", params.StatusCodeID.String()),
	)

	log.Debug("attempting to update grade")

	var score *int16
	if params.Grade != nil {
		score = mapper.Ptr(int16(*params.Grade))
	}

	arg := sqlgen.UpdateGradesParams{
		ID:           mapper.ToPgUUID(params.GradeID),
		StatusCodeID: mapper.PtrToPgUUID(params.StatusCodeID),
		Score:        score,
		Note:         params.Note,
	}

	grade, err := s.queries.UpdateGrades(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("grade row not found for update")
			return nil, domain.ErrNotFound
		}
		log.Error("failed to execute update grade query", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("grade updated successfully, fetching teaching load", slog.Any("ts_id", grade.TsID))

	ts, err := s.queries.GetTeachingLoad(ctx, grade.TsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("teaching load not found for updated grade", slog.Any("ts_id", grade.TsID))
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get teaching load", slog.Any("error", err), slog.Any("ts_id", grade.TsID))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grade fully updated and mapped")

	return &domain.Grade{
		ID:           grade.ID.Bytes,
		SubjectID:    ts.SubjectID.Bytes,
		StudentID:    grade.StudentID.Bytes,
		ClassID:      ts.ClassID.Bytes,
		DateOfGrade:  &grade.LessonDate.Time,
		Grade:        mapper.Ptr(uint32(*grade.Score)),
		StatusCodeID: grade.StatusCodeID.Bytes,
	}, nil
}

func (s *Storage) DeleteGrade(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteGrade"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id.String()),
	)

	log.Debug("attempting to delete grade")

	err := s.queries.DeleteGrade(ctx, mapper.ToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("grade row not found for deletion")
			return domain.ErrNotFound
		}
		log.Error("failed to execute delete grade query", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grade deleted successfully")
	return nil
}
