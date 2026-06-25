package model

import "github.com/Topvennie/beta-log/pkg/sqlc"

type ClimbType string

const (
	ClimbTypeBoulder ClimbType = "boulder"
	ClimbTypeLead    ClimbType = "lead"
)

type ClimbFinish string

const (
	ClimbFinishFlash  ClimbFinish = "flash"
	ClimbFinishTop    ClimbFinish = "top"
	ClimbFinishRepeat ClimbFinish = "repeat"
)

type Climb struct {
	ID         int
	UserID     int
	ExternalID string
	ClimbDayID int
	Grade      int
	Color      string
	HoldColor  string
	ClimbType  ClimbType
	FinishType ClimbFinish
}

func ClimbModel(c sqlc.Climb) *Climb {
	return &Climb{
		ID:         int(c.ID),
		UserID:     int(c.UserID),
		ExternalID: c.ExternalID,
		ClimbDayID: int(c.ClimbDayID),
		Grade:      int(c.Grade),
		Color:      c.Color,
		HoldColor:  c.HoldColor,
		ClimbType:  ClimbType(c.ClimbType),
		FinishType: ClimbFinish(c.FinishType),
	}
}

func ClimbPopulatedModel(c sqlc.Climb) *Climb {
	if c.ID == 0 {
		return nil
	}
	return ClimbModel(c)
}
