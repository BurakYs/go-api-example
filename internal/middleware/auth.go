package middleware

import (
	"context"
	"errors"

	"github.com/BurakYs/GoAPIExample/internal/database"
	"github.com/BurakYs/GoAPIExample/internal/models"

	"github.com/gofiber/fiber/v3"
	goredis "github.com/redis/go-redis/v9"
)

func AuthRequired(redis *database.Redis) fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{
				Message: "Unauthorized",
			})
		}

		userID, err := redis.Get(context.Background(), "session:"+sessionID)
		if errors.Is(err, goredis.Nil) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{
				Message: "Unauthorized",
			})
		}

		if err != nil {
			return err
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
