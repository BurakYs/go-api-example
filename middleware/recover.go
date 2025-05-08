package middleware

import (
	"net/http"

	"github.com/BurakYs/GoAPIExample/models"
	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
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
