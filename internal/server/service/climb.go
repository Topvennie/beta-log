package service

import (
	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type Climb struct {
	climb repository.Climb
	day   repository.ClimbDay
	gym   repository.ClimbGym
}

func NewClimb() *Climb {
	return &Climb{
		climb: *repository.NewClimb(),
		day:   *repository.NewClimbDay(),
		gym:   *repository.NewClimbGym(),
	}
}

func (c *Climb) GetDays(ctx fiber.Ctx, filter dto.ClimbDayFilter) ([]dto.ClimbDay, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	filter.UserID = userID
	days, err := c.day.GetAllPopulatedFiltered(ctx, filter.ToModel())
	if err != nil {
		return nil, err
	}

	return utils.SliceMap(days, dto.ClimbDayDTO), nil
}

func (c *Climb) GetStats(ctx fiber.Ctx) (dto.ClimbStats, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.ClimbStats{}, err
	}

	days, err := c.day.GetAllPopulatedByUser(ctx, userID)
	if err != nil {
		return dto.ClimbStats{}, nil
	}

	// Get all stats
	stats := dto.ClimbStats{
		Sessions: len(days),
	}

	for _, day := range days {
		stats.Total += len(day.Climbs)

		for _, climb := range day.Climbs {
			if climb.Grade > stats.Top {
				stats.Top = climb.Grade
			}

			if climb.FinishType == model.ClimbFinishFlash {
				if climb.Grade > stats.TopFlash {
					stats.TopFlash = climb.Grade
				}
			}

			if climb.FinishType == model.ClimbFinishRepeat {
				stats.Repeats++
			}
		}
	}

	return stats, nil
}
