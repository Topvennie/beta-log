package dto

import (
	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/utils"
)

type Exercise struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Variants []Variant `json:"variants,omitzero"`
}

func ExerciseDTO(e *model.Exercise) Exercise {
	variants := make([]string, len(e.Variants))
	for i, v := range e.Variants {
		variants[i] = v.Variant
	}
	return Exercise{
		ID:       e.ID,
		Name:     e.Name,
		Variants: utils.SliceMap(e.Variants, func(v model.Variant) Variant { return VariantDTO(&v) }),
	}
}

type ExerciseCreate struct {
	Name     string          `json:"name" validate:"required"`
	Variants []VariantCreate `json:"variants"`
}

func (e ExerciseCreate) ToModel() model.Exercise {
	return model.Exercise{
		Name:     e.Name,
		Variants: utils.SliceMap(e.Variants, func(e VariantCreate) model.Variant { return e.ToModel() }),
	}
}

type ExerciseUpdate struct {
	ID       int             `json:"id" validate:"required,min=1"`
	Name     string          `json:"name" validate:"required"`
	Variants []VariantUpdate `json:"variants"`
}

func (e ExerciseUpdate) ToModel() model.Exercise {
	return model.Exercise{
		ID:       e.ID,
		Name:     e.Name,
		Variants: utils.SliceMap(e.Variants, func(e VariantUpdate) model.Variant { return e.ToModel() }),
	}
}
