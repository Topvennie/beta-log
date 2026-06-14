package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Session struct {
	ID        int
	UserID    int
	Name      string
	Active    bool
	Position  int
	DeletedAt time.Time

	// Non db fields
	// Not guaranteed to be populated
	Exercises []SessionExercise
}

func SessionModel(s sqlc.Session) *Session {
	var deletedAt time.Time
	if s.DeletedAt.Valid {
		deletedAt = s.DeletedAt.Time
	}

	return &Session{
		ID:        int(s.ID),
		UserID:    int(s.UserID),
		Name:      s.Name,
		Active:    s.Active,
		Position:  int(s.Position.Int32),
		DeletedAt: deletedAt,
	}
}
