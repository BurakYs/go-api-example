package userroute

import (
	"errors"
	"time"

	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserController struct {
	collection *mongo.Collection
}

func NewUserController(collection *mongo.Collection) *UserController {
	return &UserController{
		collection: collection,
	}
}

func (uc *UserController) GetAllUsers(c fiber.Ctx) error {
	const pageSize = 10

	query := middleware.GetQuery[models.GetAllUsersQuery](c)
	skip := (query.Page - 1) * pageSize

	cursor, err := uc.collection.Find(
		c,
		bson.M{},
		options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)),
	)

	if err != nil {
		return err
	}

	defer cursor.Close(c)

	var results []models.PublicUser
	for cursor.Next(c) {
		var doc models.User

		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		results = append(results, models.PublicUser{
			ID:        doc.ID.Hex(),
			Username:  doc.Username,
			CreatedAt: doc.CreatedAt.Time().Format(time.RFC3339),
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON([]models.PublicUser{})
	}

	return c.Status(fiber.StatusOK).JSON(results)
}

func (uc *UserController) GetUserByID(c fiber.Ctx) error {
	params := middleware.GetParams[models.GetUserByIDParams](c)

	objectID, err := bson.ObjectIDFromHex(params.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIError{
			Message: "Invalid user ID",
		})
	}

	var result models.User
	err = uc.collection.FindOne(c, bson.M{
		"_id": objectID,
	}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.PublicUser{
		ID:        result.ID.Hex(),
		Username:  result.Username,
		CreatedAt: result.CreatedAt.Time().Format(time.RFC3339),
	})
}

func (uc *UserController) DeleteAccount(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	err = uc.collection.FindOneAndDelete(c, bson.M{
		"_id": objectID,
	}).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	c.ClearCookie("session_id")
	return c.SendStatus(fiber.StatusNoContent)
}
