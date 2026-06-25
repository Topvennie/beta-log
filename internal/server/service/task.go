package service

import (
	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/task"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type Task struct {
	service Service

	task repository.Task
	user repository.User
}

func (s *Service) NewTask() *Task {
	return &Task{
		service: *s,
		task:    *s.repo.NewTask(),
		user:    *s.repo.NewUser(),
	}
}

func (t *Task) GetTasks(ctx fiber.Ctx) ([]dto.Task, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := task.Manager.Tasks()
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if tasks == nil {
		return []dto.Task{}, nil
	}

	lastRuns, err := t.task.GetRunLastAllByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	lastRunMap := make(map[string]*model.Task)
	for _, lastRun := range lastRuns {
		lastRunMap[lastRun.UID] = lastRun
	}

	taskDTOs := make([]dto.Task, 0, len(tasks))
	for _, task := range tasks {
		taskDTO := dto.TaskDTO(task)

		if lastRun, ok := lastRunMap[task.TaskUID]; ok {
			lastError := ""
			if lastRun.Error != nil {
				lastError = lastRun.Error.Error()
			}

			taskDTO.LastStatus = lastRun.Result
			taskDTO.LastMessage = lastRun.Message
			taskDTO.LastError = lastError
		}

		taskDTOs = append(taskDTOs, taskDTO)
	}

	return taskDTOs, nil
}

func (t *Task) GetHistory(ctx fiber.Ctx, filter dto.TaskFilter) ([]dto.TaskHistory, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}
	filter.UserID = userID

	tasks, err := t.task.GetRunFiltered(ctx, *filter.ToModel())
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if tasks == nil {
		return []dto.TaskHistory{}, nil
	}

	return utils.SliceMap(tasks, dto.TaskHistoryDTO), nil
}

func (t *Task) Start(ctx fiber.Ctx, taskUID string) error {
	userID, err := getID(ctx)
	if err != nil {
		return err
	}

	user, err := t.user.GetByID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrUnauthorized
	}

	taskModel, err := t.task.GetByUID(ctx, taskUID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if taskModel == nil {
		return fiber.ErrNotFound
	}
	if !taskModel.Recurring || !taskModel.Active {
		return fiber.ErrBadRequest
	}

	return task.Manager.RunRecurringByUID(taskUID, *user)
}
