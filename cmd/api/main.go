package main

import (
	"log"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/BurakYs/GoAPIExample/internal/routes"
	"github.com/BurakYs/GoAPIExample/internal/routes/authroute"
	"github.com/BurakYs/GoAPIExample/internal/routes/userroute"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
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
		db.DisconnectRedis()
	}()

	app := fiber.New(fiber.Config{
		ErrorHandler:  middleware.ErrorHandler(),
		CaseSensitive: true,
		StructValidator: &structValidator{
			validator: validator.New(),
		},
	})

	app.Use(
		recover.New(),
		logger.New(logger.Config{
			Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "UTC",
		}),
	)

	router := app.Group("")
	routes.Register(router)

	usersCollection := db.GetCollection("users")

	userController := userroute.NewUserController(usersCollection)
	userroute.Register(router, userController)

	authGroup := router.Group("/auth")
	authController := authroute.NewAuthController(usersCollection)
	authroute.Register(authGroup, authController)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(models.APIError{
			Message: "Not Found",
		})
	})

	log.Println("Listening on http://localhost:" + config.App.Port)

	err := app.Listen(":"+config.App.Port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})

	if err != nil {
		log.Fatalln(err)
	}
}
