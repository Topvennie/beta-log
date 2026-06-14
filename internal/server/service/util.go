package service

import "github.com/gofiber/fiber/v3"

func getID(ctx fiber.Ctx) (int, error) {
	id, ok := ctx.Locals("id").(int)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}

	return id, nil
}
