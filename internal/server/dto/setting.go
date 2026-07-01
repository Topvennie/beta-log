package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type Setting struct {
	ClimbToploggerUserID       string `json:"climb_toplogger_user_id,omitzero"`
	ClimbToploggerAuthToken    string `json:"climb_toplogger_auth_token,omitzero"`
	ClimbToploggerRefreshToken string `json:"climb_toplogger_refresh_token,omitzero"`
}

func SettingDTO(s *model.Setting) Setting {
	return Setting{
		ClimbToploggerUserID:       s.ClimbToploggerUserID,
		ClimbToploggerAuthToken:    s.ClimbToploggerAuthToken,
		ClimbToploggerRefreshToken: s.ClimbToploggerRefreshToken,
	}
}

type SettingToploggerUpdate struct {
	ClimbToploggerUserID       string `json:"climb_toplogger_user_id" validate:"required_with=ClimbToploggerAuthToken ClimbToploggerRefreshToken"`
	ClimbToploggerAuthToken    string `json:"climb_toplogger_auth_token" validate:"required_with=ClimbToploggerUserID ClimbToploggerRefreshToken"`
	ClimbToploggerRefreshToken string `json:"climb_toplogger_refresh_token" validate:"required_with=ClimbToploggerUserID ClimbToploggerAuthToken"`
}

func (s SettingToploggerUpdate) ToModel() model.Setting {
	return model.Setting{
		ClimbToploggerUserID:       s.ClimbToploggerUserID,
		ClimbToploggerAuthToken:    s.ClimbToploggerAuthToken,
		ClimbToploggerRefreshToken: s.ClimbToploggerRefreshToken,
	}
}
