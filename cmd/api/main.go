package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/database"
	"github.com/BurakYs/GoAPIExample/internal/handlers/authhandler"
	"github.com/BurakYs/GoAPIExample/internal/handlers/userhandler"
	"github.com/BurakYs/GoAPIExample/internal/repository/authrepository"
	"github.com/BurakYs/GoAPIExample/internal/repository/userrepository"
	"github.com/BurakYs/GoAPIExample/internal/services/authservice"
	"github.com/BurakYs/GoAPIExample/internal/services/userservice"
	"github.com/joho/godotenv"
)

type dependencies struct {
	DB    *database.MongoDB
	Redis *database.Redis

	AuthHandler *authhandler.AuthHandler
	UserHandler *userhandler.UserHandler
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(".env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load env config:", err)
	}

	fmt.Printf("%+v\n", cfg)

	deps, cleanup, err := initDeps(cfg)
	if err != nil {
		log.Fatalln("Failed to initialize app", err)
	}

	defer cleanup()

	app := newServer()
	app.setupRoutes(deps)
	err = app.listen(cfg.App.Port)

	if err != nil {
		log.Fatalln(err)
	}
}

func initDeps(cfg *config.Config) (*dependencies, func(), error) {
	db, err := database.NewMongoDB(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return nil, nil, err
	}

	redis, err := database.NewRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
	if err != nil {
		return nil, func() { db.Disconnect(context.Background()) }, err
	}

	authRepository := authrepository.NewAuthRepository(db.Database(), redis)
	authService := authservice.NewAuthService(authRepository, cfg.App.Domain)

	userRepository := userrepository.NewUserRepository(db.Database())
	userService := userservice.NewUserService(userRepository)

	deps := &dependencies{
		DB:          db,
		Redis:       redis,
		AuthHandler: authhandler.NewAuthHandler(authService),
		UserHandler: userhandler.NewUserHandler(userService),
	}

	cleanup := func() {
		redis.Close()
		db.Disconnect(context.Background())
	}

	return deps, cleanup, nil
}
