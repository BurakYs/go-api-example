package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
)

type ValidationFailure struct {
	Location string `json:"location"`
	Field    string `json:"field"`
	Message  string `json:"error"`
}

type ValidationError struct {
	Message  string              `json:"error"`
	Failures []ValidationFailure `json:"validationFailures"`
}

func New() *govalidator.Validate {
	validate := govalidator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tags := []string{"json", "query", "uri", "form"}

		for _, tag := range tags {
			tagValue := field.Tag.Get(tag)
			if tagValue == "" {
				continue
			}

			tagName := strings.SplitN(tagValue, ",", 2)[0]
			if tagName != "-" {
				return tagName
			}
		}

		return field.Name
	})

	_ = validate.RegisterValidation("alpha_space", func(fl govalidator.FieldLevel) bool {
		alphaSpaceRegex := regexp.MustCompile("^[a-zA-Z ]+$")
		return alphaSpaceRegex.MatchString(fl.Field().String())
	})

	return validate
}

func ToResponse(ve govalidator.ValidationErrors, location string) ValidationError {
	failures := make([]ValidationFailure, len(ve))
	for i, fieldError := range ve {
		failures[i] = ValidationFailure{
			Location: location,
			Field:    fieldError.Field(),
			Message:  getErrorMessage(fieldError),
		}
	}

	return ValidationError{
		Message:  "Invalid parameters provided",
		Failures: failures,
	}
}

func getErrorMessage(fieldError govalidator.FieldError) string {
	var msg string

	switch fieldError.Tag() {
	case "required":
		msg = "This field is required"
	case "email":
		msg = "This field must be a valid email address"
	case "uuid":
		msg = "This field must be a valid UUID"
	case "alpha_space":
		msg = "This field can only contain alphabetic and space characters"
	case "min":
		switch fieldError.Kind() {
		case reflect.String:
			msg = fmt.Sprintf("This field must be at least %s characters long", fieldError.Param())
		case reflect.Slice, reflect.Array:
			msg = fmt.Sprintf("This field must contain at least %s items", fieldError.Param())
		default:
			msg = fmt.Sprintf("The value must be at least %s", fieldError.Param())
		}
	case "max":
		switch fieldError.Kind() {
		case reflect.String:
			msg = fmt.Sprintf("This field must be at most %s characters long", fieldError.Param())
		case reflect.Slice, reflect.Array:
			msg = fmt.Sprintf("This field must contain at most %s items", fieldError.Param())
		default:
			msg = fmt.Sprintf("The value must be at most %s", fieldError.Param())
		}
	default:
		msg = fmt.Sprintf("This field is invalid for: %s", fieldError.Tag())
	}

	return msg
}
