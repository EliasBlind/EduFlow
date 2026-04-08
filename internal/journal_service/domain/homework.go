package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecordHomework struct {
	TeacherID       uuid.UUID `validate:"is_uuid7"`
	ClassID         uuid.UUID `validate:"is_uuid7"`
	SubjectID       uuid.UUID `validate:"is_uuid7"`
	DescriptionTask string    `validate:"required"`

	// Данные ниже могут быть nil (optional)
	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type Homework struct {
	ID              uuid.UUID `validate:"is_uuid7"`
	ClassID         uuid.UUID `validate:"is_uuid7"`
	TeacherID       uuid.UUID `validate:"is_uuid7"`
	SubjectID       uuid.UUID `validate:"is_uuid7"`
	DescriptionTask string    `validate:"required, min=1"`

	// Данные ниже могут быть nil (optional)

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type UpdateHomework struct {
	ID              uuid.UUID `validate:"is_uuid7"`
	DescriptionTask *string   `validate:"omitempty,min=1"`

	// Данные ниже могут быть nil (optional)

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type ListHomeworkRequest struct {
	ClassID   uuid.UUID `validate:"is_uuid7"`
	SubjectID uuid.UUID `validate:"is_uuid7"`

	// Данные ниже могут быть nil (optional)

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type ListHomeworkResponse struct {
	TotalCount uint32
	Homeworks  []Homework
}
