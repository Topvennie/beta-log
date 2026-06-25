package api

import (
	"strconv"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/gofiber/fiber/v3"
)

type task struct {
	router fiber.Router

	task service.Task
}

func newTask(router fiber.Router, service service.Service) *task {
	api := &task{
		router: router.Group("/task"),
		task:   *service.NewTask(),
	}

	api.createRoutes()

	return api
}

func (r *task) createRoutes() {
	r.router.Get("/", r.getTasks)
	r.router.Get("/history", r.getHistory)
	r.router.Post("/start/:uid", r.start)
}

func (r *task) getTasks(c fiber.Ctx) error {
	tasks, err := r.task.GetTasks(c)
	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func (r *task) getHistory(c fiber.Ctx) error {
	uid := c.Query("uid")

	var result *model.TaskResult
	if v := c.Query("result"); v != "" {
		switch v {
		case string(model.TaskSuccess), string(model.TaskFailed):
			r := model.TaskResult(v)
			result = &r
		}
	}

	var recurring *bool
	if v := c.Query("recurring"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			recurring = &b
		}
	}

	limit := fiber.Query[int](c, "limit", 10)
	page := fiber.Query[int](c, "page", 1)
	if limit < 1 || page < 1 {
		return fiber.ErrBadRequest
	}

	tasks, err := r.task.GetHistory(c, dto.TaskFilter{
		TaskUID:   uid,
		Result:    result,
		Limit:     limit,
		Recurring: recurring,
		Offset:    (page - 1) * limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func (r *task) start(c fiber.Ctx) error {
	uid := fiber.Params[string](c, "uid")

	if err := r.task.Start(c, uid); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusAccepted)
}
