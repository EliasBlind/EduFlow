package domain

import "github.com/google/uuid"

type Student struct {
	ID       uuid.UUID  `validate:"is_uuid7"`
	ClassID  *uuid.UUID `validate:"omitempty,is_uuid7"`
	FullName string     `validate:"required"`
}

type CreateStudent struct {
	ClassID  *uuid.UUID `validate:"omitempty,is_uuid7"`
	FullName string     `validate:"required"`
}

type ListStudentsRequest struct {
	ClassID uuid.UUID `validate:"is_uuid7"`
	Limit   uint32
	Offset  uint32
}

type ListStudentsResponse struct {
	TotalCount uint32
	Students   []Student
}

type UpdateStudent struct {
	ID       uuid.UUID  `validate:"is_uuid7"`
	ClassID  *uuid.UUID `validate:"omitempty,is_uuid7"`
	FullName *string    `validate:"omitempty,required"`
}
