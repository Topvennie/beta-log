package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/sqlc"
	"github.com/Topvennie/beta-log/pkg/utils"
)

type Gym struct{}

func NewGym() *Gym {
	return &Gym{}
}

func (c *Gym) Get(ctx context.Context, id int) (*model.Gym, error) {
	gym, err := queries(ctx).GymGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get gym %d | %w", id, err)
	}

	return model.GymModel(gym), nil
}

func (c *Gym) GetByExternalSource(ctx context.Context, source model.Source, externalID string) (*model.Gym, error) {
	gym, err := queries(ctx).GymGetByExternalSource(ctx, sqlc.GymGetByExternalSourceParams{
		ExternalID: externalID,
		Source:     sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get gym by external id %s and source %s | %w", externalID, source, err)
	}

	return model.GymModel(gym), nil
}

func (c *Gym) GetByExternalSourceIDs(ctx context.Context, source model.Source, externalIDs []int) ([]*model.Gym, error) {
	gyms, err := queries(ctx).GymGetAllByExternalSource(ctx, sqlc.GymGetAllByExternalSourceParams{
		Column1: utils.SliceMap(externalIDs, func(id int) int32 { return int32(id) }),
		Source:  sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get gyms by external ids %v and source %s | %w", externalIDs, source, err)
	}

	return utils.SliceMap(gyms, model.GymModel), nil
}

func (c *Gym) GetAllByUser(ctx context.Context, userID int) ([]*model.Gym, error) {
	gyms, err := queries(ctx).GymGetAllByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get gyms by user id %d | %w", userID, err)
	}

	return utils.SliceMap(gyms, model.GymModel), nil
}

func (c *Gym) Create(ctx context.Context, gym *model.Gym) error {
	id, err := queries(ctx).GymCreate(ctx, sqlc.GymCreateParams{
		UserID:     int32(gym.UserID),
		ExternalID: gym.ExternalID,
		Name:       gym.Name,
		IconPath:   gym.IconPath,
		Source:     sqlc.ClimbSource(gym.Source),
	})
	if err != nil {
		return fmt.Errorf("create gym %+v | %w", *gym, err)
	}

	gym.ID = int(id)

	return nil
}

func (c *Gym) Update(ctx context.Context, gym model.Gym) error {
	if err := queries(ctx).GymUpdate(ctx, sqlc.GymUpdateParams{
		ID:       int32(gym.ID),
		Name:     gym.Name,
		IconPath: gym.IconPath,
	}); err != nil {
		return fmt.Errorf("update gym %+v | %w", gym, err)
	}

	return nil
}

func (c *Gym) Delete(ctx context.Context, id int) error {
	if err := queries(ctx).GymDelete(ctx, int32(id)); err != nil {
		return fmt.Errorf("delete gym with id %d | %w", id, err)
	}

	return nil
}
