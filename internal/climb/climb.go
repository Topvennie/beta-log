package climb

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/task"
	"github.com/Topvennie/beta-log/pkg/config"
)

const taskUpdateUID = "task-climb-update"

type Fetcher interface {
	Name() string
	Fetch(context.Context, model.User) ([]model.ClimbDay, error)
}

type Manager struct {
	interval time.Duration

	fetchers []Fetcher

	climb    repository.Climb
	climbDay repository.ClimbDay
	climbGym repository.ClimbGym
}

func New(repo repository.Repository) *Manager {
	return &Manager{
		interval: config.GetDefaultDurationS("climb.interval_s", 3600),
		fetchers: []Fetcher{},
		climb:    *repo.NewClimb(),
		climbDay: *repo.NewClimbDay(),
		climbGym: *repo.NewClimbGym(),
	}
}

func (m *Manager) Start(ctx context.Context) error {
	if err := task.Manager.Add(ctx, task.NewTask(
		taskUpdateUID,
		"Climbs Update",
		m.interval,
		false,
		func(ctx context.Context, users []model.User) []task.TaskResult {
			results := make([]task.TaskResult, 0, len(users))
			for _, user := range users {
				msg, err := m.updateAll(ctx, user)
				results = append(results, task.TaskResult{
					User:    user,
					Message: msg,
					Error:   err,
				})
			}

			return results
		},
	)); err != nil {
		return err
	}

	return nil
}

func (m *Manager) updateAll(ctx context.Context, user model.User) (string, error) {
	var errs []error
	totalNewClimbs := 0

	for _, fetcher := range m.fetchers {
		newClimbs, err := m.update(ctx, user, fetcher)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		totalNewClimbs += newClimbs
	}

	var msg string
	if totalNewClimbs > 0 {
		msg = fmt.Sprintf("Added %d new climbs", totalNewClimbs)
	}

	return msg, errors.Join(errs...)
}

// update gets all climb days for a single user and fetcher, updates the database and returns the amount of added climbs
func (m *Manager) update(ctx context.Context, user model.User, fetcher Fetcher) (int, error) {
	newClimbs := 0

	days, err := fetcher.Fetch(ctx, user)
	if err != nil {
		return 0, err
	}

	for _, day := range days {
		// Create gym if  necessary
		dbGym, err := m.climbGym.GetByExternalID(ctx, day.Gym.ExternalID)
		if err != nil {
			return 0, err
		}
		if dbGym == nil {
			if err := m.climbGym.Create(ctx, &day.Gym); err != nil {
				return 0, err
			}

			day.GymID = day.Gym.ID
		} else {
			day.GymID = dbGym.ID
		}

		// Create day if necessary
		dbDay, err := m.climbDay.GetPopulatedByExternal(ctx, day.ExternalID)
		if err != nil {
			return 0, err
		}
		if dbDay == nil {
			if err := m.climbDay.Create(ctx, &day); err != nil {
				return 0, err
			}
		} else {
			day.ID = dbDay.ID
		}

		// Add all missing climbs
		// Make a copy that we can safely mutate
		dbClimbs := make([]model.Climb, 0, len(dbDay.Climbs))
		copy(dbClimbs, dbDay.Climbs)

		for _, climb := range day.Climbs {
			climb.ClimbDayID = day.ID

			if idx := slices.IndexFunc(dbClimbs, func(c model.Climb) bool { return c.ExternalID == climb.ExternalID && c.ClimbDayID == climb.ClimbDayID }); idx != -1 {
				// Climb found
				// Remove it from the slice to support finishing the same climb multiple times in the same day (adding multiple entries)
				dbClimbs[idx] = dbClimbs[len(dbClimbs)-1]
				dbClimbs = dbClimbs[:len(dbClimbs)-1]

				continue
			}

			// Climb not found
			if err := m.climb.Create(ctx, &climb); err != nil {
				return 0, err
			}

			newClimbs++
		}
	}

	return newClimbs, nil
}
