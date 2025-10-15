package middleware

import (
	"errors"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/internal/httperror"
	"github.com/BurakYs/go-api-example/internal/util/validator"
)

var val = validator.New()

const (
	bindingBody   = "body"
	bindingQuery  = "query"
	bindingParams = "params"
	bindingForm   = "form"
)

func ValidateBody[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingBody)
}

func ValidateQuery[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingQuery)
}

func ValidateParams[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingParams)
}

func ValidateForm[T any](c fiber.Ctx) (*T, bool) {
	return validate[T](c, bindingForm)
}

type normalizable interface {
	Normalize()
}

func validate[T any](c fiber.Ctx, location string) (*T, bool) {
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
		_ = c.Status(fiber.StatusBadRequest).JSON(httperror.HTTPError{Message: "Invalid request data"})
		return data, false
	}

	if n, ok := any(data).(normalizable); ok {
		n.Normalize()
	}

	err = val.Struct(data)
	if err != nil {
		var ve govalidator.ValidationErrors
		if errors.As(err, &ve) {
			_ = c.Status(fiber.StatusBadRequest).JSON(validator.ToResponse(ve, location))
		} else {
			_ = c.Status(fiber.StatusBadRequest).JSON(httperror.HTTPError{Message: "Invalid request data"})
		}

		return data, false
	}

	return data, true
}
