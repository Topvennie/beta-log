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
	Day        time.Time

	// Non db fields
	Climbs []Climb
	Gym    *ClimbGym
}

func ClimbDayModel(c sqlc.ClimbDay) *ClimbDay {
	return &ClimbDay{
		ID:         int(c.ID),
		UserID:     int(c.UserID),
		ExternalID: c.ExternalID,
		GymID:      int(c.GymID),
		Day:        c.Day.Time,
	}
}
