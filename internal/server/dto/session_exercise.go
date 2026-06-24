package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type SessionExercise struct {
	ID int `json:"id"`

	Exercise Exercise `json:"exercise"`
	Variant  Variant  `json:"variant,omitzero"`

	Position  int `json:"position"`
	Sets      int `json:"sets"`
	Reps      int `json:"reps,omitzero"`
	Weight    int `json:"weight,omitzero"`
	DurationS int `json:"duration_s,omitzero"`
}

func SessionExerciseDTO(s *model.SessionExercise) SessionExercise {
	return SessionExercise{
		ID:        s.ID,
		Exercise:  ExerciseDTO(&s.Exercise),
		Variant:   VariantDTO(&s.Variant),
		Position:  s.Position,
		Sets:      s.Sets,
		Reps:      s.Reps,
		Weight:    s.Weight,
		DurationS: s.DurationS,
	}
}

type SessionExerciseCreate struct {
	ExerciseID int `json:"exercise_id" validate:"required,min=1"`
	VariantID  int `json:"variant_id"`
	Position   int `json:"position" validate:"required,min=1"`
	Sets       int `json:"sets" validate:"required,min=1"`
	Reps       int `json:"reps" validate:"min=0"`
	Weight     int `json:"weight" validate:"min=0"`
	DurationS  int `json:"duration_s" validate:"min=0"`
}

func (s SessionExerciseCreate) ToModel() model.SessionExercise {
	return model.SessionExercise{
		ExerciseID: s.ExerciseID,
		VariantID:  s.VariantID,
		Position:   s.Position,
		Sets:       s.Sets,
		Reps:       s.Reps,
		Weight:     s.Weight,
		DurationS:  s.DurationS,
	}
}

type SessionExerciseUpdate struct {
	ExerciseID int `json:"exercise_id" validate:"required,min=1"`
	VariantID  int `json:"variant_id"`
	Position   int `json:"position" validate:"required,min=1"`
	Sets       int `json:"sets" validate:"required,min=1"`
	Reps       int `json:"reps" validate:"min=0"`
	Weight     int `json:"weight" validate:"min=0"`
	DurationS  int `json:"duration_s" validate:"min=0"`
}

func (s SessionExerciseUpdate) ToModel() model.SessionExercise {
	return model.SessionExercise{
		ExerciseID: s.ExerciseID,
		VariantID:  s.VariantID,
		Position:   s.Position,
		Sets:       s.Sets,
		Reps:       s.Reps,
		Weight:     s.Weight,
		DurationS:  s.DurationS,
	}
}
