package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/fetch"
	"github.com/Topvennie/beta-log/internal/fetch/toplogger"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/task"
	"github.com/gofiber/fiber/v3"
)

type Setting struct {
	setting repository.Setting
	user    repository.User
}

func NewSetting() *Setting {
	return &Setting{
		setting: *repository.NewSetting(),
		user:    *repository.NewUser(),
	}
}

func (s *Setting) Get(ctx fiber.Ctx) (dto.Setting, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Setting{}, err
	}

	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}
	if setting == nil {
		return dto.Setting{}, fmt.Errorf("user %d has no settings", userID)
	}

	return dto.SettingDTO(setting), nil
}

func (s *Setting) UpdateGradeSystem(ctx fiber.Ctx, settingSave dto.SettingUpdateGradeSystem) (dto.Setting, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Setting{}, err
	}
	user, err := s.user.GetByID(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}
	if user == nil {
		return dto.Setting{}, fmt.Errorf("user %d not found", userID)
	}

	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}
	if setting == nil {
		return dto.Setting{}, fmt.Errorf("user %d has no settings", userID)
	}

	setting.GradeSystem = settingSave.GradeSystem

	if err := s.setting.UpdateGradeSystem(ctx, *setting); err != nil {
		return dto.Setting{}, err
	}

	// Refetch data
	newSetting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}

	return dto.SettingDTO(newSetting), nil
}

func (s *Setting) UpdateToplogger(ctx fiber.Ctx, settingSave dto.SettingUpdateToplogger) (dto.Setting, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Setting{}, err
	}
	user, err := s.user.GetByID(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}
	if user == nil {
		return dto.Setting{}, fmt.Errorf("user %d not found", userID)
	}

	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}
	if setting == nil {
		return dto.Setting{}, fmt.Errorf("user %d has no settings", userID)
	}

	setting.ClimbToploggerUserID = settingSave.ClimbToploggerUserID
	setting.ClimbToploggerAuthToken = settingSave.ClimbToploggerAuthToken
	setting.ClimbToploggerRefreshToken = settingSave.ClimbToploggerRefreshToken
	setting.ClimbToploggerExpiration = time.Time{} // Use this as placeholder

	if setting.ClimbToploggerUserID != "" {
		// The user provded data
		// Let's try it out

		tokens, err := toplogger.New().Refresh(ctx, *setting)
		if err != nil {
			if errors.Is(err, toplogger.ErrUnauthorized) || errors.Is(err, toplogger.ErrNoTokenResponse) {
				return dto.Setting{}, fiber.NewError(fiber.StatusBadRequest, "invalid user id / tokens")
			}
			return dto.Setting{}, err
		}

		// Save the new tokens
		setting.ClimbToploggerAuthToken = tokens.Access.Token
		setting.ClimbToploggerRefreshToken = tokens.Refresh.Token
		setting.ClimbToploggerExpiration = tokens.Refresh.ExpiresAt
	}

	if err := s.setting.UpdateToplogger(ctx, *setting); err != nil {
		return dto.Setting{}, err
	}

	if setting.ClimbToploggerUserID != "" {
		// Start the climb update task
		if err := task.Manager.RunRecurringByUID(fetch.TaskUpdateUID, *user); err != nil {
			return dto.Setting{}, err
		}
	}

	// Refetch data
	newSetting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}

	return dto.SettingDTO(newSetting), nil
}
