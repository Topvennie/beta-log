package api

import (
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type sessionAPI struct {
	router  fiber.Router
	session service.Session
}

func newSession(router fiber.Router, service service.Service) *sessionAPI {
	api := &sessionAPI{
		router:  router.Group("/session"),
		session: *service.NewSession(),
	}

	api.createRoutes()

	return api
}

func (s *sessionAPI) createRoutes() {
	s.router.Get("/", s.getAll)
	s.router.Put("/", s.create)
	s.router.Post("/:id", s.update)
	s.router.Delete("/:id", s.delete)
}

func (s *sessionAPI) getAll(c fiber.Ctx) error {
	id, ok := c.Locals("id").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	sessions, err := s.session.GetAll(c, id)
	if err != nil {
		return err
	}

	return c.JSON(sessions)
}

func (s *sessionAPI) create(c fiber.Ctx) error {
	var session dto.SessionCreate
	if err := c.Bind().Body(&session); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(session); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newSession, err := s.session.Create(c, session)
	if err != nil {
		return err
	}

	return c.JSON(newSession)
}

func (s *sessionAPI) update(c fiber.Ctx) error {
	var session dto.SessionUpdate
	if err := c.Bind().Body(&session); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(session); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id := fiber.Params[int](c, "id")
	if id != session.ID {
		return fiber.NewError(fiber.StatusBadRequest, "params id doesn't match body id")
	}

	newSession, err := s.session.Update(c, session)
	if err != nil {
		return err
	}

	return c.JSON(newSession)
}

func (s *sessionAPI) delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := s.session.Delete(c, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
