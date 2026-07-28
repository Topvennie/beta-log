package service

import (
	"slices"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Gym struct {
	gym repository.Gym
	day repository.ClimbDay
}

func NewGym() *Gym {
	return &Gym{
		gym: *repository.NewGym(),
		day: *repository.NewClimbDay(),
	}
}

func (g *Gym) GetAll(ctx fiber.Ctx) ([]dto.Gym, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	gyms, err := g.gym.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.SliceMap(gyms, dto.GymDTO), nil
}

func (g *Gym) Create(ctx fiber.Ctx, gymSave dto.GymCreate) (dto.Gym, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Gym{}, err
	}

	gym := gymSave.ToModel()
	gym.UserID = userID
	gym.ExternalID = uuid.NewString()
	gym.Source = model.SourceManual

	if err := g.gym.Create(ctx, &gym); err != nil {
		return dto.Gym{}, err
	}

	return dto.GymDTO(&gym), nil
}

func (g *Gym) Update(ctx fiber.Ctx, gymSave dto.GymUpdate) (dto.Gym, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Gym{}, err
	}

	oldGym, err := g.gym.Get(ctx, gymSave.ID)
	if err != nil {
		return dto.Gym{}, err
	}
	if oldGym == nil || oldGym.UserID != userID {
		return dto.Gym{}, fiber.ErrNotFound
	}
	if oldGym.Source != model.SourceManual {
		return dto.Gym{}, fiber.NewError(fiber.StatusBadRequest, "only manual gyms can be updated")
	}

	gym := gymSave.ToModel()
	gym.UserID = userID
	gym.ExternalID = oldGym.ExternalID
	gym.Source = oldGym.Source

	if err := g.gym.Update(ctx, gym); err != nil {
		return dto.Gym{}, err
	}

	return dto.GymDTO(&gym), nil
}

func (g *Gym) Delete(ctx fiber.Ctx, id int) error {
	userID, err := getID(ctx)
	if err != nil {
		return err
	}

	gym, err := g.gym.Get(ctx, id)
	if err != nil {
		return err
	}
	if gym == nil || gym.UserID != userID {
		return fiber.ErrNotFound
	}
	if gym.Source != model.SourceManual {
		return fiber.NewError(fiber.StatusBadRequest, "only manual gyms can be deleted")
	}

	if err := g.gym.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (g *Gym) GetStats(ctx fiber.Ctx) (dto.GymStats, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.GymStats{}, err
	}

	gyms, err := g.gym.GetAllByUser(ctx, userID)
	if err != nil {
		return dto.GymStats{}, err
	}

	days, err := g.day.GetAllPopulatedByUser(ctx, userID)
	if err != nil {
		return dto.GymStats{}, err
	}

	stats := dto.GymStats{
		GraphVisits:       make([]dto.GymStatsVisits, 0, len(gyms)),
		GraphTop:          make([]dto.GymStatsTop, 0, len(gyms)),
		GraphDistribution: make([]dto.GymStatsDistribution, 0, len(gyms)),
	}
	stats.Total = len(gyms)
	stats.Sessions = len(days)

	visits := make(map[int]int)
	bestTop := make(map[int]int)
	bestFlash := make(map[int]int)
	gradeCounts := make(map[int]map[int]int, len(gyms))

	for _, gym := range gyms {
		gradeCounts[gym.ID] = make(map[int]int)
	}

	for _, day := range days {
		visits[day.GymID]++

		for _, climb := range day.Climbs {
			gradeCounts[day.GymID][climb.Grade]++

			switch climb.FinishType {
			case model.ClimbFinishTop:
				if climb.Grade > bestTop[day.GymID] {
					bestTop[day.GymID] = climb.Grade
				}
			case model.ClimbFinishFlash:
				if climb.Grade > bestFlash[day.GymID] {
					bestFlash[day.GymID] = climb.Grade
				}
			case model.ClimbFinishRepeat:
			}
		}
	}

	mostVisited := ""
	mostVisitedAmount := 0
	for _, gym := range gyms {
		v := visits[gym.ID]
		if v > mostVisitedAmount {
			mostVisitedAmount = v
			mostVisited = gym.Name
		}
	}

	stats.MostVisited = mostVisited
	stats.MostVisitedAmount = mostVisitedAmount

	for _, gym := range gyms {
		stats.GraphVisits = append(stats.GraphVisits, dto.GymStatsVisits{
			Gym:    gym.Name,
			Amount: visits[gym.ID],
		})

		stats.GraphTop = append(stats.GraphTop, dto.GymStatsTop{
			Gym:   gym.Name,
			Top:   bestTop[gym.ID],
			Flash: bestFlash[gym.ID],
		})

		total := 0
		for _, c := range gradeCounts[gym.ID] {
			total += c
		}

		dist := make(map[int]int)
		if total > 0 {
			type entry struct {
				grade   int
				percent int
				frac    float64
			}
			entries := make([]entry, 0, len(gradeCounts[gym.ID]))
			sum := 0
			for grade, c := range gradeCounts[gym.ID] {
				exact := float64(c) * 100 / float64(total)
				floor := int(exact)
				entries = append(entries, entry{grade: grade, percent: floor, frac: exact - float64(floor)})
				sum += floor
			}

			remainder := 100 - sum
			slices.SortFunc(entries, func(a, b entry) int {
				if a.frac != b.frac {
					if a.frac > b.frac {
						return -1
					}
					return 1
				}
				return a.grade - b.grade
			})
			for i := 0; i < remainder; i++ {
				entries[i].percent++
			}

			for _, e := range entries {
				dist[e.grade] = e.percent
			}
		}

		stats.GraphDistribution = append(stats.GraphDistribution, dto.GymStatsDistribution{
			Gym:          gym.Name,
			Distribution: dist,
		})
	}

	return stats, nil
}
