package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func New() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("is_uuid7", func(fl validator.FieldLevel) bool {
		id, ok := fl.Field().Interface().(uuid.UUID)
		if !ok {
			return false
		}
		return id.Version() == 7
	})

	v.RegisterValidation("year_gte_now", func(fl validator.FieldLevel) bool {
		year, ok := fl.Field().Interface().(uint32)
		if !ok {
			return false
		}
		currentYear := uint32(time.Now().Year())
		return year >= currentYear
	})

	return v
}
