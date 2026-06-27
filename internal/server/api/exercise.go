package api

import (
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type exercise struct {
	router   fiber.Router
	exercise *service.Exercise
}

func newExercise(router fiber.Router) *exercise {
	api := &exercise{
		router:   router.Group("/exercise"),
		exercise: service.NewExercise(),
	}

	api.routes()

	return api
}

func (e *exercise) routes() {
	e.router.Get("/", e.getAll)
	e.router.Post("/", e.create)
	e.router.Put("/:id", e.update)
	e.router.Delete("/:id", e.delete)
}

func (e *exercise) getAll(c fiber.Ctx) error {
	exercises, err := e.exercise.GetAll(c)
	if err != nil {
		return err
	}

	return c.JSON(exercises)
}

func (e *exercise) create(c fiber.Ctx) error {
	var exercise dto.ExerciseCreate
	if err := c.Bind().Body(&exercise); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(exercise); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newExercise, err := e.exercise.Create(c, exercise)
	if err != nil {
		return err
	}

	return c.JSON(newExercise)
}

func (e *exercise) update(c fiber.Ctx) error {
	var exercise dto.ExerciseUpdate
	if err := c.Bind().Body(&exercise); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(exercise); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id := fiber.Params[int](c, "id")
	if id != exercise.ID {
		return fiber.NewError(fiber.StatusBadRequest, "params id doesn't match body id")
	}

	newExercise, err := e.exercise.Update(c, exercise)
	if err != nil {
		return err
	}

	return c.JSON(newExercise)
}

func (e *exercise) delete(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	if id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := e.exercise.Delete(c, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
