package service

import (
	"slices"
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type Climb struct {
	climb repository.Climb
	day   repository.ClimbDay
}

func NewClimb() *Climb {
	return &Climb{
		climb: *repository.NewClimb(),
		day:   *repository.NewClimbDay(),
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

func (c *Climb) GetStats(ctx fiber.Ctx, start, end time.Time) (dto.ClimbStats, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.ClimbStats{}, err
	}

	days, err := c.day.GetAllPopulatedByUser(ctx, userID)
	if err != nil {
		return dto.ClimbStats{}, err
	}

	// Get all stats
	stats := dto.ClimbStats{
		GraphProgress: make([]dto.ClimbStatsGraphProgress, 0, len(days)),
		GraphPerGrade: []dto.ClimbStatsGraphGrade{},
	}

	graphGrades := make(map[int]dto.ClimbStatsGraphGrade)
	perSessionClimbs := make([]int, 0, len(days))

	for _, day := range days {
		if !start.IsZero() && day.Date.Before(start) {
			continue
		}
		if !end.IsZero() && day.Date.After(end) {
			break
		}

		stats.Sessions++
		stats.Total += len(day.Climbs)
		perSessionClimbs = append(perSessionClimbs, len(day.Climbs))

		dayBest := 0

		for _, climb := range day.Climbs {
			graphGrade, ok := graphGrades[climb.Grade]
			if !ok {
				graphGrade = dto.ClimbStatsGraphGrade{
					Grade: climb.Grade,
				}
			}

			switch climb.FinishType {
			case model.ClimbFinishFlash:
				stats.Flash++
				stats.TotalUnique++
				graphGrade.Flash++
			case model.ClimbFinishTop:
				stats.Top++
				stats.TotalUnique++
				graphGrade.Top++
			case model.ClimbFinishRepeat:
				stats.Repeat++
				graphGrade.Repeat++
			}

			// Best grade is advanced only by tops and flashes
			// Repeats don't move your best.
			// BestAmount counts every send (top + flash + repeat) at the current best grade.
			if (climb.FinishType == model.ClimbFinishFlash || climb.FinishType == model.ClimbFinishTop) && climb.Grade > stats.Best {
				stats.Best = climb.Grade
				stats.BestAmount = 0
			}
			if stats.Best > 0 && climb.Grade == stats.Best {
				stats.BestAmount++
			}

			// Best flash is flashes only.
			if climb.FinishType == model.ClimbFinishFlash {
				if climb.Grade > stats.BestFlash {
					stats.BestFlash = climb.Grade
					stats.BestFlashAmount = 0
				}
				if stats.BestFlash > 0 && climb.Grade == stats.BestFlash {
					stats.BestFlashAmount++
				}
			}

			if climb.Grade > dayBest {
				dayBest = climb.Grade
			}

			graphGrades[climb.Grade] = graphGrade
		}

		if dayBest > 0 {
			stats.GraphProgress = append(stats.GraphProgress, dto.ClimbStatsGraphProgress{
				Date:   day.Date.Format("Jan 02"),
				Grade:  dayBest,
				Volume: len(day.Climbs),
			})
		}
	}

	slices.Sort(perSessionClimbs)
	if n := len(perSessionClimbs); n > 0 {
		if n%2 == 1 {
			stats.MedianClimbsPerSession = perSessionClimbs[n/2]
		} else {
			stats.MedianClimbsPerSession = (perSessionClimbs[n/2-1] + perSessionClimbs[n/2] + 1) / 2
		}
	}
	stats.GraphPerGrade = utils.MapValues(graphGrades)
	slices.SortFunc(stats.GraphPerGrade, func(a, b dto.ClimbStatsGraphGrade) int { return a.Grade - b.Grade })

	return stats, nil
}
