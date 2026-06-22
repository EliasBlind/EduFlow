package domain

import "github.com/google/uuid"

type Teacher struct {
	ID       *uuid.UUID `validate:"required"`
	FullName string     `validate:"required,min=2,max=20,ascii|multibyte"`
}

type ListTeachersRequest struct {
	Limit  uint32
	Offset uint32
}

type ListTeachersResponse struct {
	TotalCount uint32
	Teachers   []Teacher
}

type UpdateTeacher struct {
	ID       uuid.UUID `validate:"required"`
	FullName string    `validate:"required"`
}
