package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type SessionExercise struct {
	ID int `json:"id"`

	// Exercise
	Name    string `json:"name"`
	Variant string `json:"variant"`

	// Session Exercise
	Position  int `json:"position"`
	Sets      int `json:"sets"`
	Reps      int `json:"reps,omitzero"`
	Weight    int `json:"weight,omitzero"`
	DurationS int `json:"duration_s,omitzero"`
}

func SessionExerciseDTO(s *model.SessionExercise, e *model.Exercise) SessionExercise {
	return SessionExercise{
		ID:        s.ID,
		Name:      e.Name,
		Variant:   e.Variant,
		Position:  s.Position,
		Sets:      s.Sets,
		Reps:      s.Reps,
		Weight:    s.Weight,
		DurationS: s.DurationS,
	}
}

type SessionExerciseCreate struct {
	ExerciseID int `json:"exercise_id" validate:"required,min=1"`
	Position   int `json:"position" validate:"required,min=1"`
	Sets       int `json:"sets" validate:"required,min=1"`
	Reps       int `json:"reps"`
	Weight     int `json:"weight"`
	DurationS  int `json:"duration_s"`
}

func (s SessionExerciseCreate) ToModel() model.SessionExercise {
	return model.SessionExercise{
		ExerciseID: s.ExerciseID,
		Position:   s.Position,
		Sets:       s.Sets,
		Reps:       s.Reps,
		Weight:     s.Weight,
		DurationS:  s.DurationS,
	}
}

type SessionExerciseUpdate struct {
	ID        int `json:"id" validate:"required,min=1"`
	Position  int `json:"position" validate:"required,min=1"`
	Sets      int `json:"sets" validate:"required,min=1"`
	Reps      int `json:"reps"`
	Weight    int `json:"weight"`
	DurationS int `json:"duration_s"`
}

func (s SessionExerciseUpdate) ToModel() model.SessionExercise {
	return model.SessionExercise{
		ID:        s.ID,
		Position:  s.Position,
		Sets:      s.Sets,
		Reps:      s.Reps,
		Weight:    s.Weight,
		DurationS: s.DurationS,
	}
}
