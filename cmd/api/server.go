package main

import (
	"log"

	"github.com/BurakYs/go-api-example/internal/middleware"
	"github.com/BurakYs/go-api-example/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type structValidator struct {
	validator *validator.Validate
}

func (v *structValidator) Validate(i any) error {
	return v.validator.Struct(i)
}

type server struct {
	app *fiber.App
}

func newServer() *server {
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

	return &server{
		app: app,
	}
}

func (s *server) setupRoutes(deps *dependencies) {
	s.app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	authGroup := s.app.Group("/auth")
	authGroup.Post(
		"/login",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(5))),
		deps.AuthHandler.Login,
	)

	authGroup.Post(
		"/logout",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(5))),
		middleware.AuthRequired(deps.Redis),
		deps.AuthHandler.Logout,
	)

	authGroup.Post(
		"/register",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(3))),
		deps.AuthHandler.Register,
	)

	authGroup.Delete(
		"/delete-account",
		middleware.AuthRequired(deps.Redis),
		deps.AuthHandler.DeleteAccount,
	)

	userGroup := s.app.Group("/users")
	userGroup.Get(
		"/",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(50))),
		deps.UserHandler.GetAll,
	)

	userGroup.Get(
		"/:id",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(50))),
		deps.UserHandler.GetByID,
	)

	s.app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(models.APIError{
			Message: "Page not found",
		})
	})
}

func (s *server) listen(port string) error {
	log.Println("Listening on http://localhost:" + port)

	return s.app.Listen(":"+port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}
