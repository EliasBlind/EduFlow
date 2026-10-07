package usercalimas

import (
	"fmt"

	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenParser struct {
	secret string `validate:"required"`
}

func NewToken(key string) *TokenParser {
	return &TokenParser{secret: key}
}

type jwtClaims struct {
	ID                   string `json:"Id"`
	Login                string `json:"Login"`
	Role                 string `json:"Role"`
	jwt.RegisteredClaims        // exp, iat и т.д.
}

func (tp *TokenParser) ParseToken(tokenStr string) (*UserClaims, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tp.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	return &UserClaims{
		ID:    id,
		Login: claims.Login,
		Role:  roles.GetRole(claims.Role),
	}, nil
}
