package config

import (
	"time"

	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type LimiterOption func(*limiter.Config)

func BaseLimiter(opts ...LimiterOption) limiter.Config {
	cfg := limiter.Config{
		Max:        250,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.APIError{
				Message: "Too many requests",
			})
		},
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

func LimiterWithMax(max int) LimiterOption {
	return func(cfg *limiter.Config) {
		cfg.Max = max
	}
}

func LimiterWithExpiration(d time.Duration) LimiterOption {
	return func(cfg *limiter.Config) {
		cfg.Expiration = d
	}
}
