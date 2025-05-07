package main

import (
	"github.com/BurakYs/GoAPIExample/config"
	"github.com/BurakYs/GoAPIExample/db"
	"github.com/BurakYs/GoAPIExample/routes/userroute"

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
