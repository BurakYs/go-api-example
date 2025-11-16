package middleware

import (
	"errors"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/httperror"
	"github.com/BurakYs/go-api-example/util/validator"
)

type ValidationFailure struct {
	Location string `json:"location"`
	Field    string `json:"field"`
	Message  string `json:"error"`
}

var val = validator.New()

const (
	bindingBody   = "body"
	bindingQuery  = "query"
	bindingParams = "params"
	bindingForm   = "form"
)

func ValidateBody[T any](c fiber.Ctx) (*T, error) {
	return validate[T](c, bindingBody)
}

func ValidateQuery[T any](c fiber.Ctx) (*T, error) {
	return validate[T](c, bindingQuery)
}

func ValidateParams[T any](c fiber.Ctx) (*T, error) {
	return validate[T](c, bindingParams)
}

func ValidateForm[T any](c fiber.Ctx) (*T, error) {
	return validate[T](c, bindingForm)
}

type normalizable interface {
	Normalize()
}

func validate[T any](c fiber.Ctx, location string) (*T, error) {
	data := new(T)
	var err error

	switch location {
	case bindingBody:
		err = c.Bind().Body(data)
	case bindingQuery:
		err = c.Bind().Query(data)
	case bindingParams:
		err = c.Bind().URI(data)
	case bindingForm:
		err = c.Bind().Form(data)
	}

	if err != nil {
		return data, httperror.New(fiber.StatusBadRequest, "Invalid parameters provided").WithExtra("validationFailures", []ValidationFailure{})
	}

	if n, ok := any(data).(normalizable); ok {
		n.Normalize()
	}

	err = val.Struct(data)
	if err != nil {
		var ve govalidator.ValidationErrors
		if errors.As(err, &ve) {
			failures := make([]ValidationFailure, len(ve))
			for i, fieldError := range ve {
				failures[i] = ValidationFailure{
					Location: location,
					Field:    fieldError.Field(),
					Message:  validator.GetErrorMessage(fieldError),
				}
			}

			return data, httperror.New(fiber.StatusBadRequest, "Invalid parameters provided").WithExtra("validationFailures", failures)
		}

		return data, httperror.New(fiber.StatusBadRequest, "Invalid parameters provided").WithExtra("validationFailures", []ValidationFailure{})
	}

	return data, nil
}
