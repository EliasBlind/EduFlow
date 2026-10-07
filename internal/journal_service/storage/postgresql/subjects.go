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

func (s *Storage) CreateSubject(
	ctx context.Context,
	subjectName string,
) (*domain.Subject, error) {
	const op = "storage.postgresql.CreateSubject"

	log := s.log.With(
		"op", op,
		"subject_name", subjectName,
	)

	res, err := s.queries.CreateSubject(ctx, subjectName)
	if err != nil {
		log.Error("failed to create subject", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Subject{
		ID:       res.ID.Bytes,
		FullName: res.FullName,
	}, nil
}

func (s *Storage) UpdateSubject(
	ctx context.Context,
	params *domain.UpdateSubject,
) (*domain.Subject, error) {
	const op = "storage.postgresql.UpdateSubject"

	log := s.log.With(
		"op", op,
		"id", params.ID.String(),
		"full_name", params.FullName,
	)

	arg := sqlgen.UpdateSubjectParams{
		ID:       mapper.ToPgUUID(params.ID),
		FullName: params.FullName,
	}

	res, err := s.queries.UpdateSubject(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("subject not found for update", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update subject", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.Subject{
		ID:       res.ID.Bytes,
		FullName: res.FullName,
	}, nil
}

func (s *Storage) DeleteSubject(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteSubject"

	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	err := s.queries.DeleteSubject(ctx, mapper.ToPgUUID(id))
	if err != nil {
		log.Error("failed to delete subject", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *Storage) ListSubjects(ctx context.Context) ([]domain.Subject, error) {
	const op = "storage.postgresql.ListSubjects"

	log := s.log.With(
		"op", op,
	)

	subjects, err := s.queries.ListSubjects(ctx)
	if err != nil {
		log.Error("failed to delete subject", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return pkgmapper.MapSliceRef(
		subjects,
		func(f *sqlgen.Subject) domain.Subject {
			return domain.Subject{
				ID:       f.ID.Bytes,
				FullName: f.FullName,
			}
		},
	), nil
}
