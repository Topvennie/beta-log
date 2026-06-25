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

type Climb struct {
	repo Repository
}

func (r *Repository) NewClimb() *Climb {
	return &Climb{repo: *r}
}

func (c *Climb) Get(ctx context.Context, id int) (*model.Climb, error) {
	climb, err := c.repo.queries(ctx).ClimbGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb %d | %w", id, err)
	}

	return model.ClimbModel(climb), nil
}

func (c *Climb) GetByExternalID(ctx context.Context, externalID string) (*model.Climb, error) {
	climb, err := c.repo.queries(ctx).ClimbGetByExternal(ctx, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb by external id %s | %w", externalID, err)
	}

	return model.ClimbModel(climb), nil
}

func (c *Climb) GetAllByClimbDayID(ctx context.Context, climbDayID int) ([]*model.Climb, error) {
	climbs, err := c.repo.queries(ctx).ClimbGetAllByClimbDay(ctx, int32(climbDayID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all climbs by climb day id %d | %w", climbDayID, err)
	}

	return utils.SliceMap(climbs, model.ClimbModel), nil
}

func (c *Climb) Create(ctx context.Context, climb *model.Climb) error {
	id, err := c.repo.queries(ctx).ClimbCreate(ctx, sqlc.ClimbCreateParams{
		UserID:     int32(climb.UserID),
		ExternalID: climb.ExternalID,
		ClimbDayID: int32(climb.ClimbDayID),
		Grade:      int32(climb.Grade),
		Color:      climb.Color,
		HoldColor:  climb.HoldColor,
		ClimbType:  sqlc.ClimbType(climb.ClimbType),
		FinishType: sqlc.FinishType(climb.FinishType),
	})
	if err != nil {
		return fmt.Errorf("create climb %+v | %w", *climb, err)
	}

	climb.ID = int(id)

	return nil
}

func (c *Climb) Update(ctx context.Context, climb model.Climb) error {
	if err := c.repo.queries(ctx).ClimbUpdate(ctx, sqlc.ClimbUpdateParams{
		ID:         int32(climb.ID),
		Grade:      int32(climb.Grade),
		Color:      climb.Color,
		HoldColor:  climb.HoldColor,
		ClimbType:  sqlc.ClimbType(climb.ClimbType),
		FinishType: sqlc.FinishType(climb.FinishType),
	}); err != nil {
		return fmt.Errorf("update climb %+v | %w", climb, err)
	}

	return nil
}
