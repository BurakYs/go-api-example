package main

import (
	"go-api-example/config"
	"go-api-example/db"
	"go-api-example/routes/userroute"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db.SetupMongo()
	db.SetupRedis()

	defer func() {
		db.DisconnectMongo()
	}()

	gin.SetMode(config.AppConfig.GinMode)

	router := gin.Default()

	userroute.RegisterRoutes(router)

	router.Run(":" + config.AppConfig.Port)
}
