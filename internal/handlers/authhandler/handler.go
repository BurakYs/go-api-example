package authhandler

import (
	"errors"
	"time"

	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/BurakYs/GoAPIExample/internal/services/authservice"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service *authservice.AuthService
}

func NewAuthHandler(service *authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	body, ok := middleware.ValidateBody[models.RegisterUserBody](c)
	if !ok {
		return nil
	}

	user, sessionID, err := h.service.Register(c, body)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return c.Status(fiber.StatusConflict).JSON(models.APIError{
				Message: "A user with the same username or e-mail already exists",
			})
		}

		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   int(15 * time.Minute / time.Second),
		Path:     "/",
		Domain:   h.service.Domain(),
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	body, ok := middleware.ValidateBody[models.LoginUserBody](c)
	if !ok {
		return nil
	}

	user, sessionID, err := h.service.Login(c, body)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		case errors.Is(err, models.ErrForbidden):
			return c.Status(fiber.StatusForbidden).JSON(models.APIError{
				Message: "Invalid e-mail or password",
			})
		default:
			return err
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   int(15 * time.Minute / time.Second),
		Path:     "/",
		Domain:   h.service.Domain(),
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if err := h.service.Logout(c, sessionID); err != nil {
		return err
	}

	c.ClearCookie("session_id")
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) DeleteAccount(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	err := h.service.DeleteAccount(c, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIError{
				Message: "User not found",
			})
		}

		return err
	}

	c.ClearCookie("session_id")
	return c.SendStatus(fiber.StatusNoContent)
}
