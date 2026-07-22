package api

import (
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
	cl.router.Get("/days", cl.getDays)
	cl.router.Get("/stats", cl.getStats)
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
	stats, err := cl.climb.GetStats(c)
	if err != nil {
		return err
	}

	return c.JSON(stats)
}
