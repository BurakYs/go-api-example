package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App       AppConfig
	Cookie    CookieConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Port     string `env:"PORT"      envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

type CookieConfig struct {
	Name       string        `env:"COOKIE_NAME,required"`
	Expiration time.Duration `env:"COOKIE_EXPIRATION"         envDefault:"24h"`
	Domain     string        `env:"COOKIE_DOMAIN,required"`
	Secure     bool          `env:"COOKIE_SECURE,required"`
	SameSite   string        `env:"COOKIE_SAME_SITE,required"`
}

type DatabaseConfig struct {
	Name string `env:"MONGODB_DBNAME,required"`
	URI  string `env:"MONGODB_URI,required"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     string `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB,required"`
}

type RateLimitConfig struct {
	Enabled  bool          `env:"RATE_LIMIT_ENABLED"  envDefault:"true"`
	Requests int           `env:"RATE_LIMIT_REQUESTS" envDefault:"50"`
	Window   time.Duration `env:"RATE_LIMIT_WINDOW"   envDefault:"60s"`
}

func Load() (*Config, error) {
	var config Config

	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
