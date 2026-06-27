package api

import (
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type setting struct {
	router fiber.Router

	setting service.Setting
}

func newSetting(router fiber.Router, service service.Service) *setting {
	api := &setting{
		router:  router.Group("/setting"),
		setting: *service.NewSetting(),
	}

	api.createRoutes()

	return api
}

func (s *setting) createRoutes() {
	s.router.Get("/", s.get)
	s.router.Post("/", s.update)
}

func (s *setting) get(c fiber.Ctx) error {
	setting, err := s.setting.Get(c)
	if err != nil {
		return err
	}

	return c.JSON(setting)
}

func (s *setting) update(c fiber.Ctx) error {
	var setting dto.SettingUpdate
	if err := c.Bind().Body(&setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newSetting, err := s.setting.Update(c, setting)
	if err != nil {
		return err
	}

	return c.JSON(newSetting)
}
