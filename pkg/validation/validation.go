package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(v interface{}) error {
	err := validator.New().Struct(v)
	if err != nil {
		var errors []string
		for _, item := range err.(validator.ValidationErrors) {
			ve := getValidationErr(item)
			errors = append(errors, ve.Error())
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func getValidationErr(err validator.FieldError) error {
	switch err.Tag() {
	case "required":
		return fmt.Errorf("field %s is required", err.Field())
	default:
		return fmt.Errorf("field %s is invalid", err.Field())
	}
}
