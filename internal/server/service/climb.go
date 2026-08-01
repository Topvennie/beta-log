package service

import (
	"context"
	"slices"
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/grade"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Climb struct {
	climb   repository.Climb
	day     repository.ClimbDay
	setting repository.Setting
}

func NewClimb() *Climb {
	return &Climb{
		climb:   *repository.NewClimb(),
		day:     *repository.NewClimbDay(),
		setting: *repository.NewSetting(),
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

	setting, err := c.setting.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.SliceMap(days, func(c *model.ClimbDay) dto.ClimbDay { return dto.ClimbDayDTO(c, *setting) }), nil
}

func (c *Climb) CreateDay(ctx fiber.Ctx, daySave dto.ClimbDayCreate) (dto.ClimbDay, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.ClimbDay{}, err
	}
	setting, err := c.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.ClimbDay{}, err
	}

	day := daySave.ToModel(*setting)
	day.UserID = userID
	day.ExternalID = uuid.NewString()
	day.Source = model.SourceManual

	climbs := day.Climbs
	// Validate the climb grades
	for _, c := range climbs {
		if c.Grade == -1 {
			return dto.ClimbDay{}, fiber.NewError(fiber.StatusBadRequest, "Invalid climb grade")
		}
	}

	if err := withRollback(ctx, func(ctx context.Context) error {
		if err := c.day.Create(ctx, &day); err != nil {
			return err
		}

		for i := range climbs {
			climbs[i].UserID = userID
			climbs[i].ExternalID = uuid.NewString()
			climbs[i].Source = model.SourceManual
			climbs[i].ClimbDayID = day.ID

			if err := c.climb.Create(ctx, &climbs[i]); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return dto.ClimbDay{}, err
	}

	populated, err := c.day.GetPopulated(ctx, day.ID)
	if err != nil {
		return dto.ClimbDay{}, err
	}
	if populated == nil {
		return dto.ClimbDay{}, fiber.ErrNotFound
	}

	return dto.ClimbDayDTO(populated, *setting), nil
}

func (c *Climb) UpdateDay(ctx fiber.Ctx, daySave dto.ClimbDayUpdate) (dto.ClimbDay, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.ClimbDay{}, err
	}
	setting, err := c.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.ClimbDay{}, err
	}

	oldDay, err := c.day.Get(ctx, daySave.ID)
	if err != nil {
		return dto.ClimbDay{}, err
	}
	if oldDay == nil || oldDay.UserID != userID {
		return dto.ClimbDay{}, fiber.ErrNotFound
	}
	if oldDay.Source != model.SourceManual {
		return dto.ClimbDay{}, fiber.NewError(fiber.StatusBadRequest, "only manual climb days can be updated")
	}

	day := daySave.ToModel(*setting)
	day.UserID = userID
	day.ExternalID = oldDay.ExternalID
	day.Source = oldDay.Source

	climbs := day.Climbs
	// Validate the climb grades
	for _, c := range climbs {
		if c.Grade == -1 {
			return dto.ClimbDay{}, fiber.NewError(fiber.StatusBadRequest, "Invalid climb grade")
		}
	}

	if err := withRollback(ctx, func(ctx context.Context) error {
		if err := c.day.Update(ctx, day); err != nil {
			return err
		}

		if err := c.climb.DeleteByClimbDay(ctx, day.ID); err != nil {
			return err
		}

		for i := range climbs {
			climbs[i].UserID = userID
			climbs[i].ExternalID = uuid.NewString()
			climbs[i].Source = model.SourceManual
			climbs[i].ClimbDayID = day.ID

			if err := c.climb.Create(ctx, &climbs[i]); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return dto.ClimbDay{}, err
	}

	populated, err := c.day.GetPopulated(ctx, day.ID)
	if err != nil {
		return dto.ClimbDay{}, err
	}
	if populated == nil {
		return dto.ClimbDay{}, fiber.ErrNotFound
	}

	return dto.ClimbDayDTO(populated, *setting), nil
}

func (c *Climb) DeleteDay(ctx fiber.Ctx, id int) error {
	userID, err := getID(ctx)
	if err != nil {
		return err
	}

	day, err := c.day.Get(ctx, id)
	if err != nil {
		return err
	}
	if day == nil || day.UserID != userID {
		return fiber.ErrNotFound
	}
	if day.Source != model.SourceManual {
		return fiber.NewError(fiber.StatusBadRequest, "only manual climb days can be deleted")
	}

	return withRollback(ctx, func(ictx context.Context) error {
		if err := c.climb.DeleteByClimbDay(ictx, id); err != nil {
			return err
		}
		if err := c.day.Delete(ictx, id); err != nil {
			return err
		}
		return nil
	})
}

func (c *Climb) GetStats(ctx fiber.Ctx, start, end time.Time) (dto.ClimbStats, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.ClimbStats{}, err
	}
	setting, err := c.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.ClimbStats{}, err
	}

	days, err := c.day.GetAllPopulatedByUser(ctx, userID)
	if err != nil {
		return dto.ClimbStats{}, err
	}
	// Filter out the days outside of the range
	utils.SliceFilter(days, func(d *model.ClimbDay) bool {
		if !start.IsZero() && d.Date.Before(start) {
			return false
		}
		if !end.IsZero() && d.Date.After(end) {
			return false
		}

		return true
	})

	// Get all stats
	stats := getStatsBasic(days, *setting)
	stats.GraphProgress = getStatsProgress(days, *setting)
	stats.GraphPerGrade = getStatsGrade(days, *setting)

	return stats, nil
}

