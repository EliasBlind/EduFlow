package domain

import (
	"fmt"

	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenParser struct {
	secret string `validate:"required"`
}

func NewToken(token string) *TokenParser {
	return &TokenParser{secret: token}
}

func (tp *TokenParser) ParseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(tp.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	idStr, okId := claims["Id"].(string)
	login, okLogin := claims["Login"].(string)
	role, okRole := claims["Role"].(string)
	if !okId || !okRole || !okLogin {
		return nil, fmt.Errorf("missing fields in token")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	return &UserClaims{
		ID:    id,
		Login: login,
		Role:  roles.GetRole(role),
	}, nil
}
