package middleware

import (
	"log"
	"time"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/gofiber/fiber/v3"
)

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		ip := c.IP()

		isRelease := config.App.GoEnv == config.EnvRelease
		isLocal := ip == "127.0.0.1" || ip == "::1"
		skip := isRelease && isLocal
		if skip {
			return c.Next()
		}

		c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		log.Printf(
			"%3d | %13v | %15s | %-7s | %#v\n",
			status,
			latency,
			ip,
			method,
			path,
		)

		return nil
	}

}
