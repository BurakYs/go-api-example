package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	GinMode     string `env:"GIN_MODE" envDefault:"debug"`
	Port        string `env:"PORT,required"`
	Domain      string `env:"DOMAIN,required"`
	MongoDBName string `env:"MONGODB_DBNAME,required"`
	MongoURI    string `env:"MONGODB_URI,required"`
}

var App Config

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Invalid .env file")
	}

	if err := env.Parse(&App); err != nil {
		log.Fatalln("Failed to parse the .env file:", err)
	}
}
