package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/sqlc"
)

type Setting struct{}

func NewSetting() *Setting {
	return &Setting{}
}

func (s *Setting) GetByUser(ctx context.Context, userID int) (*model.Setting, error) {
	setting, err := queries(ctx).SettingGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get setting by user id %d | %w", userID, err)
	}

	return model.SettingModel(setting), nil
}

func (s *Setting) Create(ctx context.Context, setting *model.Setting) error {
	id, err := queries(ctx).SettingCreate(ctx, int32(setting.UserID))
	if err != nil {
		return fmt.Errorf("create setting %+v | %w", *setting, err)
	}

	setting.ID = int(id)

	return nil
}

func (s *Setting) ToploggerUpdate(ctx context.Context, setting model.Setting) error {
	if err := queries(ctx).SettingToploggerUpdate(ctx, sqlc.SettingToploggerUpdateParams{
		ID:                         int32(setting.ID),
		ClimbToploggerUserID:       toString(setting.ClimbToploggerUserID),
		ClimbToploggerAuthToken:    toString(setting.ClimbToploggerAuthToken),
		ClimbToploggerRefreshToken: toString(setting.ClimbToploggerRefreshToken),
		ClimbToploggerExpiration:   toTime(setting.ClimbTopLoggerExpiration),
	}); err != nil {
		return fmt.Errorf("update setting %+v | %w", setting, err)
	}

	return nil
}
