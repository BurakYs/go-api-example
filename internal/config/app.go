package config

import "github.com/caarlos0/env/v11"

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type AppConfig struct {
	GoEnv  string `env:"GO_ENV" envDefault:"debug"` // debug or release
	Port   string `env:"PORT" envDefault:"8080"`
	Domain string `env:"DOMAIN,required"`
}

type DatabaseConfig struct {
	Name string `env:"MONGODB_DBNAME,required"`
	URI  string `env:"MONGODB_URI,required"`
}

type RedisConfig struct {
	Host string `env:"REDIS_HOST,required"`
	Port string `env:"REDIS_PORT,required"`
	DB   int    `env:"REDIS_DB,required"`
}

func Load() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
