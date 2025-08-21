package middleware

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/internal/models"
)

const (
	bindingLocationBody   = "body"
	bindingLocationQuery  = "query"
	bindingLocationParams = "params"
	bindingLocationForm   = "form"
)

func ValidateBody[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingLocationBody)
}

func ValidateQuery[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingLocationQuery)
}

func ValidateParams[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingLocationParams)
}

func ValidateForm[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingLocationForm)
}

type transformable interface {
	Transform()
}

func validate[T any](c fiber.Ctx, location string) (*T, bool) {
	data := new(T)
	var err error

	switch location {
	case bindingLocationBody:
		err = c.Bind().Body(data)
	case bindingLocationQuery:
		err = c.Bind().Query(data)
	case bindingLocationParams:
		err = c.Bind().URI(data)
	case bindingLocationForm:
		err = c.Bind().Form(data)
	}

	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(formatValidationError(err, location, data))
		return data, false
	}

	if t, ok := any(data).(transformable); ok {
		t.Transform()
	}

	return data, true
}

func formatValidationError(err error, location string, data any) models.ValidationError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return models.ValidationError{
			Message:            "Invalid parameters provided",
			ValidationFailures: []models.ValidationFailure{},
		}
	}

	failures := make([]models.ValidationFailure, 0, len(ve))

	for _, fe := range ve {
		field := getFieldName(fe.StructField(), data)
		failures = append(failures, formatFieldError(fe, location, field))
	}

	return models.ValidationError{
		Message:            "Invalid parameters provided",
		ValidationFailures: failures,
	}
}

func getFieldName(structField string, data any) string {
	t := reflect.TypeOf(data).Elem()

	if f, ok := t.FieldByName(structField); ok {
		tags := []string{"json", "query", "uri", "form", "header", "cookie", "cbor", "respHeader", "xml"}

		for _, tag := range tags {
			if tagValue := f.Tag.Get(tag); tagValue != "" && tagValue != "-" {
				return strings.Split(tagValue, ",")[0]
			}
		}
	}

	return structField
}

func formatFieldError(fe validator.FieldError, location, field string) models.ValidationFailure {
	var msg string

	switch fe.Tag() {
	case "required":
		msg = "This field is required"
	case "email":
		msg = "This field must be a valid email address"
	case "uuid":
		msg = "This field must be a valid UUID"
	case "min":
		switch fe.Kind() {
		case reflect.String:
			msg = fmt.Sprintf("This field must be at least %s characters long", fe.Param())
		case reflect.Slice, reflect.Array:
			msg = fmt.Sprintf("This field must contain at least %s items", fe.Param())
		default:
			msg = fmt.Sprintf("The value must be at least %s", fe.Param())
		}
	case "max":
		switch fe.Kind() {
		case reflect.String:
			msg = fmt.Sprintf("This field must be at most %s characters long", fe.Param())
		case reflect.Slice, reflect.Array:
			msg = fmt.Sprintf("This field must contain at most %s items", fe.Param())
		default:
			msg = fmt.Sprintf("The value must be at most %s", fe.Param())
		}
	default:
		msg = fmt.Sprintf("This field is invalid for tag: %s", fe.Tag())
	}

	return models.ValidationFailure{
		Location: location,
		Field:    field,
		Message:  msg,
	}
}
