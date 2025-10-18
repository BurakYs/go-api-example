package user

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/httperror"
)

type User struct {
	ID        string    `json:"id"        bson:"_id"`
	Name      string    `json:"name"      bson:"name"`
	Email     string    `json:"email"     bson:"email"`
	Password  string    `json:"-"         bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

var (
	ErrAlreadyExists      = httperror.New(fiber.StatusConflict, "This email is already registered")
	ErrInvalidCredentials = httperror.New(fiber.StatusUnauthorized, "Invalid email or password")
	ErrNotFound           = httperror.New(fiber.StatusNotFound, "User not found")
)
