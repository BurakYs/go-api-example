package main

import (
	"context"

	"github.com/gofiber/fiber/v3"
	loggermi "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"

	"github.com/BurakYs/go-api-example/httperror"
	"github.com/BurakYs/go-api-example/middleware"
)

type Server struct {
	app *fiber.App
}

func NewServer(logger *zap.Logger) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler:  middleware.ErrorHandler,
		CaseSensitive: true,
	})

	app.Use(
		recover.New(),
		loggermi.New(loggermi.Config{
			LoggerFunc: func(c fiber.Ctx, data *loggermi.Data, _ *loggermi.Config) error {
				logger.Info("HTTP Request",
					zap.String("ip", c.IP()),
					zap.Int("status", c.Response().StatusCode()),
					zap.Duration("latency", data.Stop.Sub(data.Start)),
					zap.String("method", c.Method()),
					zap.String("path", c.OriginalURL()),
					zap.NamedError("error", data.ChainErr),
				)
				return nil
			},
		}),
	)

	return &Server{
		app: app,
	}
}

func (s *Server) SetupRoutes(deps *Dependencies) {
	s.app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	auth := s.app.Group("/auth")
	auth.Post("/register", deps.RateLimiter.Middleware(), deps.UserHandler.Register)
	auth.Post("/login", deps.RateLimiter.Middleware(), deps.UserHandler.Login)
	auth.Post("/logout", deps.RequireAuth.Middleware(), deps.UserHandler.Logout)

	users := s.app.Group("/users")
	users.Get("/me", deps.RequireAuth.Middleware(), deps.UserHandler.Me)

	s.app.Use(func(c fiber.Ctx) error {
		return httperror.New(fiber.StatusNotFound, "Page not found")
	})
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(":"+port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
