package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Exercise struct {
	ID        int
	UserID    int
	Name      string
	Variants  []string
	DeletedAt time.Time
}

func ExerciseModel(e sqlc.Exercise) *Exercise {
	return &Exercise{
		ID:        int(e.ID),
		UserID:    int(e.UserID),
		Name:      e.Name,
		Variants:  e.Variants,
		DeletedAt: fromTime(e.DeletedAt),
	}
}
