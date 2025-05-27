package routes

import (
	"runtime"

	"github.com/BurakYs/GoAPIExample/internal/utils"
	"github.com/gofiber/fiber/v3"
)

func Register(router fiber.Router) {
	router.Get(
		"/",
		func(c fiber.Ctx) error {
			return c.SendString("OK")
		},
	)

	router.Get(
		"/health-check",
		func(c fiber.Ctx) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			return c.JSON(fiber.Map{
				"Alloc":      utils.BytesToMB(int(m.Alloc)),
				"TotalAlloc": utils.BytesToMB(int(m.TotalAlloc)),
				"Sys":        utils.BytesToMB(int(m.Sys)),
				"NumGC":      m.NumGC,
			})
		},
	)
}
