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
