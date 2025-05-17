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
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIError{
				Message: "Unauthorized",
			})
			return
		}

		userID, err := db.Redis.Get(context.Background(), "session:"+sessionID).Result()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIError{
				Message: "Unauthorized",
			})
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.APIError{
				Message: "Internal Server Error",
			})
			return
		}

		c.Set("userId", userID)
		c.Next()
	}
}
