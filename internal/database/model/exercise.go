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
	var deletedAt time.Time
	if e.DeletedAt.Valid {
		deletedAt = e.DeletedAt.Time
	}

	return &Exercise{
		ID:        int(e.ID),
		UserID:    int(e.UserID),
		Name:      e.Name,
		Variants:  e.Variants,
		DeletedAt: deletedAt,
	}
}
