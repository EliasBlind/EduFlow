package domain

import "github.com/google/uuid"

type StatusCode struct {
	ID       uuid.UUID `validate:"required"`
	FullName string    `validate:"required"`
}

type UpdateStatusCode struct {
	ID       uuid.UUID `validate:"required"`
	FullName string    `validate:"required"`
}

type ListStatusCode struct {
	TotalCount  uint32
	StatusCodes []StatusCode
}
