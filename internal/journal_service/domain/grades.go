package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecordGrade struct {
	SubjectID   uuid.UUID  `validate:"is_uuid7"`
	StudentID   uuid.UUID  `validate:"is_uuid7"`
	DateOfGrade *time.Time `validate:"required"`

	// Может быть только одно поле (oneof)
	Grade        *uint32    `validate:"required_without=StatusCodeID,omitempty,min=1,max=5"`
	StatusCodeID *uuid.UUID `validate:"required_without=Grade,omitempty,is_uuid7"`

	LessonNumber uint32  `validate:"required"`
	Note         *string `validate:"omitempty,min=1"`
}

type Grade struct {
	ID          uuid.UUID `validate:"is_uuid7"`
	SubjectID   uuid.UUID `validate:"is_uuid7"`
	StudentID   uuid.UUID `validate:"is_uuid7"`
	ClassID     uuid.UUID `validate:"is_uuid7"`
	DateOfGrade *time.Time

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID uuid.UUID `validate:"is_uuid7"`

	LessonNumber *uint32
	Note         *string `validate:"omitempty,min=1"`
}

type ListGradesRequest struct {
	SubjectID *uuid.UUID `validate:"omitempty, is_uuid7"`

	// Может быть только одно поле (oneof)
	StudentID *uuid.UUID `validate:"required_without=ClassID,omitempty,is_uuid7"`
	ClassID   *uuid.UUID `validate:"required_without=StudentID,omitempty,is_uuid7"`

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type ListGradesResponse struct {
	TotalCount uint32
	Grades     []Grade
}

type UpdateGrade struct {
	GradeID uuid.UUID `validate:"is_uuid7"`

	// Может быть только одно поле (oneof)
	Grade        *uint32    `validate:"required_without=StatusCodeID,omitempty,min=1,max=5"`
	StatusCodeID *uuid.UUID `validate:"required_without=Grade,omitempty,is_uuid7"`

	Note *string `validate:"omitempty,gtfield=Start"`
}
