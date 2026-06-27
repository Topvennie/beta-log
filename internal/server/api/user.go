package api

import (
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type user struct {
	router fiber.Router
	user   *service.User
}

func newUser(router fiber.Router) *user {
	api := &user{
		router: router.Group("/user"),
		user:   service.NewUser(),
	}

	api.routes()

	return api
}

func (u *user) routes() {
	u.router.Get("/me", u.getMe)
}

func (u *user) getMe(c fiber.Ctx) error {
	user, err := u.user.GetMe(c)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
