package domain

import (
	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `validate:"required"`
	Login string `validate:"required,min=2,max=20,ascii|multibyte"`
	Email string `validate:"required,email"`
	Role  roles.Role
}
