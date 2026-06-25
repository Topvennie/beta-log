// Package middlewares contains the server middlewares
package middlewares

import (
	"errors"

	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const errMsg = "Internal Server Error"

func ErrorHandler(repo repository.Repository) fiber.ErrorHandler {
	repoUser := repo.NewUser()

	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
		}

		if code < 500 {
			if fiberErr != nil {
				return c.Status(code).SendString(fiberErr.Message)
			}
			return c.SendStatus(code)
		}

		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("route", c.Route().Path),
		}

		params := map[string]interface{}{}
		for _, param := range c.Route().Params {
			params[param] = c.Params(param, "")
		}
		if len(params) > 0 {
			fields = append(fields, zap.Any("params", params))
		}

		if data, ok := c.Locals("body").(map[string]interface{}); ok {
			fields = append(fields, zap.Any("body", data))
		}

		if id, ok := c.Locals("id").(int); ok {
			fields = append(fields, zap.Int("id", id))

			if user, err := repoUser.GetByID(c.Context(), id); err == nil && user != nil {
				fields = append(fields, zap.String("userName", user.Name))
			}
		}

		if uid, ok := c.Locals("uid").(string); ok {
			fields = append(fields, zap.String("uid", uid))
		}

		zap.L().Error(err.Error(), fields...)

		return c.Status(500).SendString(errMsg)
	}
}
