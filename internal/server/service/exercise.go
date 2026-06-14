package service

import (
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type Exercise struct {
	service Service

	exercise repository.Exercise
}

func (s *Service) NewExercise() *Exercise {
	return &Exercise{
		service:  *s,
		exercise: *s.repo.NewExercise(),
	}
}

func (e *Exercise) GetAll(ctx fiber.Ctx, userID int) ([]dto.Exercise, error) {
	exercises, err := e.exercise.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.SliceMap(exercises, dto.ExerciseDTO), nil
}

func (e *Exercise) Create(ctx fiber.Ctx, exerciseCreate dto.ExerciseCreate) (dto.Exercise, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Exercise{}, err
	}

	exercise := exerciseCreate.ToModel()
	exercise.UserID = userID

	if err := e.exercise.Create(ctx, &exercise); err != nil {
		return dto.Exercise{}, err
	}

	// Refetch data

	newExercise, err := e.exercise.Get(ctx, exercise.ID)
	if err != nil {
		return dto.Exercise{}, err
	}

	return dto.ExerciseDTO(newExercise), nil
}

func (e *Exercise) Update(ctx fiber.Ctx, exerciseUpdate dto.ExerciseUpdate) (dto.Exercise, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Exercise{}, err
	}

	exercise := exerciseUpdate.ToModel()
	exercise.UserID = userID

	exerciseOld, err := e.exercise.Get(ctx, exercise.ID)
	if err != nil {
		return dto.Exercise{}, err
	}
	if exerciseOld == nil {
		return dto.Exercise{}, fiber.ErrNotFound
	}
	if exerciseOld.UserID != userID {
		return dto.Exercise{}, fiber.ErrNotFound
	}

	if err := e.exercise.Update(ctx, exercise); err != nil {
		return dto.Exercise{}, err
	}

	// Refetch data

	newExercise, err := e.exercise.Get(ctx, exercise.ID)
	if err != nil {
		return dto.Exercise{}, err
	}

	return dto.ExerciseDTO(newExercise), nil
}

func (e *Exercise) Delete(ctx fiber.Ctx, id int) error {
	userID, err := getID(ctx)
	if err != nil {
		return err
	}

	exercise, err := e.exercise.Get(ctx, id)
	if err != nil {
		return err
	}
	if exercise == nil {
		return fiber.ErrNotFound
	}
	if exercise.UserID != userID {
		return fiber.ErrNotFound
	}

	if err := e.exercise.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
