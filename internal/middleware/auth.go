package middleware

import (
	"context"
	"net/http"

	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.APIError{
				Message: "Unauthorized",
			})
			return
		}

		userID, err := db.Redis.Get(context.Background(), "session:"+sessionID).Result()
		if err == redis.Nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.APIError{
				Message: "Unauthorized",
			})
			return
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.APIError{
				Message: "Internal Server Error",
			})
			return
		}

		ctx.Set("userId", userID)
		ctx.Next()
	}
}
