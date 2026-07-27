package api

import (
	"time"

	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type climb struct {
	router fiber.Router
	climb  *service.Climb
}

func newClimb(router fiber.Router) *climb {
	api := &climb{
		router: router.Group("/climb"),
		climb:  service.NewClimb(),
	}

	api.routes()

	return api
}

func (cl *climb) routes() {
	cl.router.Get("/day", cl.getDays)
	cl.router.Get("/stat", cl.getStats)
	cl.router.Get("/gym", cl.getGyms)
	cl.router.Post("/gym", cl.createGym)
	cl.router.Put("/gym/:id", cl.updateGym)
}

func (cl *climb) getDays(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
	page := fiber.Query[int](c, "page", 1)

	days, err := cl.climb.GetDays(c, dto.ClimbDayFilter{
		Limit:  limit,
		Offset: (page - 1) * limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(days)
}

func (cl *climb) getStats(c fiber.Ctx) error {
	startStr := fiber.Query[string](c, "start", "")
	endStr := fiber.Query[string](c, "end", "")

	var err error

	start := time.Time{}
	if startStr != "" {
		start, err = time.Parse("02-01-2006", startStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid start format")
		}
	}

	end := time.Time{}
	if endStr != "" {
		end, err = time.Parse("02-01-2006", endStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid end format")
		}
	}

	stats, err := cl.climb.GetStats(c, start, end)
	if err != nil {
		return err
	}

	return c.JSON(stats)
}

func (cl *climb) getGyms(c fiber.Ctx) error {
	gyms, err := cl.climb.GetGyms(c)
	if err != nil {
		return err
	}

	return c.JSON(gyms)
}

func (cl *climb) createGym(c fiber.Ctx) error {
	var gym dto.ClimbGymCreate
	if err := c.Bind().Body(&gym); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(gym); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newGym, err := cl.climb.CreateGym(c, gym)
	if err != nil {
		return err
	}

	return c.JSON(newGym)
}

func (cl *climb) updateGym(c fiber.Ctx) error {
	var gym dto.ClimbGymUpdate
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

	newGym, err := cl.climb.UpdateGym(c, gym)
	if err != nil {
		return err
	}

	return c.JSON(newGym)
}
