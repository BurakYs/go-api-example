package main

import (
	"log"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/BurakYs/GoAPIExample/internal/routes"
	"github.com/BurakYs/GoAPIExample/internal/routes/userroute"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type structValidator struct {
	validator *validator.Validate
}

func (v *structValidator) Validate(i any) error {
	return v.validator.Struct(i)
}

func main() {
	config.LoadEnv()

	db.SetupMongo()
	db.SetupRedis()

	defer func() {
		db.DisconnectMongo()
	}()

	app := fiber.New(fiber.Config{
		ErrorHandler:  middleware.ErrorHandler(),
		CaseSensitive: true,
		StructValidator: &structValidator{
			validator: validator.New(),
		},
	})

	app.Use(recover.New(), middleware.Logger())

	router := app.Group("")
	routes.Register(&router)
	userroute.Register(&router)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(models.APIError{
			Message: "Not Found",
		})
	})

	err := app.Listen(":"+config.App.Port, fiber.ListenConfig{
		DisableStartupMessage: config.App.Mode == config.ModeRelease,
	})

	if err != nil {
		log.Fatalln(err)
	}
}
