package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrors(err error) map[string]string {
	validationErrors := err.(validator.ValidationErrors)

	errors :=
		make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())

		switch fieldErr.Tag() {
		case "required":
			errors[field] = fmt.Sprintf("%s is required", field)

		case "min":
			errors[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldErr.Param())

		case "max":
			errors[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldErr.Param())

		case "gte":
			errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, fieldErr.Param())

		case "gt":
			errors[field] = fmt.Sprintf("%s must be greater than %s", field, fieldErr.Param())

		case "email":
			errors[field] = fmt.Sprintf("%s must be a valid email", field)

		default:
			errors[field] = fmt.Sprintf("%s is invalid", field)

		}

	}

	return errors

}
