package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type Setting struct {
	ID                         int    `json:"id"`
	ClimbToploggerUserID       string `json:"climb_toplogger_user_id,omitzero"`
	ClimbToploggerAuthToken    string `json:"climb_toplogger_auth_token,omitzero"`
	ClimbToploggerRefreshToken string `json:"climb_toplogger_refresh_token,omitzero"`
}

func SettingDTO(s *model.Setting) Setting {
	return Setting{
		ID:                         s.ID,
		ClimbToploggerUserID:       s.ClimbToploggerUserID,
		ClimbToploggerAuthToken:    s.ClimbToploggerAuthToken,
		ClimbToploggerRefreshToken: s.ClimbToploggerRefreshToken,
	}
}

type SettingUpdate struct {
	ID                         int    `json:"id" validate:"required,min=1"`
	ClimbToploggerUserID       string `json:"climb_toplogger_user_id"`
	ClimbToploggerAuthToken    string `json:"climb_toplogger_auth_token" validate:"required_with=ClimbToploggerAuthToken"`
	ClimbToploggerRefreshToken string `json:"climb_toplogger_refresh_token" validate:"required_with=ClimbToploggerRefreshToken"`
}

func (s SettingUpdate) ToModel() model.Setting {
	return model.Setting{
		ID:                         s.ID,
		ClimbToploggerUserID:       s.ClimbToploggerUserID,
		ClimbToploggerAuthToken:    s.ClimbToploggerAuthToken,
		ClimbToploggerRefreshToken: s.ClimbToploggerRefreshToken,
	}
}
