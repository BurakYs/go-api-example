package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func SetupRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if err := Redis.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("Failed to connect to Redis:", err)
	}
}

func DisconnectRedis() {
	Redis.Close()
}
