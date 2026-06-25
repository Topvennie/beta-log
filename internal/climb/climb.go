package climb

import (
	"context"
	"errors"
	"fmt"
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
				msg, err := m.update(ctx, user)
				results = append(results, task.TaskResult{
					User:    user,
					Message: msg,
					Error:   err,
				})
			}

			return results
		},
	)); err != nil {
		if errors.Is(err, task.ErrTaskExists) {
			return fmt.Errorf("task is already running %w", err)
		}
		return err
	}

	return nil
}

func (m *Manager) update(ctx context.Context, user model.User) (string, error) {
}
