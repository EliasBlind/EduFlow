package domain

import "github.com/google/uuid"

type Subject struct {
	ID       uuid.UUID `validate:"required"`
	FullName string    `validate:"required"`
}

type UpdateSubject struct {
	ID       uuid.UUID `validate:"required"`
	FullName string    `validate:"required"`
}
