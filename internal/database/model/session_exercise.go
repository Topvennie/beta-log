package model

import (
	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type SessionExercise struct {
	ID         int
	SessionID  int
	ExerciseID int
	Position   int
	Sets       int
	Reps       int
	Weight     int
	DurationS  int
}

func SessionExerciseModel(s sqlc.SessionsExercise) *SessionExercise {
	var reps int
	if s.Reps.Valid {
		r := int(s.Reps.Int32)
		reps = r
	}

	var weight int
	if s.Weight.Valid {
		w := int(s.Weight.Int32)
		weight = w
	}

	var durationS int
	if s.DurationS.Valid {
		d := int(s.DurationS.Int32)
		durationS = d
	}

	return &SessionExercise{
		ID:         int(s.ID),
		SessionID:  int(s.SessionID),
		ExerciseID: int(s.ExerciseID),
		Position:   int(s.Position),
		Sets:       int(s.Sets),
		Reps:       reps,
		Weight:     weight,
		DurationS:  durationS,
	}
}
