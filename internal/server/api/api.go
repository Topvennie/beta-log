package api

import (
	"github.com/Topvennie/beta-log/internal/server/middlewares"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

func New(router fiber.Router, service service.Service) error {
	// Authentication
	if _, err := newAuth(router, service); err != nil {
		return err
	}

	// Protected routes
	routerProtected := router.Group("/auth", middlewares.AuthRoute)

	newUser(routerProtected, service)
	newExercise(routerProtected, service)
	newSession(routerProtected, service)
	newTask(routerProtected, service)
	newSetting(routerProtected, service)

	return nil
}
