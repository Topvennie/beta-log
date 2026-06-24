package model

import (
	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type SessionExercise struct {
	ID         int
	SessionID  int
	ExerciseID int
	VariantID  int
	Position   int
	Sets       int
	Reps       int
	Weight     int
	DurationS  int

	// Non db fields
	Exercise Exercise
	Variant  Variant
}

func SessionExerciseModel(s sqlc.SessionExercise) *SessionExercise {
	return &SessionExercise{
		ID:         int(s.ID),
		SessionID:  int(s.SessionID),
		ExerciseID: int(s.ExerciseID),
		VariantID:  fromInt(s.VariantID),
		Position:   int(s.Position),
		Sets:       int(s.Sets),
		Reps:       fromInt(s.Reps),
		Weight:     fromInt(s.Weight),
		DurationS:  fromInt(s.DurationS),
	}
}

func SessionExerciseViewModel(s sqlc.SessionExercisesView) *SessionExercise {
	return &SessionExercise{
		ID:         fromInt(s.ID),
		SessionID:  fromInt(s.SessionID),
		ExerciseID: fromInt(s.ExerciseID),
		VariantID:  fromInt(s.VariantID),
		Position:   fromInt(s.Position),
		Sets:       fromInt(s.Sets),
		Reps:       fromInt(s.Reps),
		Weight:     fromInt(s.Weight),
		DurationS:  fromInt(s.DurationS),
	}
}
