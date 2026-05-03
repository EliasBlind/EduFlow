package validator

import (
	"time"

	"github.com/EliasBlind/EduFlow/pkg/roles"
	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("year_gte_now", func(fl validator.FieldLevel) bool {
		year, ok := fl.Field().Interface().(uint32)
		if !ok {
			return false
		}
		currentYear := uint32(time.Now().Year())
		return year >= currentYear
	})

	v.RegisterValidation("is_role", func(fl validator.FieldLevel) bool {
		roleStr, ok := fl.Field().Interface().(string)
		if !ok {
			r, ok := fl.Field().Interface().(roles.Role)
			if !ok {
				return false
			}
			return r.IsRole()
		}

		return roles.Role(roleStr).IsRole()
	})
	return v
}
