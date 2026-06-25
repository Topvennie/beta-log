package model

import (
	"time"

	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Setting struct {
	ID                         int
	UserID                     int
	ClimbToploggerUserID       string
	ClimbToploggerAuthToken    string
	ClimbToploggerRefreshToken string
	ClimbTopLoggerExpiration   time.Time
}

func SettingModel(s sqlc.Setting) *Setting {
	return &Setting{
		ID:                         int(s.ID),
		UserID:                     int(s.UserID),
		ClimbToploggerUserID:       fromString(s.ClimbToploggerUserID),
		ClimbToploggerAuthToken:    fromString(s.ClimbToploggerAuthToken),
		ClimbToploggerRefreshToken: fromString(s.ClimbToploggerRefreshToken),
		ClimbTopLoggerExpiration:   fromTime(s.ClimbToploggerExpiration),
	}
}
