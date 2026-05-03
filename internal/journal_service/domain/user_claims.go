package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ctxKey struct{}

type UserClaims struct {
	ID   uuid.UUID `validate:"is_uuid7"`
	Role string    `validate:"required"`
}

func ContextWithClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, ctxKey{}, claims)
}

func GetUserClaims(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(ctxKey{}).(*UserClaims)
	if !ok {
		return nil, errors.New("user claims not found in context")
	}
	return claims, nil
}