func getStatsBasic(days []*model.ClimbDay, setting model.Setting) dto.ClimbStats {
	var stats dto.ClimbStats

	var bestTop, bestTopAmount int
	var bestFlash, bestFlashAmount int

	climbsPerSession := make([]int, 0, len(days)) // Used to calculate the median

	for _, day := range days {
		stats.Sessions++
		climbsPerSession = append(climbsPerSession, len(day.Climbs))

		for _, climb := range day.Climbs {
			stats.Total++

			switch climb.FinishType {
			case model.ClimbFinishFlash:
				stats.Flash++
				stats.TotalUnique++

				if climb.Grade >= bestFlash {
					if climb.Grade > bestFlash {
						bestFlash = climb.Grade
						bestFlashAmount = 0
					}

					bestFlashAmount++
				}

				if climb.Grade >= bestTop {
					if climb.Grade > bestTop {
						bestTop = climb.Grade
						bestTopAmount = 0
					}

					bestTopAmount++
				}

			case model.ClimbFinishTop:
				stats.Top++
				stats.TotalUnique++

				if climb.Grade >= bestTop {
					if climb.Grade > bestTop {
						bestTop = climb.Grade
						bestTopAmount = 0
					}

					bestTopAmount++
				}

			case model.ClimbFinishRepeat:
				stats.Repeat++
			}

			switch climb.ClimbType {
			case model.ClimbTypeBoulder:
				stats.Boulder++

			case model.ClimbTypeLead:
				stats.Lead++
			}
		}
	}

	stats.Best = grade.Grade(bestTop).Format(setting.GradeSystem)
	stats.BestAmount = bestTopAmount
	stats.BestFlash = grade.Grade(bestFlash).Format(setting.GradeSystem)
	stats.BestFlashAmount = bestFlashAmount

	// Get median
	slices.Sort(climbsPerSession)
	if n := len(climbsPerSession); n > 0 {
		if n%2 == 1 {
			stats.MedianClimbsPerSession = climbsPerSession[n/2]
		} else {
			stats.MedianClimbsPerSession = (climbsPerSession[n/2-1] + climbsPerSession[n/2] + 1) / 2
		}
	}

	return stats
}

func getStatsProgress(days []*model.ClimbDay, setting model.Setting) []dto.ClimbStatsProgress {
	stats := make([]dto.ClimbStatsProgress, 0, len(days))

	for _, day := range days {
		best := 0
		for _, c := range day.Climbs {
			if c.Grade > best {
				best = c.Grade
			}
		}

		stats = append(stats, dto.ClimbStatsProgress{
			Date:   day.Date.Format("Jan 02"),
			Grade:  grade.Grade(best).Format(setting.GradeSystem),
			Volume: len(day.Climbs),
		})
	}

	return stats
}

func getStatsGrade(days []*model.ClimbDay, setting model.Setting) []dto.ClimbStatsGrade {
	type climbStatsGrade struct {
		dto      dto.ClimbStatsGrade
		gradeRaw int
	}

	stats := make(map[string]climbStatsGrade)

	for _, day := range days {
		for _, c := range day.Climbs {
			g := grade.Grade(c.Grade).Format(setting.GradeSystem)

			stat, ok := stats[g]
			if !ok {
				stat = climbStatsGrade{
					dto: dto.ClimbStatsGrade{
						Grade: g,
					},
					gradeRaw: c.Grade,
				}
			}

			switch c.FinishType {
			case model.ClimbFinishFlash:
				stat.dto.Flash++
			case model.ClimbFinishTop:
				stat.dto.Top++
			case model.ClimbFinishRepeat:
				stat.dto.Repeat++
			}

			stats[g] = stat
		}
	}

	values := utils.MapValues(stats)
	slices.SortFunc(values, func(a, b climbStatsGrade) int { return a.gradeRaw - b.gradeRaw })
	return utils.SliceMap(values, func(c climbStatsGrade) dto.ClimbStatsGrade { return c.dto })
}
