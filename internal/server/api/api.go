package api

import (
	"github.com/Topvennie/beta-log/internal/server/middlewares"
	"github.com/gofiber/fiber/v3"
)

func New(router fiber.Router) error {
	// Authentication
	if _, err := newAuth(router); err != nil {
		return err
	}

	// Protected routes
	routerProtected := router.Group("/auth", middlewares.AuthRoute)

	newUser(routerProtected)
	newExercise(routerProtected)
	newSession(routerProtected)
	newTask(routerProtected)
	newSetting(routerProtected)

	return nil
}
