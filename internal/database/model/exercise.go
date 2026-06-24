package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Exercise struct {
	ID        int
	UserID    int
	Name      string
	DeletedAt time.Time

	// Non db fields
	Variants []Variant
}

func ExerciseModel(e sqlc.Exercise) *Exercise {
	return &Exercise{
		ID:        int(e.ID),
		UserID:    int(e.UserID),
		Name:      e.Name,
		DeletedAt: fromTime(e.DeletedAt),
	}
}

func ExerciseViewModel(e sqlc.ExercisesView) *Exercise {
	return &Exercise{
		ID:        fromInt(e.ID),
		UserID:    fromInt(e.UserID),
		Name:      fromString(e.Name),
		DeletedAt: fromTime(e.DeletedAt),
	}
}
