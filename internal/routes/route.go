package routes

import (
	"github.com/gofiber/fiber/v3"
)

func Register(router fiber.Router) {
	router.Get(
		"/",
		func(c fiber.Ctx) error {
			return c.SendString("OK")
		},
	)
}
