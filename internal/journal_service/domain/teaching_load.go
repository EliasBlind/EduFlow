package domain

import "github.com/google/uuid"

type CreateTeachingLoad struct {
	TeacherID uuid.UUID `validate:"is_uuid7"`
	SubjectID uuid.UUID `validate:"is_uuid7"`
	ClassID   uuid.UUID `validate:"is_uuid7"`
}

type TeachingLoad struct {
	ID        uuid.UUID `validate:"is_uuid7"`
	TeacherID uuid.UUID `validate:"is_uuid7"`
	SubjectID uuid.UUID `validate:"is_uuid7"`
	ClassID   uuid.UUID `validate:"is_uuid7"`
}

type UpdateTeachingLoad struct {
	ID        uuid.UUID `validate:"is_uuid7"`
	TeacherID uuid.UUID `validate:"is_uuid7"`
}

type ListTeachingLoadRequest struct {
	// Может быть только одно поле (oneof)
	TeacherID *uuid.UUID `validate:"omitempty,is_uuid7"`
	ClassID   *uuid.UUID `validate:"omitempty,is_uuid7"`
}

type ListTeachingLoadResponse struct {
	TotalCount    uint32
	TeachingLoads []TeachingLoad
}
