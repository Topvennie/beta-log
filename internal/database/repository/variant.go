package repository

import (
	"context"
	"fmt"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Variant struct{}

func NewVariant() *Variant {
	return &Variant{}
}

func (v *Variant) Create(ctx context.Context, variant *model.Variant) error {
	id, err := queries(ctx).VariantCreate(ctx, sqlc.VariantCreateParams{
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
	err := queries(ctx).VariantUpdate(ctx, sqlc.VariantUpdateParams{
		ID:      int32(variant.ID),
		Variant: variant.Variant,
	})
	if err != nil {
		return fmt.Errorf("update variant %+v | %w", variant, err)
	}

	return nil
}

func (v *Variant) Delete(ctx context.Context, id int) error {
	if err := queries(ctx).VariantDelete(ctx, int32(id)); err != nil {
		return fmt.Errorf("delete variant id %d | %w", id, err)
	}

	return nil
}

func (v *Variant) DeleteByExerciseID(ctx context.Context, exerciseID int) error {
	if err := queries(ctx).VariantDeleteByExercise(ctx, int32(exerciseID)); err != nil {
		return fmt.Errorf("delete variants by exercise %d | %w", exerciseID, err)
	}

	return nil
}
