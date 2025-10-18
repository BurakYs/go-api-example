package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/BurakYs/go-api-example/app/session"
	"github.com/BurakYs/go-api-example/app/user"
	"github.com/BurakYs/go-api-example/config"
	"github.com/BurakYs/go-api-example/database"
	"github.com/BurakYs/go-api-example/middleware"
)

type Dependencies struct {
	config *config.Config
	db     *database.DB
	redis  *database.Redis
	logger *zap.Logger

	userRepository    *user.Repository
	sessionRepository *session.Repository

	userService    *user.Service
	sessionService *session.Service

	RateLimiter *middleware.RateLimiter
	RequireAuth *middleware.RequireAuth

	UserHandler *user.Handler
}

func NewDependencies(cfg *config.Config, db *database.DB, redis *database.Redis, logger *zap.Logger) *Dependencies {
	d := &Dependencies{
		config: cfg,
		db:     db,
		redis:  redis,
		logger: logger,
	}

	d.sessionRepository = session.NewRepository(d.redis)
	d.sessionService = session.NewService(d.sessionRepository, d.config.Cookie.Expiration)

	d.userRepository = user.NewRepository(d.db)
	d.userService = user.NewService(d.userRepository)
	d.UserHandler = user.NewHandler(d.userService, d.sessionService, &d.config.Cookie)

	rateLimiterCfg := middleware.RateLimiterConfig{
		Enabled:     d.config.RateLimit.Enabled,
		Window:      d.config.RateLimit.Window,
		Max:         d.config.RateLimit.Requests,
		SendHeaders: true,
		KeyFunc: func(c fiber.Ctx) string {
			return "rate_limit:" + c.IP()
		},
	}

	d.RateLimiter = middleware.NewRateLimiter(d.redis, rateLimiterCfg, d.logger)
	d.RequireAuth = middleware.NewRequireAuth(d.sessionService, d.config.Cookie.Name)

	return d
}

func (c *Dependencies) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.userRepository.CreateIndexes(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	return nil
}
