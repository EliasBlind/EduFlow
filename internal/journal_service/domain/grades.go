package domain

import (
	"time"

	"github.com/google/uuid"
)

type RecordGrade struct {
	TsID        uuid.UUID
	StudentID   uuid.UUID
	DateOfGrade *time.Time `validate:"required"`

	// Может быть только одно поле (oneof)
	Grade        *uint32    `validate:"required_without=StatusCodeID,omitempty,min=1,max=5"`
	StatusCodeID *uuid.UUID `validate:"required_without=Grade,omitempty,uuid4"`

	LessonNumber uint32  `validate:"required"`
	Note         *string `validate:"omitempty,min=1"`
}

type Grade struct {
	ID          uuid.UUID `validate:"required"`
	SubjectID   uuid.UUID `validate:"required"`
	StudentID   uuid.UUID `validate:"required"`
	ClassID     uuid.UUID `validate:"required"`
	DateOfGrade *time.Time

	// Может быть только одно поле (oneof)
	Grade        *uint32
	StatusCodeID uuid.UUID `validate:"required"`

	LessonNumber *uint32
	Note         *string `validate:"omitempty,min=1"`
}

type ListGradesRequest struct {
	SubjectID *uuid.UUID `validate:"omitempty"`

	// Может быть только одно поле (oneof)
	StudentID *uuid.UUID `validate:"required_without=ClassID"`
	ClassID   *uuid.UUID `validate:"required_without=StudentID,omitempty,uuid4"`

	Start *time.Time `validate:"omitempty"`
	End   *time.Time `validate:"omitempty,gtfield=Start"`
}

type ListGradesResponse struct {
	TotalCount uint32
	Grades     []Grade
}

type UpdateGrade struct {
	GradeID uuid.UUID `validate:"required"`

	// Может быть только одно поле (oneof)
	Grade        *uint32    `validate:"required_without=StatusCodeID,omitempty,min=1,max=5"`
	StatusCodeID *uuid.UUID `validate:"required_without=Grade,omitempty,uuid4"`

	Note *string `validate:"omitempty,gtfield=Start"`
}
