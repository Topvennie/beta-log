package api

import (
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type gym struct {
	router fiber.Router
	gym    *service.Gym
}

func newGym(router fiber.Router) *gym {
	api := &gym{
		router: router.Group("/gym"),
		gym:    service.NewGym(),
	}

	api.routes()

	return api
}

func (g *gym) routes() {
	g.router.Get("/", g.getAll)
	g.router.Post("/", g.create)
	g.router.Put("/:id", g.update)
}

func (g *gym) getAll(c fiber.Ctx) error {
	gyms, err := g.gym.GetAll(c)
	if err != nil {
		return err
	}

	return c.JSON(gyms)
}

func (g *gym) create(c fiber.Ctx) error {
	var gym dto.GymCreate
	if err := c.Bind().Body(&gym); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(gym); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newGym, err := g.gym.Create(c, gym)
	if err != nil {
		return err
	}

	return c.JSON(newGym)
}

func (g *gym) update(c fiber.Ctx) error {
	var gym dto.GymUpdate
	if err := c.Bind().Body(&gym); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(gym); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id := fiber.Params[int](c, "id")
	if id != gym.ID {
		return fiber.NewError(fiber.StatusBadRequest, "params id doesn't match body id")
	}

	newGym, err := g.gym.Update(c, gym)
	if err != nil {
		return err
	}

	return c.JSON(newGym)
}
