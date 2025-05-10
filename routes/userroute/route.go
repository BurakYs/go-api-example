package userroute

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/BurakYs/GoAPIExample/config"
	"github.com/BurakYs/GoAPIExample/db"
	"github.com/BurakYs/GoAPIExample/middleware"
	"github.com/BurakYs/GoAPIExample/models"
	"github.com/BurakYs/GoAPIExample/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/health-check", func(ctx *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		ctx.JSON(200, gin.H{
			"Alloc":      utils.BytesToMB(int(m.Alloc)),
			"TotalAlloc": utils.BytesToMB(int(m.TotalAlloc)),
			"Sys":        utils.BytesToMB(int(m.Sys)),
			"NumGC":      m.NumGC,
		})
	})

	router.GET("/users", func(ctx *gin.Context) {
		const pageSize = 10

		pageStr := ctx.DefaultQuery("page", "1")
		page, err := strconv.Atoi(pageStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.APIError{
				Message: "Invalid page number",
			})
			return
		}

		skip := (int(page) - 1) * pageSize

		cursor, err := db.Collections.Users.Find(
			ctx,
			bson.M{},
			options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)),
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error fetching users",
			})
			return
		}

		defer cursor.Close(ctx)

		var results []models.PublicUser

		for cursor.Next(ctx) {
			var doc models.PublicUser

			if err := cursor.Decode(&doc); err != nil {
				continue
			}

			results = append(results, doc)
		}

		if len(results) == 0 {
			ctx.JSON(http.StatusOK, []models.PublicUser{})
			return
		}

		ctx.JSON(http.StatusOK, results)
	})

	router.GET("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var result models.PublicUser
		err := db.Collections.Users.FindOne(ctx, bson.M{
			"id": id,
		}).Decode(&result)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, models.APIError{
					Message: "User not found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error fetching user",
			})
			return
		}

		ctx.JSON(http.StatusOK, result)
	})

	router.POST("/register", middleware.ValidateBody[models.RegisterUserBody](), func(ctx *gin.Context) {
		body := ctx.MustGet("body").(models.RegisterUserBody)

		userID := uuid.NewString()
		createdAt := time.Now().Format(time.RFC3339)

		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Internal Server Error",
			})
			return
		}

		_, err := db.Collections.Users.InsertOne(ctx, bson.M{
			"id":        userID,
			"username":  body.Username,
			"email":     body.Email,
			"password":  hashedPassword,
			"createdAt": createdAt,
		})

		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				ctx.JSON(http.StatusConflict, models.APIError{
					Message: "A user with the same username or e-mail already exists",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error creating user",
			})
			return
		}

		sessionID := uuid.NewString()

		if err := db.Redis.Set(ctx, "session:"+sessionID, userID, 15*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Could not create session",
			})
			return
		}

		ctx.SetCookie("session_id", sessionID, 900, "/", config.AppConfig.Domain, true, true)

		ctx.JSON(http.StatusCreated, gin.H{
			"id":        userID,
			"username":  body.Username,
			"email":     body.Email,
			"createdAt": createdAt,
		})
	})

	router.POST("/login", middleware.ValidateBody[models.LoginUserBody](), func(ctx *gin.Context) {
		body := ctx.MustGet("body").(models.LoginUserBody)

		var result models.User
		err := db.Collections.Users.FindOne(ctx, bson.M{
			"email": body.Email,
		}).Decode(&result)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, models.APIError{
					Message: "User not found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error fetching user",
			})
			return
		}

		compareErr := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(body.Password))
		if compareErr != nil {
			ctx.JSON(http.StatusForbidden, models.APIError{
				Message: "Invalid e-mail or password",
			})
			return
		}

		sessionID := uuid.NewString()

		if err := db.Redis.Set(ctx, "session:"+sessionID, result.ID, 15*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Could not create session",
			})
			return
		}

		ctx.SetCookie("session_id", sessionID, 900, "/", config.AppConfig.Domain, true, true)

		ctx.JSON(http.StatusOK, gin.H{
			"id":        result.ID,
			"username":  result.Username,
			"email":     result.Email,
			"createdAt": result.CreatedAt,
		})
	})

	router.POST("/logout", middleware.AuthRequired(), func(ctx *gin.Context) {
		sessionID, _ := ctx.Cookie("session_id")

		if err := db.Redis.Del(ctx, "session:"+sessionID).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error deleting session",
			})
			return
		}

		ctx.SetCookie("session_id", "", -1, "/", config.AppConfig.Domain, true, true)
		ctx.Status(http.StatusNoContent)
	})

	router.DELETE("/delete-account", middleware.AuthRequired(), func(ctx *gin.Context) {
		userID, _ := ctx.Get("userId")

		err := db.Collections.Users.FindOneAndDelete(ctx, bson.M{
			"id": userID,
		}).Err()

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, models.APIError{
					Message: "User not found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error deleting user",
			})
			return
		}

		ctx.JSON(http.StatusOK, models.APIError{
			Message: "User deleted successfully",
		})
	})
}
