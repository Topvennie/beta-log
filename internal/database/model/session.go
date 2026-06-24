package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Session struct {
	ID        int
	UserID    int
	Name      string
	DeletedAt time.Time

	// Non db fields
	Exercises []SessionExercise
}

func SessionModel(s sqlc.Session) *Session {
	return &Session{
		ID:        int(s.ID),
		UserID:    int(s.UserID),
		Name:      s.Name,
		DeletedAt: fromTime(s.DeletedAt),
	}
}
