package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/httperror"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	var httpErr *httperror.HTTPError
	if errors.As(err, &httpErr) {
		return c.Status(httpErr.Code).JSON(httperror.HTTPError{
			Message: httpErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(httperror.HTTPError{
		Message: "Internal server error",
	})
}
