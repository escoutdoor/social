package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(v interface{}) error {
	err := validator.New().Struct(v)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]
		switch fieldErr.Tag() {
		case "required":
			return fmt.Errorf("field %s is required", fieldErr.Field())
		default:
			return fmt.Errorf("field %s is invalid", fieldErr.Field())
		}
	}

	return nil
}
