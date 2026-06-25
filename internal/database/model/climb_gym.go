package model

import "github.com/Topvennie/beta-log/pkg/sqlc"

type ClimbGym struct {
	ID         int
	UserID     int
	ExternalID string
	Name       string
	IconPath   string
}

func ClimbGymModel(c sqlc.ClimbGym) *ClimbGym {
	return &ClimbGym{
		ID:         int(c.ID),
		UserID:     int(c.UserID),
		ExternalID: c.ExternalID,
		Name:       c.Name,
		IconPath:   c.IconPath,
	}
}
