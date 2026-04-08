package domain

import "github.com/google/uuid"

type Teacher struct {
	ID       uuid.UUID `validate:"is_uuid7"`
	FullName string
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
	ID       uuid.UUID `validate:"is_uuid7"`
	FullName string    `validate:"required"`
}
