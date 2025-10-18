package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/BurakYs/go-api-example/app/session"
	"github.com/BurakYs/go-api-example/httperror"
	"github.com/BurakYs/go-api-example/util/rctx"
)

type RequireAuth struct {
	service    *session.Service
	cookieName string
}

func NewRequireAuth(service *session.Service, cookieName string) *RequireAuth {
	return &RequireAuth{
		service:    service,
		cookieName: cookieName,
	}
}

func (m *RequireAuth) Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		sid := c.Cookies(m.cookieName)
		if sid == "" {
			return httperror.New(fiber.StatusUnauthorized, "Unauthorized")
		}

		userID, err := m.service.GetUserID(c, sid)
		if err != nil {
			if errors.Is(err, session.ErrNotFound) {
				return httperror.New(fiber.StatusUnauthorized, "Unauthorized")
			}

			return err
		}

		rctx.SetUserID(c, userID)
		rctx.SetSessionID(c, sid)

		return c.Next()
	}
}
