package middleware

import (
	"log"

	"github.com/BurakYs/go-api-example/internal/models"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		log.Printf("%s: %v", c.Path(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Message: "Internal server error",
		})
	}
}
