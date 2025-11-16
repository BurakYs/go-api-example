package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
)

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

func GetErrorMessage(fieldError govalidator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "This field must be a valid email address"
	case "uuid":
		return "This field must be a valid UUID"
	case "alpha_space":
		return "This field can only contain alphabetic and space characters"
	case "min":
		switch fieldError.Kind() {
		case reflect.String:
			return fmt.Sprintf("This field must be at least %s characters long", fieldError.Param())
		case reflect.Slice, reflect.Array:
			return fmt.Sprintf("This field must contain at least %s items", fieldError.Param())
		default:
			return fmt.Sprintf("The value must be at least %s", fieldError.Param())
		}
	case "max":
		switch fieldError.Kind() {
		case reflect.String:
			return fmt.Sprintf("This field must be at most %s characters long", fieldError.Param())
		case reflect.Slice, reflect.Array:
			return fmt.Sprintf("This field must contain at most %s items", fieldError.Param())
		default:
			return fmt.Sprintf("The value must be at most %s", fieldError.Param())
		}
	default:
		return fmt.Sprintf("This field is invalid for: %s", fieldError.Tag())
	}
}
