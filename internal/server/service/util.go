package service

import (
	"context"

	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/gofiber/fiber/v3"
)

func withRollback(ctx fiber.Ctx, fn func(context.Context) error) error {
	return repository.WithRollback(ctx, fn)
}

func getID(ctx fiber.Ctx) (int, error) {
	id, ok := ctx.Locals("id").(int)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}

	return id, nil
}
