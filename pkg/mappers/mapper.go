package pkgmapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/EliasBlind/EduFlow/internal/journal_service/domain"
)

func MapSliceRef[F any, T any](items []F, mapper func(*F) T) []T {
	if items == nil {
		return nil
	}

	res := make([]T, len(items))

	for i := range items {
		res[i] = mapper(&items[i])
	}
	return res
}

func MapSlice[F any, T any](items []F, mapper func(F) T) []T {
	if items == nil {
		return nil
	}

	res := make([]T, len(items))

	for i, v := range items {
		res[i] = mapper(v)
	}
	return res
}

func ToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func PtrToPgUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func PtrToPgInt2(data *uint32) pgtype.Int2 {
	if data == nil {
		return pgtype.Int2{Valid: false}
	}
	return pgtype.Int2{
		Int16: int16(*data),
		Valid: true,
	}
}

func PtrToPgText(data *string) pgtype.Text {
	if data == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *data,
		Valid:  true,
	}
}

func Ptr[T any](v T) *T {
	return &v
}

func ToTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}

func ToDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func FromTimestamp(t pgtype.Date) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func ToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func ToUUID(data *string) (*uuid.UUID, error) {
	if data == nil {
		return nil, nil
	}

	buffer, err := uuid.Parse(*data)
	if err != nil {
		return nil, domain.ErrInvalidData
	}
	return &buffer, nil
}

func ToPtrUUID(id uuid.UUID) *uuid.UUID {
	return &id
}
