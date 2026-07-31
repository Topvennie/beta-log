package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type GradeSystem string

const (
	GradeSystemFont GradeSystem = "font"
	GradeSystemV    GradeSystem = "v"
)

type Setting struct {
	ID                         int
	UserID                     int
	ClimbToploggerUserID       string
	ClimbToploggerAuthToken    string
	ClimbToploggerRefreshToken string
	ClimbToploggerExpiration   time.Time
	GradeSystem                GradeSystem
}

func SettingModel(s sqlc.Setting) *Setting {
	return &Setting{
		ID:                         int(s.ID),
		UserID:                     int(s.UserID),
		ClimbToploggerUserID:       fromString(s.ClimbToploggerUserID),
		ClimbToploggerAuthToken:    fromString(s.ClimbToploggerAuthToken),
		ClimbToploggerRefreshToken: fromString(s.ClimbToploggerRefreshToken),
		ClimbToploggerExpiration:   fromTime(s.ClimbToploggerExpiration),
		GradeSystem:                GradeSystem(s.GradeSystem),
	}
}
