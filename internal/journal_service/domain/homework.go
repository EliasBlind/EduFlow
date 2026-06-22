package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecordHomework struct {
	TeacherID       uuid.UUID `validate:"required"`
	ClassID         uuid.UUID `validate:"required"`
	SubjectID       uuid.UUID `validate:"required"`
	DescriptionTask string    `validate:"required"`

	// Данные ниже могут быть nil (optional)
	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type Homework struct {
	ID              uuid.UUID `validate:"required"`
	ClassID         uuid.UUID `validate:"required"`
	TeacherID       uuid.UUID `validate:"required"`
	SubjectID       uuid.UUID `validate:"required"`
	DescriptionTask string    `validate:"required, min=1"`

	// Данные ниже могут быть nil (optional)

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type UpdateHomework struct {
	ID              uuid.UUID `validate:"required"`
	DescriptionTask *string   `validate:"omitempty,min=1"`

	// Данные ниже могут быть nil (optional)

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type ListHomeworkRequest struct {
	ClassID   uuid.UUID `validate:"required"`
	SubjectID uuid.UUID `validate:"required"`

	// Данные ниже могут быть nil (optional)

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type ListHomeworkResponse struct {
	TotalCount uint32
	Homeworks  []Homework
}
