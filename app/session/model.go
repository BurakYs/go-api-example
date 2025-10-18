package session

import (
	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/httperror"
)

var ErrNotFound = httperror.New(fiber.StatusNotFound, "Session not found")
