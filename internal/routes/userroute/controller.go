package userroute

import (
	"context"
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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (*UserController) GetAllUsers(c fiber.Ctx) error {
	const pageSize = 10

	query := c.Locals(middleware.BindingLocationQuery).(models.GetAllUsersQuery)
	skip := (query.Page - 1) * pageSize

	cursor, err := db.Collections.Users.Find(
		context.TODO(),
		bson.M{},
		options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)),
	)

	if err != nil {
		return err
	}

	defer cursor.Close(context.TODO())

	var results []models.PublicUser
	for cursor.Next(context.TODO()) {
		var doc models.PublicUser

		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		results = append(results, doc)
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON([]models.PublicUser{})
	}

	return c.Status(fiber.StatusOK).JSON(results)
}

func (*UserController) GetUserByID(c fiber.Ctx) error {
	params := c.Locals(middleware.BindingLocationParams).(models.GetUserByIDParams)

	var result models.PublicUser
	err := db.Collections.Users.FindOne(context.TODO(), bson.M{
		"id": params.ID,
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (*UserController) Register(c fiber.Ctx) error {
	body := c.Locals(middleware.BindingLocationBody).(models.RegisterUserBody)

	userID := uuid.NewString()
	createdAt := time.Now().Format(time.RFC3339)

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}

	_, err := db.Collections.Users.InsertOne(context.TODO(), bson.M{
		"id":        userID,
		"username":  body.Username,
		"email":     body.Email,
		"password":  hashedPassword,
		"createdAt": createdAt,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).JSON(models.APIError{
				Message: "A user with the same username or e-mail already exists",
			})
		}

		return err
	}

	sessionID := uuid.NewString()

	if err := db.Redis.Set(context.TODO(), "session:"+sessionID, userID, 15*time.Minute).Err(); err != nil {
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        userID,
		"username":  body.Username,
		"email":     body.Email,
		"createdAt": createdAt,
	})
}

func (*UserController) Login(c fiber.Ctx) error {
	body := c.Locals(middleware.BindingLocationBody).(models.LoginUserBody)

	var result models.User
	err := db.Collections.Users.FindOne(context.TODO(), bson.M{
		"email": body.Email,
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(body.Password))
	if compareErr != nil {
		return c.Status(fiber.StatusForbidden).JSON(models.APIError{
			Message: "Invalid e-mail or password",
		})
	}

	sessionID := uuid.NewString()

	if err := db.Redis.Set(context.TODO(), "session:"+sessionID, result.ID, 15*time.Minute).Err(); err != nil {
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        result.ID,
		"username":  result.Username,
		"email":     result.Email,
		"createdAt": result.CreatedAt,
	})
}

func (*UserController) Logout(c fiber.Ctx) error {
	sessionID := c.Cookies("session_id")

	if err := db.Redis.Del(context.TODO(), "session:"+sessionID).Err(); err != nil {
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

func (*UserController) DeleteAccount(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	err := db.Collections.Users.FindOneAndDelete(context.TODO(), bson.M{
		"id": userID,
	}).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
