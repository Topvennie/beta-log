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

type Variant struct {
	repo Repository
}

func (r *Repository) NewVariant() *Variant {
	return &Variant{
		repo: *r,
	}
}

func (v *Variant) Get(ctx context.Context, id int) (*model.Variant, error) {
	variant, err := v.repo.queries(ctx).VariantGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get variant by id %d | %w", id, err)
	}

	return model.VariantModel(variant), nil
}

func (v *Variant) GetByExerciseID(ctx context.Context, exerciseID int) ([]*model.Variant, error) {
	variants, err := v.repo.queries(ctx).VariantGetByExercise(ctx, int32(exerciseID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get variants by exercise %d | %w", exerciseID, err)
	}

	return utils.SliceMap(variants, model.VariantModel), nil
}

func (v *Variant) GetByIDs(ctx context.Context, ids []int) ([]*model.Variant, error) {
	variants, err := v.repo.queries(ctx).VariantGetByIDs(ctx, utils.SliceMap(ids, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get variants by ids %v | %w", ids, err)
	}

	return utils.SliceMap(variants, model.VariantModel), nil
}

func (v *Variant) Create(ctx context.Context, variant *model.Variant) error {
	id, err := v.repo.queries(ctx).VariantCreate(ctx, sqlc.VariantCreateParams{
		ExerciseID: int32(variant.ExerciseID),
		Variant:    variant.Variant,
	})
	if err != nil {
		return fmt.Errorf("create variant %+v | %w", *variant, err)
	}

	variant.ID = int(id)

	return nil
}

func (v *Variant) Update(ctx context.Context, variant model.Variant) error {
	err := v.repo.queries(ctx).VariantUpdate(ctx, sqlc.VariantUpdateParams{
		ID:      int32(variant.ID),
		Variant: variant.Variant,
	})
	if err != nil {
		return fmt.Errorf("update variant %+v | %w", variant, err)
	}

	return nil
}

func (v *Variant) Delete(ctx context.Context, id int) error {
	err := v.repo.queries(ctx).VariantDelete(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("delete variant id %d | %w", id, err)
	}

	return nil
}

func (v *Variant) DeleteByExerciseID(ctx context.Context, exerciseID int) error {
	err := v.repo.queries(ctx).VariantDeleteByExercise(ctx, int32(exerciseID))
	if err != nil {
		return fmt.Errorf("delete variants by exercise %d | %w", exerciseID, err)
	}

	return nil
}
