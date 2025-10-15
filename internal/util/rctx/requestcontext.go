package rctx

import "github.com/gofiber/fiber/v3"

const (
	UserIDKey    = "userID"
	SessionIDKey = "sessionID"
)

func SetUserID(c fiber.Ctx, userID string) {
	c.Locals(UserIDKey, userID)
}

func GetUserID(c fiber.Ctx) string {
	id, ok := c.Locals(UserIDKey).(string)
	if !ok {
		panic("userID not set in context")
	}

	return id
}

func SetSessionID(c fiber.Ctx, sessionID string) {
	c.Locals(SessionIDKey, sessionID)
}

func GetSessionID(c fiber.Ctx) string {
	id, ok := c.Locals(SessionIDKey).(string)
	if !ok {
		panic("sessionID not set in context")
	}

	return id
}
