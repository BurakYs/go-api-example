package middleware

import (
	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Skip: func(c *gin.Context) bool {
			ip := c.ClientIP()

			isRelease := config.App.GinMode == gin.ReleaseMode
			isLocal := ip == "127.0.0.1" || ip == "::1"
			return isRelease && isLocal
		},
	})
}
