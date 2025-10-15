package user

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/internal/app/session"
	"github.com/BurakYs/go-api-example/internal/config"
	"github.com/BurakYs/go-api-example/internal/httperror"
	"github.com/BurakYs/go-api-example/internal/middleware"
	"github.com/BurakYs/go-api-example/internal/util/rctx"
)

type Handler struct {
	svc        *Service
	sessionSvc *session.Service
	cookieCfg  *config.CookieConfig
}

func NewHandler(svc *Service, sessionSvc *session.Service, cookieCfg *config.CookieConfig) *Handler {
	return &Handler{
		svc:        svc,
		sessionSvc: sessionSvc,
		cookieCfg:  cookieCfg,
	}
}

func (h *Handler) Register(c fiber.Ctx) error {
	body, ok := middleware.ValidateBody[RegistrationBody](c)
	if !ok {
		return nil
	}

	user, err := h.svc.Register(c, body.Name, body.Email, body.Password)
	if err != nil {
		return err
	}

	sessionID, err := h.sessionSvc.Create(c, user.ID)
	if err != nil {
		return err
	}

	h.setSessionCookie(c, sessionID)
	return c.JSON(NewAuthResponse(user))
}

func (h *Handler) Login(c fiber.Ctx) error {
	body, ok := middleware.ValidateBody[LoginBody](c)
	if !ok {
		return nil
	}

	user, err := h.svc.Login(c, body.Email, body.Password)
	if err != nil {
		return err
	}

	sessionID, err := h.sessionSvc.Create(c, user.ID)
	if err != nil {
		return err
	}

	h.setSessionCookie(c, sessionID)
	return c.JSON(NewAuthResponse(user))
}

func (h *Handler) Me(c fiber.Ctx) error {
	userID := rctx.GetUserID(c)

	user, err := h.svc.GetByID(c, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return httperror.New(fiber.StatusNotFound, "User not found")
		}
		return err
	}

	return c.JSON(NewAuthResponse(user))
}

func (h *Handler) Logout(c fiber.Ctx) error {
	_ = h.sessionSvc.Delete(c, rctx.GetSessionID(c))
	h.clearSessionCookie(c)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) setSessionCookie(c fiber.Ctx, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     h.cookieCfg.Name,
		Value:    value,
		MaxAge:   int(h.cookieCfg.Expiration.Seconds()),
		HTTPOnly: true,
		Secure:   h.cookieCfg.Secure,
		SameSite: h.cookieCfg.SameSite,
	})
}

func (h *Handler) clearSessionCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     h.cookieCfg.Name,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.cookieCfg.Secure,
		SameSite: h.cookieCfg.SameSite,
	})
}
