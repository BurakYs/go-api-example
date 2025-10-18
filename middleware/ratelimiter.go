package middleware

import (
	_ "embed"
	"strconv"
	"time"

	"github.com/BurakYs/go-api-example/database"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/BurakYs/go-api-example/httperror"
)

var (
	//go:embed ratelimiter_sliding.lua
	slidingLimiterScriptContent string

	//go:embed ratelimiter_fixed.lua
	fixedLimiterScriptContent string

	slidingLimiterScript = redis.NewScript(slidingLimiterScriptContent)
	fixedLimiterScript   = redis.NewScript(fixedLimiterScriptContent)
)

type RateLimiterConfig struct {
	Enabled     bool
	Window      time.Duration
	Max         int
	SendHeaders bool
	KeyFunc     func(c fiber.Ctx) string
}

type RateLimiter struct {
	redis      *database.Redis
	defaultCfg RateLimiterConfig
	logger     *zap.Logger
}

type RateLimiterBuilder struct {
	rateLimiter *RateLimiter
	config      RateLimiterConfig
	sliding     bool
}

func NewRateLimiter(redis *database.Redis, defaultCfg RateLimiterConfig, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		redis:      redis,
		defaultCfg: defaultCfg,
		logger:     logger,
	}
}

func (rl *RateLimiter) Fixed() *RateLimiterBuilder {
	return &RateLimiterBuilder{
		rateLimiter: rl,
		config:      rl.defaultCfg,
		sliding:     false,
	}
}

func (rl *RateLimiter) Sliding() *RateLimiterBuilder {
	return &RateLimiterBuilder{
		rateLimiter: rl,
		config:      rl.defaultCfg,
		sliding:     true,
	}
}

func (rl *RateLimiter) Middleware() fiber.Handler {
	return rl.Fixed().Middleware()
}

func (b *RateLimiterBuilder) WithEnabled(enabled bool) *RateLimiterBuilder {
	b.config.Enabled = enabled
	return b
}

func (b *RateLimiterBuilder) WithWindow(window time.Duration) *RateLimiterBuilder {
	b.config.Window = window
	return b
}

func (b *RateLimiterBuilder) WithMax(maxRequests int) *RateLimiterBuilder {
	b.config.Max = maxRequests
	return b
}

func (b *RateLimiterBuilder) WithSendHeaders(sendHeaders bool) *RateLimiterBuilder {
	b.config.SendHeaders = sendHeaders
	return b
}

func (b *RateLimiterBuilder) WithKeyFunc(keyFunc func(c fiber.Ctx) string) *RateLimiterBuilder {
	b.config.KeyFunc = keyFunc
	return b
}

func (b *RateLimiterBuilder) Middleware() fiber.Handler {
	var script *redis.Script
	if b.sliding {
		script = slidingLimiterScript
	} else {
		script = fixedLimiterScript
	}

	if b.config.Window < time.Second {
		b.rateLimiter.logger.Warn("RateLimiter window is less than 1 second, setting to 1 second")
		b.config.Window = time.Second
	}

	windowSecs := int(b.config.Window.Seconds())

	return func(c fiber.Ctx) error {
		if !b.config.Enabled {
			return c.Next()
		}

		key := []string{b.config.KeyFunc(c)}
		now := time.Now().UTC().UnixMilli()

		result, err := b.rateLimiter.redis.EvalScript(c.Context(), script, key, b.config.Max, windowSecs, now)
		if err != nil {
			return err
		}

		remaining := result[0]
		resetSeconds := result[1]
		allowed := remaining >= 0

		if remaining < 0 {
			remaining = 0
		}

		if resetSeconds < 0 {
			resetSeconds = 0
		}

		if b.config.SendHeaders {
			if !allowed {
				c.Set("Retry-After", strconv.FormatInt(resetSeconds, 10))
			} else {
				c.Set("X-RateLimit-Limit", strconv.Itoa(b.config.Max))
				c.Set("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
				c.Set("X-RateLimit-Reset", strconv.FormatInt(resetSeconds, 10))
			}
		}

		if !allowed {
			return httperror.New(fiber.StatusTooManyRequests, "Too many requests")
		}

		return c.Next()
	}
}
