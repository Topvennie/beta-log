package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type ClimbDay struct {
	ID         int
	UserID     int
	ExternalID string
	GymID      int
	Date       time.Time
	Source     Source

	// Non db fields
	Climbs []Climb
	Gym    Gym
}

func ClimbDayModel(c sqlc.ClimbDay) *ClimbDay {
	return &ClimbDay{
		ID:         int(c.ID),
		UserID:     int(c.UserID),
		ExternalID: c.ExternalID,
		GymID:      int(c.GymID),
		Date:       c.Date.Time,
		Source:     Source(c.Source),
	}
}

type ClimbDayFilter struct {
	UserID int
	Limit  int
	Offset int
}
