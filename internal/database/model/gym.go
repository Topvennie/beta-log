package model

import "github.com/Topvennie/beta-log/pkg/sqlc"

type Source string

const (
	SourceToplogger Source = "toplogger"
	SourceManual    Source = "manual"
)

type Gym struct {
	ID         int
	UserID     int
	ExternalID string
	Name       string
	IconPath   string
	Source     Source
}

func GymModel(c sqlc.Gym) *Gym {
	return &Gym{
		ID:         int(c.ID),
		UserID:     int(c.UserID),
		ExternalID: c.ExternalID,
		Name:       c.Name,
		IconPath:   c.IconPath,
		Source:     Source(c.Source),
	}
}
