package api

import (
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type setting struct {
	router fiber.Router

	setting *service.Setting
}

func newSetting(router fiber.Router) *setting {
	api := &setting{
		router:  router.Group("/setting"),
		setting: service.NewSetting(),
	}

	api.routes()

	return api
}

func (s *setting) routes() {
	s.router.Get("/", s.get)
	s.router.Put("/grade_system", s.updateGradeSystem)
	s.router.Put("/toplogger", s.updateToplogger)
}

func (s *setting) get(c fiber.Ctx) error {
	setting, err := s.setting.Get(c)
	if err != nil {
		return err
	}

	return c.JSON(setting)
}

func (s *setting) updateGradeSystem(c fiber.Ctx) error {
	var setting dto.SettingUpdateGradeSystem
	if err := c.Bind().Body(&setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newSetting, err := s.setting.UpdateGradeSystem(c, setting)
	if err != nil {
		return err
	}

	return c.JSON(newSetting)
}

func (s *setting) updateToplogger(c fiber.Ctx) error {
	var setting dto.SettingUpdateToplogger
	if err := c.Bind().Body(&setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newSetting, err := s.setting.UpdateToplogger(c, setting)
	if err != nil {
		return err
	}

	return c.JSON(newSetting)
}
