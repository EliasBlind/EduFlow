package domain

import "github.com/google/uuid"

type Student struct {
	ID       *uuid.UUID `validate:"required"`
	ClassID  *uuid.UUID `validate:"required"`
	FullName string     `validate:"required,min=2,max=20,ascii|multibyte"`
}

type ListStudentsRequest struct {
	ClassID uuid.UUID
	Limit   uint32
	Offset  uint32
}

type ListStudentsResponse struct {
	TotalCount uint32
	Students   []Student
}

type UpdateStudent struct {
	ID       uuid.UUID
	ClassID  *uuid.UUID
	FullName *string `validate:"omitempty,required"`
}
