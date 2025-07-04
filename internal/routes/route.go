package routes

import (
	"github.com/gofiber/fiber/v3"
)

func Register(router fiber.Router) {
	router.Add(
		[]string{"GET", "HEAD"},
		"/",
		func(c fiber.Ctx) error {
			return c.SendString("OK")
		},
	)
}
