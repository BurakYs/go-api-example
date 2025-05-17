package userroute

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/BurakYs/GoAPIExample/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/health-check", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(200, gin.H{
			"Alloc":      utils.BytesToMB(int(m.Alloc)),
			"TotalAlloc": utils.BytesToMB(int(m.TotalAlloc)),
			"Sys":        utils.BytesToMB(int(m.Sys)),
			"NumGC":      m.NumGC,
		})
	})

	router.GET("/users", func(c *gin.Context) {
		const pageSize = 10

		pageStr := c.DefaultQuery("page", "1")
		page, err := strconv.Atoi(pageStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIError{
				Message: "Invalid page number",
			})
			return
		}

		skip := (int(page) - 1) * pageSize

		cursor, err := db.Collections.Users.Find(
			c,
			bson.M{},
			options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error fetching users",
			})
			return
		}

		defer cursor.Close(c)

		var results []models.PublicUser

		for cursor.Next(c) {
			var doc models.PublicUser

			if err := cursor.Decode(&doc); err != nil {
				continue
			}

			results = append(results, doc)
		}

		if len(results) == 0 {
			c.JSON(http.StatusOK, []models.PublicUser{})
			return
		}

		c.JSON(http.StatusOK, results)
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var result models.PublicUser
		err := db.Collections.Users.FindOne(c, bson.M{
			"id": id,
		}).Decode(&result)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, models.APIError{
					Message: "User not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error fetching user",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	router.POST("/register", middleware.ValidateBody[models.RegisterUserBody](), func(c *gin.Context) {
		body := c.MustGet("body").(models.RegisterUserBody)

		userID := uuid.NewString()
		createdAt := time.Now().Format(time.RFC3339)

		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Internal Server Error",
			})
			return
		}

		_, err := db.Collections.Users.InsertOne(c, bson.M{
			"id":        userID,
			"username":  body.Username,
			"email":     body.Email,
			"password":  hashedPassword,
			"createdAt": createdAt,
		})

		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, models.APIError{
					Message: "A user with the same username or e-mail already exists",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error creating user",
			})
			return
		}

		sessionID := uuid.NewString()

		if err := db.Redis.Set(c, "session:"+sessionID, userID, 15*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Could not create session",
			})
			return
		}

		c.SetCookie("session_id", sessionID, 900, "/", config.App.Domain, true, true)

		c.JSON(http.StatusCreated, gin.H{
			"id":        userID,
			"username":  body.Username,
			"email":     body.Email,
			"createdAt": createdAt,
		})
	})

	router.POST("/login", middleware.ValidateBody[models.LoginUserBody](), func(c *gin.Context) {
		body := c.MustGet("body").(models.LoginUserBody)

		var result models.User
		err := db.Collections.Users.FindOne(c, bson.M{
			"email": body.Email,
		}).Decode(&result)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, models.APIError{
					Message: "User not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error fetching user",
			})
			return
		}

		compareErr := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(body.Password))
		if compareErr != nil {
			c.JSON(http.StatusForbidden, models.APIError{
				Message: "Invalid e-mail or password",
			})
			return
		}

		sessionID := uuid.NewString()

		if err := db.Redis.Set(c, "session:"+sessionID, result.ID, 15*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Could not create session",
			})
			return
		}

		c.SetCookie("session_id", sessionID, 900, "/", config.App.Domain, true, true)

		c.JSON(http.StatusOK, gin.H{
			"id":        result.ID,
			"username":  result.Username,
			"email":     result.Email,
			"createdAt": result.CreatedAt,
		})
	})

	router.POST("/logout", middleware.AuthRequired(), func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")

		if err := db.Redis.Del(c, "session:"+sessionID).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error deleting session",
			})
			return
		}

		c.SetCookie("session_id", "", -1, "/", config.App.Domain, true, true)
		c.Status(http.StatusNoContent)
	})

	router.DELETE("/delete-account", middleware.AuthRequired(), func(c *gin.Context) {
		userID, _ := c.Get("userId")

		err := db.Collections.Users.FindOneAndDelete(c, bson.M{
			"id": userID,
		}).Err()

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, models.APIError{
					Message: "User not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, models.APIError{
				Message: "Error deleting user",
			})
			return
		}

		c.JSON(http.StatusOK, models.APIError{
			Message: "User deleted successfully",
		})
	})
}
