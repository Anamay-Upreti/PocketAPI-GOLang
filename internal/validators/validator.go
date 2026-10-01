package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func FormatValidationError(err error) string {

	validationErrors := err.(validator.ValidationErrors)

	firstError := validationErrors[0]

	switch firstError.Tag() {

	case "required":
		return fmt.Sprintf("%s is required", firstError.Field())

	case "email":
		return "Invalid email address"

	case "min":
		return fmt.Sprintf(
			"%s must contain at least %s characters",
			firstError.Field(),
			firstError.Param(),
		)

	case "max":
		return fmt.Sprintf(
			"%s cannot exceed %s characters",
			firstError.Field(),
			firstError.Param(),
		)
	}

	return firstError.Error()
}