package main

import (
	"net/http"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/BurakYs/GoAPIExample/internal/routes/userroute"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db.SetupMongo()
	db.SetupRedis()

	defer func() {
		db.DisconnectMongo()
	}()

	gin.SetMode(config.App.GinMode)

	router := gin.New()

	router.Use(middleware.Logger(), middleware.Recovery())
	userroute.RegisterRoutes(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.APIError{
			Message: "Page not found",
		})
	})

	router.Run(":" + config.App.Port)
}
