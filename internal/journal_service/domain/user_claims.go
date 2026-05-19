package domain

import (
	"context"

	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/google/uuid"
)

type ctxKey struct{}

type UserClaims struct {
	ID    uuid.UUID  `validate:"required"`
	Login string     `validate:"required"`
	Role  roles.Role `validate:"required"`
}

func ContextWithClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, ctxKey{}, claims)
}

func GetUserClaims(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(ctxKey{}).(*UserClaims)
	if !ok || claims == nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}
