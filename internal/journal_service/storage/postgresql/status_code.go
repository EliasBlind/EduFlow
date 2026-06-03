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

func (s *Storage) CreateStatusCode(
	ctx context.Context,
	statusCodeName string,
) (*domain.StatusCode, error) {
	const op = "storage.postgresql.CreateStatusCode"
	log := s.log.With(
		"op", op,
		"status_code_name", statusCodeName,
	)

	res, err := s.queries.CreateStatusCode(ctx, statusCodeName)
	if err != nil {
		log.Error("failed to create status code", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.StatusCode{
		ID:       res.ID.Bytes,
		FullName: res.Abbreviation,
	}, nil
}

func (s *Storage) UpdateStatusCode(
	ctx context.Context,
	params *domain.UpdateStatusCode,
) (*domain.StatusCode, error) {
	const op = "storage.postgresql.UpdateStatusCode"
	log := s.log.With(
		"op", op,
		"id", params.ID.String(),
		"full_name", params.FullName,
	)

	arg := sqlgen.UpdateStatusCodeParams{
		ID:           mapper.ToPgUUID(params.ID),
		Abbreviation: params.FullName,
	}

	res, err := s.queries.UpdateStatusCode(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("status code not found for update", "id", params.ID)
			return nil, domain.ErrNotFound
		}
		log.Error("failed to update status code", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.StatusCode{
		ID:       res.ID.Bytes,
		FullName: res.Abbreviation,
	}, nil
}

func (s *Storage) ListStatusCode(
	ctx context.Context,
) ([]domain.StatusCode, error) {
	const op = "storage.postgresql.ListStatusCode"
	log := s.log.With(
		"op", op,
	)

	statusCodes, err := s.queries.ListStatusCode(ctx)
	if err != nil {
		log.Error("failed to list status codes", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mapper.MapSliceRef(
		statusCodes,
		func(f *sqlgen.StatusCode) domain.StatusCode {
			return domain.StatusCode{
				ID:       f.ID.Bytes,
				FullName: f.Abbreviation,
			}
		},
	), nil
}

func (s *Storage) DeleteStatusCode(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgresql.DeleteStatusCode"
	log := s.log.With(
		"op", op,
		"id", id.String(),
	)

	err := s.queries.DeleteStatusCode(ctx, mapper.ToPgUUID(id))
	if err != nil {
		log.Error("failed to delete status code", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
