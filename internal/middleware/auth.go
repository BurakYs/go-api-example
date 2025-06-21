package middleware

import (
	"context"
	"errors"

	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{
				Message: "Unauthorized",
			})
		}

		userID, err := db.Redis.Get(context.Background(), "session:"+sessionID).Result()
		if errors.Is(err, redis.Nil) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIError{
				Message: "Unauthorized",
			})
		}

		if err != nil {
			panic(err)
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
