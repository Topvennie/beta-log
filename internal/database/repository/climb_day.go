package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type ClimbDay struct {
	repo Repository
}

func (r *Repository) NewClimbDay() *ClimbDay {
	return &ClimbDay{repo: *r}
}

func (c *ClimbDay) Get(ctx context.Context, id int) (*model.ClimbDay, error) {
	day, err := c.repo.queries(ctx).ClimbDayGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day %d | %w", id, err)
	}

	return model.ClimbDayModel(day), nil
}

func (c *ClimbDay) GetByExternalID(ctx context.Context, externalID string) (*model.ClimbDay, error) {
	day, err := c.repo.queries(ctx).ClimbDayGetByExternal(ctx, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day by external id %s | %w", externalID, err)
	}

	return model.ClimbDayModel(day), nil
}

func (c *ClimbDay) GetPopulated(ctx context.Context, id int) (*model.ClimbDay, error) {
	rows, err := c.repo.queries(ctx).ClimbDayGetPopulated(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day populated by id %d | %w", id, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	day := model.ClimbDayModel(rows[0].ClimbDay)
	day.Gym = model.ClimbGymPopulatedModel(rows[0].ClimbGym)

	for _, row := range rows {
		climb := model.ClimbPopulatedModel(row.Climb)
		day.Climbs = append(day.Climbs, *climb)
	}

	return day, nil
}

func (c *ClimbDay) GetPopulatedByExternal(ctx context.Context, externalID string) (*model.ClimbDay, error) {
	rows, err := c.repo.queries(ctx).ClimbDayGetPopulatedByExternal(ctx, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day populated by external id %s | %w", externalID, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	day := model.ClimbDayModel(rows[0].ClimbDay)
	day.Gym = model.ClimbGymPopulatedModel(rows[0].ClimbGym)

	for _, row := range rows {
		climb := model.ClimbPopulatedModel(row.Climb)
		day.Climbs = append(day.Climbs, *climb)
	}

	return day, nil
}

func (c *ClimbDay) Create(ctx context.Context, day *model.ClimbDay) error {
	id, err := c.repo.queries(ctx).ClimbDayCreate(ctx, sqlc.ClimbDayCreateParams{
		UserID:     int32(day.UserID),
		ExternalID: day.ExternalID,
		GymID:      int32(day.GymID),
		Day:        pgtype.Timestamptz{Time: day.Day, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create climb day %+v | %w", *day, err)
	}

	day.ID = int(id)

	return nil
}

func (c *ClimbDay) Update(ctx context.Context, day model.ClimbDay) error {
	err := c.repo.queries(ctx).ClimbDayUpdate(ctx, sqlc.ClimbDayUpdateParams{
		ID:    int32(day.ID),
		GymID: int32(day.GymID),
		Day:   pgtype.Timestamptz{Time: day.Day, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("update climb day %+v | %w", day, err)
	}

	return nil
}
