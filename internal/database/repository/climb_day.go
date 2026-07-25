package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/sqlc"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type ClimbDay struct{}

func NewClimbDay() *ClimbDay {
	return &ClimbDay{}
}

func (c *ClimbDay) Get(ctx context.Context, id int) (*model.ClimbDay, error) {
	day, err := queries(ctx).ClimbDayGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day %d | %w", id, err)
	}

	return model.ClimbDayModel(day), nil
}

func (c *ClimbDay) GetByExternalSource(ctx context.Context, source model.ClimbSource, externalID string) (*model.ClimbDay, error) {
	day, err := queries(ctx).ClimbDayGetByExternalSource(ctx, sqlc.ClimbDayGetByExternalSourceParams{
		ExternalID: externalID,
		Source:     sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day by external id %s and source %s | %w", externalID, source, err)
	}

	return model.ClimbDayModel(day), nil
}

func (c *ClimbDay) GetPopulated(ctx context.Context, id int) (*model.ClimbDay, error) {
	rows, err := queries(ctx).ClimbDayGetPopulated(ctx, int32(id))
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
	day.Gym = *model.ClimbGymModel(rows[0].ClimbGym)

	for _, row := range rows {
		climb := model.ClimbPopulatedModel(row.Climb)
		day.Climbs = append(day.Climbs, *climb)
	}

	return day, nil
}

func (c *ClimbDay) GetPopulatedByExternalSource(ctx context.Context, source model.ClimbSource, externalID string) (*model.ClimbDay, error) {
	rows, err := queries(ctx).ClimbDayGetPopulatedByExternalSource(ctx, sqlc.ClimbDayGetPopulatedByExternalSourceParams{
		ExternalID: externalID,
		Source:     sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day populated by external id %s and source %s | %w", externalID, source, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	day := model.ClimbDayModel(rows[0].ClimbDay)
	day.Gym = *model.ClimbGymModel(rows[0].ClimbGym)

	for _, row := range rows {
		climb := model.ClimbPopulatedModel(row.Climb)
		day.Climbs = append(day.Climbs, *climb)
	}

	return day, nil
}

func (c *ClimbDay) GetAllPopulatedFiltered(ctx context.Context, filter model.ClimbDayFilter) ([]*model.ClimbDay, error) {
	rows, err := queries(ctx).ClimbDayGetAllPopulatedFiltered(ctx, sqlc.ClimbDayGetAllPopulatedFilteredParams{
		UserID: int32(filter.UserID),
		Limit:  int32(filter.Limit),
		Offset: int32(filter.Offset),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day all populated filtered %+v | %w", filter, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	dayMap := make(map[int]*model.ClimbDay)

	for _, row := range rows {
		day, ok := dayMap[int(row.ClimbDay.ID)]
		if !ok {
			day = model.ClimbDayModel(row.ClimbDay)
			day.Gym = *model.ClimbGymModel(row.ClimbGym)
		}

		day.Climbs = append(day.Climbs, *model.ClimbPopulatedModel(row.Climb))
		dayMap[day.ID] = day
	}

	days := utils.MapValues(dayMap)
	slices.SortFunc(days, func(a, b *model.ClimbDay) int { return a.Date.Compare(b.Date) })

	return days, nil
}

func (c *ClimbDay) GetAllPopulatedByExternalSource(ctx context.Context, source model.ClimbSource, externalIDs []int) ([]*model.ClimbDay, error) {
	rows, err := queries(ctx).ClimbDayGetAllPopulatedByExternalSource(ctx, sqlc.ClimbDayGetAllPopulatedByExternalSourceParams{
		Column1: utils.SliceMap(externalIDs, func(id int) int32 { return int32(id) }),
		Source:  sqlc.ClimbSource(source),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all climb days populated by external ids %v and source %s | %w", externalIDs, source, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	dayMap := make(map[int]*model.ClimbDay)

	for _, row := range rows {
		day, ok := dayMap[int(row.ClimbDay.ID)]
		if !ok {
			day = model.ClimbDayModel(row.ClimbDay)
			day.Gym = *model.ClimbGymModel(row.ClimbGym)
		}

		if climb := model.ClimbPopulatedModel(row.Climb); climb != nil {
			day.Climbs = append(day.Climbs, *climb)
		}

		dayMap[day.ID] = day
	}

	days := utils.MapValues(dayMap)
	slices.SortFunc(days, func(a, b *model.ClimbDay) int { return a.Date.Compare(b.Date) })

	return days, nil
}

func (c *ClimbDay) GetAllPopulatedByUser(ctx context.Context, userID int) ([]*model.ClimbDay, error) {
	rows, err := queries(ctx).ClimbDayGetAllPopulatedByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get climb day all populated by user %d | %w", userID, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	dayMap := make(map[int]*model.ClimbDay)

	for _, row := range rows {
		day, ok := dayMap[int(row.ClimbDay.ID)]
		if !ok {
			day = model.ClimbDayModel(row.ClimbDay)
			day.Gym = *model.ClimbGymModel(row.ClimbGym)
		}

		day.Climbs = append(day.Climbs, *model.ClimbPopulatedModel(row.Climb))
		dayMap[day.ID] = day
	}

	days := utils.MapValues(dayMap)
	slices.SortFunc(days, func(a, b *model.ClimbDay) int { return a.Date.Compare(b.Date) })

	return days, nil
}

func (c *ClimbDay) Create(ctx context.Context, day *model.ClimbDay) error {
	id, err := queries(ctx).ClimbDayCreate(ctx, sqlc.ClimbDayCreateParams{
		UserID:     int32(day.UserID),
		ExternalID: day.ExternalID,
		GymID:      int32(day.GymID),
		Date:       pgtype.Timestamptz{Time: day.Date, Valid: true},
		Source:     sqlc.ClimbSource(day.Source),
	})
	if err != nil {
		return fmt.Errorf("create climb day %+v | %w", *day, err)
	}

	day.ID = int(id)

	return nil
}

func (c *ClimbDay) Update(ctx context.Context, day model.ClimbDay) error {
	if err := queries(ctx).ClimbDayUpdate(ctx, sqlc.ClimbDayUpdateParams{
		ID:    int32(day.ID),
		GymID: int32(day.GymID),
		Date:  pgtype.Timestamptz{Time: day.Date, Valid: true},
	}); err != nil {
		return fmt.Errorf("update climb day %+v | %w", day, err)
	}

	return nil
}
