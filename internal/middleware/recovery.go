package middleware

import (
	"net/http"

	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, models.APIError{
					Message: "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
