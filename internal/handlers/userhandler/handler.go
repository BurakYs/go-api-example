package userhandler

import (
	"errors"

	"github.com/BurakYs/go-api-example/internal/middleware"
	"github.com/BurakYs/go-api-example/internal/models"
	"github.com/BurakYs/go-api-example/internal/services/userservice"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service *userservice.UserService
}

func NewUserHandler(service *userservice.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetAll(c fiber.Ctx) error {
	query, ok := middleware.ValidateQuery[models.GetAllUsersQuery](c)
	if !ok {
		return nil
	}

	users, err := h.service.GetAll(c, query.Page)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) GetByID(c fiber.Ctx) error {
	params, ok := middleware.ValidateParams[models.GetUserByIDParams](c)
	if !ok {
		return nil
	}

	user, err := h.service.GetByID(c, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		case errors.Is(err, models.ErrValidation):
			return c.Status(fiber.StatusBadRequest).JSON(models.APIError{
				Message: "Invalid user ID format",
			})
		default:
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
