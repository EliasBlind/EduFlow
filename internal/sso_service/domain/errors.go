package domain

import (
	"fmt"
)

type AppError struct {
	Key     string
	Message string
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Key, e.Message)
}

func newErr(key, msg string) AppError {
	return AppError{Key: key, Message: msg}
}

var (
	ErrInvalidData = newErr("common.invalid_data", "invalid data")
	ErrInternal    = newErr("common.internal_error", "internal server error")

	ErrWeakPassword      = newErr("auth.weak_password", "the password is too weak")
	ErrUserAlreadyExists = newErr("auth.user_already_exists", "user with this email already exists")
	ErrUserNotFound      = newErr("auth.user_not_found", "user not found")
	ErrUnauthenticated   = newErr("auth.unauthenticated", "the token is not valid or expired")
	ErrAccessDenied      = newErr("auth.access_denied", "access denied")

	ErrInvalidCode  = newErr("verify.invalid_code", "invalid verification code")
	ErrCodeExpired  = newErr("verify.code_expired", "verification code expired")
	ErrCodeNotFound = newErr("verify.code_not_found", "verification code not found")

	ErrStorageUnavailable = newErr("infra.storage_unavailable", "storage service unavailable")
)
