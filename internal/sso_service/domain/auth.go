package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Login    string `validate:"required,min=5,max=20,alphanumunicode"`
	Password string `validate:"required,max=40"`
	AppId    int
}

type VerifyRequest struct {
	Email string `validate:"required,email"`
	Code  string `validate:"required,min=6,max=6"`
}

type LoginRequest struct {
	Login    string `validate:"required,min=5,max=20,alphanumunicode"`
	Password string `validate:"required,max=40"`
	AppId    int
}

type TokenPair struct {
	AccessToken  string `validate:"required"`
	RefreshToken string `validate:"required"`
}

type User struct {
	Id           uuid.UUID `validate:"required,uuid4"`
	Email        string    `validate:"required,email"`
	Login        string    `validate:"required,min=5,max=20,alphanumunicode"`
	PasswordHash []byte    `validate:"required"`
	Role         *string   `validate:"is_role"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Id    uuid.UUID `validate:"required,uuid4"`
	Login string    `validate:"required,min=5,max=20,alphanumunicode"`
	Role  string    `validate:"is_role"`
}

type RefreshSession struct {
	ID        uuid.UUID `validate:"required,uuid4"`
	UserID    uuid.UUID `validate:"required,uuid4"`
	AppID     int
	TokenID   uuid.UUID `validate:"required,uuid4"`
	ExpiresAt time.Time `validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `validate:"required"`
	AppId        int
}
