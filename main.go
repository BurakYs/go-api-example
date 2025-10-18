package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/BurakYs/go-api-example/config"
	"github.com/BurakYs/go-api-example/database"
	loggerpkg "github.com/BurakYs/go-api-example/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(".env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	logger := loggerpkg.New(cfg.App.LogLevel)
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to sync logger:", err)
		}
	}()

	db, err := database.NewDB(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}

	defer func() {
		err := db.Disconnect(context.Background())
		if err != nil {
			logger.Error("Failed to disconnect from DB", zap.Error(err))
		}
	}()

	redis, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	defer func() {
		err := redis.Close()
		if err != nil {
			logger.Error("Failed to close Redis connection", zap.Error(err))
		}
	}()

	deps := NewDependencies(cfg, db, redis, logger)
	err = deps.Init()
	if err != nil {
		logger.Fatal("Failed to initialize dependencies", zap.Error(err))
	}

	server := NewServer(logger)
	server.SetupRoutes(deps)

	go func(port string) {
		logger.Info("Listening on http://localhost:" + port)

		err = server.Listen(port)
		if err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}(cfg.App.Port)

	gracefulShutdown(server, logger)
}

func gracefulShutdown(server *Server, logger *zap.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logger.Fatal("Failed to gracefully shutdown server", zap.Error(err))
	}
}
