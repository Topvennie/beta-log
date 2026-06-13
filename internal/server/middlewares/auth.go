// Package middlewares contains the server middlewares
package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.uber.org/zap"
)

func AuthRoute(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess == nil {
		zap.S().Warn("Tried to authenticate a user without a session")
		return fiber.ErrUnauthorized
	}
	id := sess.Get("id")
	if id == nil || id == "" {
		return fiber.ErrUnauthorized
	}

	uid := sess.Get("uid")
	if uid == nil || uid == "" {
		return fiber.ErrUnauthorized
	}

	c.Locals("id", id)
	c.Locals("uid", uid)

	return c.Next()
}
