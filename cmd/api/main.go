package main

import (
	"net/http"

	"github.com/BurakYs/GoAPIExample/config"
	"github.com/BurakYs/GoAPIExample/db"
	"github.com/BurakYs/GoAPIExample/models"
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

	router.Use(recovery())
	userroute.RegisterRoutes(router)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.APIError{
			Message: "Page not found",
		})
	})

	router.Run(":" + config.AppConfig.Port)
}

func recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.APIError{
					Message: "Internal server error",
				})
			}
		}()

		ctx.Next()
	}
}
