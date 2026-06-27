package api

import (
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type session struct {
	router  fiber.Router
	session *service.Session
}

func newSession(router fiber.Router) *session {
	api := &session{
		router:  router.Group("/session"),
		session: service.NewSession(),
	}

	api.routes()

	return api
}

func (s *session) routes() {
	s.router.Get("/", s.getAll)
	s.router.Post("/", s.create)
	s.router.Put("/:id", s.update)
	s.router.Delete("/:id", s.delete)
}

func (s *session) getAll(c fiber.Ctx) error {
	sessions, err := s.session.GetAll(c)
	if err != nil {
		return err
	}

	return c.JSON(sessions)
}

func (s *session) create(c fiber.Ctx) error {
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

func (s *session) update(c fiber.Ctx) error {
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

func (s *session) delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := s.session.Delete(c, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
