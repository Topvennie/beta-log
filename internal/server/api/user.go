package api

import (
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type user struct {
	router fiber.Router
	user   service.User
}

func newUser(router fiber.Router, service service.Service) *user {
	api := &user{
		router: router.Group("/user"),
		user:   *service.NewUser(),
	}

	api.routes()

	return api
}

func (u *user) routes() {
	u.router.Get("/me", u.getMeHandler)
}

func (u *user) getMeHandler(c fiber.Ctx) error {
	id, ok := c.Locals("id").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, err := u.user.GetByID(c.RequestCtx(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
