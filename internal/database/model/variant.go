package model

import "github.com/Topvennie/beta-log/pkg/sqlc"

type Variant struct {
	ID         int
	ExerciseID int
	Variant    string
}

func VariantModel(e sqlc.Variant) *Variant {
	return &Variant{
		ID:         int(e.ID),
		ExerciseID: int(e.ExerciseID),
		Variant:    e.Variant,
	}
}

func VariantViewModel(e sqlc.VariantsView) *Variant {
	return &Variant{
		ID:         fromInt(e.ID),
		ExerciseID: fromInt(e.ExerciseID),
		Variant:    fromString(e.Variant),
	}
}
