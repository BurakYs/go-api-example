package middleware

import (
	"log"
	"time"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		start := time.Now()

		ctx.Next()

		if config.App.GinMode == gin.ReleaseMode && (clientIP == "127.0.0.1" || clientIP == "::1") {
			return
		}

		latency := time.Since(start)
		status := ctx.Writer.Status()

		log.Printf("[GIN] | %3d | %13v | %15s | %-7s | %#v\n",
			status,
			latency,
			clientIP,
			ctx.Request.Method,
			ctx.Request.URL.Path,
		)
	}
}
