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

type Exercise struct {
	repo Repository
}

func (r *Repository) NewExercise() *Exercise {
	return &Exercise{
		repo: *r,
	}
}

func (e *Exercise) Get(ctx context.Context, id int) (*model.Exercise, error) {
	rows, err := e.repo.queries(ctx).ExerciseGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get exercise populated by id %d | %w", id, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	exercise := model.ExerciseModel(rows[0].Exercise)

	for _, row := range rows {
		if !row.VariantsView.ID.Valid {
			continue
		}

		exercise.Variants = append(exercise.Variants, *model.VariantViewModel(row.VariantsView))
	}

	return exercise, nil
}

func (e *Exercise) GetByUserID(ctx context.Context, userID int) ([]*model.Exercise, error) {
	rows, err := e.repo.queries(ctx).ExerciseGetAll(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all exercises with variants for user %d | %w", userID, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	exerciseMap := make(map[int]*model.Exercise)

	for _, row := range rows {
		exercise, ok := exerciseMap[int(row.Exercise.ID)]
		if !ok {
			exercise = model.ExerciseModel(row.Exercise)
		}

		if row.VariantsView.ID.Valid {
			exercise.Variants = append(exercise.Variants, *model.VariantViewModel(row.VariantsView))
		}

		exerciseMap[int(row.Exercise.ID)] = exercise
	}

	return utils.MapValues(exerciseMap), nil
}

func (e *Exercise) GetByIDs(ctx context.Context, ids []int) ([]*model.Exercise, error) {
	rows, err := e.repo.queries(ctx).ExerciseGetByIDs(ctx, utils.SliceMap(ids, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get exercises by ids %v | %w", ids, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	exerciseMap := make(map[int]*model.Exercise)

	for _, row := range rows {
		exercise, ok := exerciseMap[int(row.Exercise.ID)]
		if !ok {
			exercise = model.ExerciseModel(row.Exercise)
		}

		if row.VariantsView.ID.Valid {
			exercise.Variants = append(exercise.Variants, *model.VariantViewModel(row.VariantsView))
		}

		exerciseMap[int(row.Exercise.ID)] = exercise
	}

	return utils.MapValues(exerciseMap), nil
}

func (e *Exercise) Create(ctx context.Context, exercise *model.Exercise) error {
	id, err := e.repo.queries(ctx).ExerciseCreate(ctx, sqlc.ExerciseCreateParams{
		UserID: int32(exercise.UserID),
		Name:   exercise.Name,
	})
	if err != nil {
		return fmt.Errorf("create exercise %+v | %w", *exercise, err)
	}

	exercise.ID = int(id)

	return nil
}

func (e *Exercise) Update(ctx context.Context, exercise model.Exercise) error {
	err := e.repo.queries(ctx).ExerciseUpdate(ctx, sqlc.ExerciseUpdateParams{
		ID:   int32(exercise.ID),
		Name: exercise.Name,
	})
	if err != nil {
		return fmt.Errorf("update exercise %+v | %w", exercise, err)
	}

	return nil
}

func (e *Exercise) Delete(ctx context.Context, id int) error {
	err := e.repo.queries(ctx).ExerciseDelete(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("delete exercise with id %d | %w", id, err)
	}

	return nil
}
