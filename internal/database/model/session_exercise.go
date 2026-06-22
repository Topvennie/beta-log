package model

import (
	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type SessionExercise struct {
	ID         int
	SessionID  int
	ExerciseID int
	Variant    string
	Position   int
	Sets       int
	Reps       int
	Weight     int
	DurationS  int
}

func SessionExerciseModel(s sqlc.SessionsExercise) *SessionExercise {
	return &SessionExercise{
		ID:         int(s.ID),
		SessionID:  int(s.SessionID),
		ExerciseID: int(s.ExerciseID),
		Variant:    fromString(s.Variant),
		Position:   int(s.Position),
		Sets:       int(s.Sets),
		Reps:       fromInt(s.Reps),
		Weight:     fromInt(s.Weight),
		DurationS:  fromInt(s.DurationS),
	}
}
