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

type ClimbGym struct{}

func NewClimbGym() *ClimbGym {
	return &ClimbGym{}
}

func (c *ClimbGym) Get(ctx context.Context, id int) (*model.ClimbGym, error) {
	gym, err := queries(ctx).ClimbGymGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb gym %d | %w", id, err)
	}

	return model.ClimbGymModel(gym), nil
}

func (c *ClimbGym) GetByExternalSource(ctx context.Context, source model.ClimbSource, externalID string) (*model.ClimbGym, error) {
	gym, err := queries(ctx).ClimbGymGetByExternalSource(ctx, sqlc.ClimbGymGetByExternalSourceParams{
		ExternalID: externalID,
		Source:     sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb gym by external id %s and source %s | %w", externalID, source, err)
	}

	return model.ClimbGymModel(gym), nil
}

func (c *ClimbGym) GetByExternalSourceIDs(ctx context.Context, source model.ClimbSource, externalIDs []int) ([]*model.ClimbGym, error) {
	gyms, err := queries(ctx).ClimbGymGetAllByExternalSource(ctx, sqlc.ClimbGymGetAllByExternalSourceParams{
		Column1: utils.SliceMap(externalIDs, func(id int) int32 { return int32(id) }),
		Source:  sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb gyms by external ids %v and source %s | %w", externalIDs, source, err)
	}

	return utils.SliceMap(gyms, model.ClimbGymModel), nil
}

func (c *ClimbGym) Create(ctx context.Context, gym *model.ClimbGym) error {
	id, err := queries(ctx).ClimbGymCreate(ctx, sqlc.ClimbGymCreateParams{
		UserID:     int32(gym.UserID),
		ExternalID: gym.ExternalID,
		Name:       gym.Name,
		IconPath:   gym.IconPath,
		Source:     sqlc.ClimbSource(gym.Source),
	})
	if err != nil {
		return fmt.Errorf("create climb gym %+v | %w", *gym, err)
	}

	gym.ID = int(id)

	return nil
}

func (c *ClimbGym) Update(ctx context.Context, gym model.ClimbGym) error {
	if err := queries(ctx).ClimbGymUpdate(ctx, sqlc.ClimbGymUpdateParams{
		ID:       int32(gym.ID),
		Name:     gym.Name,
		IconPath: gym.IconPath,
	}); err != nil {
		return fmt.Errorf("update climb gym %+v | %w", gym, err)
	}

	return nil
}
