package validator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidDateFormat = errors.New("invalid date format")
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

func ValidateDate(dobStr string) (time.Time, error) {
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return dob, nil
}

func getValidationErr(err validator.FieldError) error {
	switch err.Tag() {
	case "required":
		return fmt.Errorf("field %s is required", err.Field())
	default:
		return fmt.Errorf("field %s is invalid", err.Field())
	}
}
