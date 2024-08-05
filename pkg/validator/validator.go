package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidDateFormat = errors.New("invalid date format")
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{v: validate}
}

func (vl *Validator) Validate(b interface{}) error {
	err := vl.v.Struct(b)
	if err != nil {
		var errors []string
		for _, item := range err.(validator.ValidationErrors) {
			ve := vl.getValidationErr(item)
			errors = append(errors, ve.Error())
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func (vl *Validator) ValidateDate(dobStr string) (time.Time, error) {
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}
	return dob, nil
}

func (vl *Validator) getValidationErr(err validator.FieldError) error {
	var (
		field = err.Field()
		param = err.Param()
	)
	switch err.Tag() {
	case "required":
		return fmt.Errorf("field %s is required", field)
	case "min":
		return fmt.Errorf("field %s must be at least %s characters long", field, param)
	case "max":
		return fmt.Errorf("field %s must be at most %s characters long", field, param)
	case "email":
		return fmt.Errorf("field %s must be a valid email address", field)
	case "url":
		return fmt.Errorf("field %s must be a valid URL", field)
	case "uuid":
		return fmt.Errorf("field %s must be a valid UUID", field)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}
