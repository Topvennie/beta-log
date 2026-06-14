package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Exercise struct {
	ID        int
	UserID    int
	Name      string
	Variant   string
	DeletedAt time.Time
}

func ExerciseModel(e sqlc.Exercise) *Exercise {
	var variant string
	if e.Variant.Valid {
		variant = e.Variant.String
	}

	var deletedAt time.Time
	if e.DeletedAt.Valid {
		deletedAt = e.DeletedAt.Time
	}

	return &Exercise{
		ID:        int(e.ID),
		UserID:    int(e.UserID),
		Name:      e.Name,
		Variant:   variant,
		DeletedAt: deletedAt,
	}
}
