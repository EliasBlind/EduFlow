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

func (e *AppError) WithKey(newKey string) *AppError {
	return &AppError{
		Key:     newKey,
		Message: e.Message,
	}
}

/*
Подменные ключи которые могут заменить стандартные:

учитель не выдал право на выставление пропусков старосте
common.forbidden.did_not_give_permission -> common.forbidden
*/
var (
	ErrInvalidData  = newErr("common.invalid_data", "invalid data")
	ErrInternal     = newErr("common.internal_error", "internal server error")
	ErrUnauthorized = newErr("common.unauthorized", "user claims not found in context")
	ErrForbidden    = newErr("common.forbidden", "permission denied")

	ErrNotFound      = newErr("common.not_found", "resource not found")
	ErrAlreadyExists = newErr("common.already_exists", "resource already exists")
)
