package authroute

import (
	"errors"
	"time"

	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/db"
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	collection *mongo.Collection
}

func NewAuthController(collection *mongo.Collection) *AuthController {
	return &AuthController{
		collection: collection,
	}
}

func (uc *AuthController) Register(c fiber.Ctx) error {
	body := middleware.GetBody[models.RegisterUserBody](c)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	existingUserFilter := bson.M{
		"$or": []bson.M{
			{"username": body.Username},
			{"email": body.Email},
		},
	}

	err = uc.collection.FindOne(c, existingUserFilter).Err()
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(models.APIError{
			Message: "A user with the same username or e-mail already exists",
		})
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	userID := bson.NewObjectID()
	createdAt := time.Now().UTC()
	_, err = uc.collection.InsertOne(c, models.User{
		ID:        userID,
		Username:  body.Username,
		Email:     body.Email,
		Password:  string(hashedPassword),
		CreatedAt: bson.DateTime(createdAt.UnixMilli()),
	})

	if err != nil {
		return err
	}

	sessionID := uuid.NewString()

	if err := db.Redis.Set(c, "session:"+sessionID, userID.Hex(), 15*time.Minute).Err(); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   900,
		Path:     "/",
		Domain:   config.App.Domain,
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(models.PublicUser{
		ID:        userID.Hex(),
		Username:  body.Username,
		CreatedAt: createdAt.Format(time.RFC3339),
	})
}

func (uc *AuthController) Login(c fiber.Ctx) error {
	body := middleware.GetBody[models.LoginUserBody](c)

	var result models.User
	err := uc.collection.FindOne(c, bson.M{"email": body.Email}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(models.APIError{
			Message: "Invalid e-mail or password",
		})
	}

	userID := result.ID.Hex()
	sessionID := uuid.NewString()

	if err := db.Redis.Set(c, "session:"+sessionID, userID, 15*time.Minute).Err(); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   900,
		Path:     "/",
		Domain:   config.App.Domain,
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(models.PublicUser{
		ID:        userID,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.Time().Format(time.RFC3339),
	})
}

func (uc *AuthController) Logout(c fiber.Ctx) error {
	sessionID := c.Cookies("session_id")

	if err := db.Redis.Del(c, "session:"+sessionID).Err(); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   config.App.Domain,
		Secure:   true,
		HTTPOnly: true,
	})

	return c.SendStatus(fiber.StatusNoContent)
}
