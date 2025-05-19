package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Mode        string `env:"MODE" envDefault:"debug"` // debug or release
	Port        string `env:"PORT,required"`
	Domain      string `env:"DOMAIN,required"`
	MongoDBName string `env:"MONGODB_DBNAME,required"`
	MongoURI    string `env:"MONGODB_URI,required"`
}

var App Config

const (
	ModeDebug   = "debug"
	ModeRelease = "release"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Invalid .env file")
	}

	if err := env.Parse(&App); err != nil {
		log.Fatalln("Failed to parse the .env file:", err)
	}
}
