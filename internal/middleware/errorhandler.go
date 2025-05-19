package middleware

import (
	"log"

	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Message: "Internal server error",
		})
	}
}
