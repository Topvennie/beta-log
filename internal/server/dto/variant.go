package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type Variant struct {
	ID      int    `json:"id"`
	Variant string `json:"variant"`
}

func VariantDTO(v *model.Variant) Variant {
	return Variant{
		ID:      v.ID,
		Variant: v.Variant,
	}
}

type VariantCreate struct {
	Variant string `json:"variant" validate:"required"`
}

func (e VariantCreate) ToModel() model.Variant {
	return model.Variant{
		Variant: e.Variant,
	}
}

type VariantUpdate struct {
	ID      int    `json:"id" validate:"required,min=1"`
	Variant string `json:"variant" validate:"required"`
}

func (e VariantUpdate) ToModel() model.Variant {
	return model.Variant{
		ID:      e.ID,
		Variant: e.Variant,
	}
}
