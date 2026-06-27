package service

import (
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/gofiber/fiber/v3"
)

type Setting struct {
	service Service

	setting repository.Setting
}

func (s *Service) NewSetting() *Setting {
	return &Setting{
		service: *s,
		setting: *s.repo.NewSetting(),
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
		return dto.Setting{}, fiber.NewError(fiber.StatusInternalServerError, "user without settings")
	}

	return dto.SettingDTO(setting), nil
}

func (s *Setting) Update(ctx fiber.Ctx, settingSave dto.SettingUpdate) (dto.Setting, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Setting{}, err
	}

	oldSetting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}
	if oldSetting == nil || oldSetting.UserID != userID {
		return dto.Setting{}, fiber.ErrNotFound
	}

	setting := settingSave.ToModel()
	setting.UserID = userID

	if err := s.setting.Update(ctx, setting); err != nil {
		return dto.Setting{}, err
	}

	// Refetch data
	newSetting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		return dto.Setting{}, err
	}

	return dto.SettingDTO(newSetting), nil
}
