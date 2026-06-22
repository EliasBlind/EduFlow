package domain

import "github.com/google/uuid"

type CreateTeachingLoad struct {
	TeacherID uuid.UUID `validate:"required"`
	SubjectID uuid.UUID `validate:"required"`
	ClassID   uuid.UUID `validate:"required"`
}

type TeachingLoad struct {
	ID        uuid.UUID `validate:"required"`
	TeacherID uuid.UUID `validate:"required"`
	SubjectID uuid.UUID `validate:"required"`
	ClassID   uuid.UUID `validate:"required"`
}

type UpdateTeachingLoad struct {
	ID        uuid.UUID `validate:"required"`
	TeacherID uuid.UUID `validate:"required"`
}

type ListTeachingLoadRequest struct {
	// Может быть только одно поле (oneof)
	TeacherID *uuid.UUID `validate:"omitempty,required"`
	ClassID   *uuid.UUID `validate:"omitempty,required"`
}

type ListTeachingLoadResponse struct {
	TotalCount    uint32
	TeachingLoads []TeachingLoad
}
