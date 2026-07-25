package service

import (
	"math"
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

func (c *Climb) GetStats(ctx fiber.Ctx, start time.Time, end time.Time) (dto.ClimbStats, error) {
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

				if climb.Grade >= stats.Best {
					if climb.Grade > stats.Best {
						stats.Best = climb.Grade
						stats.BestAmount = 0
					}
					stats.BestAmount++
				}
				if climb.Grade >= stats.BestFlash {
					if climb.Grade > stats.BestFlash {
						stats.BestFlash = climb.Grade
						stats.BestFlashAmount = 1
					}
					stats.BestFlashAmount++
				}

			case model.ClimbFinishTop:
				stats.Top++
				stats.TotalUnique++
				graphGrade.Top++

				if climb.Grade > stats.Best {
					stats.Best = climb.Grade
				}

			case model.ClimbFinishRepeat:
				graphGrade.Repeat++
				stats.Repeat++
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
	if len(perSessionClimbs) > 0 {
		n := len(perSessionClimbs)
		if n%2 == 1 {
			stats.MedianClimbsPerSession = float64(perSessionClimbs[n/2])
		} else {
			stats.MedianClimbsPerSession = math.Round(float64(perSessionClimbs[n/2-1]+perSessionClimbs[n/2]) / 2)
		}
	}
	stats.GraphPerGrade = utils.MapValues(graphGrades)
	slices.SortFunc(stats.GraphPerGrade, func(a, b dto.ClimbStatsGraphGrade) int { return a.Grade - b.Grade })

	return stats, nil
}
