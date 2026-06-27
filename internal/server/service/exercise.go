package service

import (
	"context"
	"slices"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type Exercise struct {
	exercise repository.Exercise
	session  repository.Session
	variant  repository.Variant
}

func NewExercise() *Exercise {
	return &Exercise{
		exercise: *repository.NewExercise(),
		session:  *repository.NewSession(),
		variant:  *repository.NewVariant(),
	}
}

func (e *Exercise) GetAll(ctx fiber.Ctx) ([]dto.Exercise, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	exercises, err := e.exercise.GetByUserID(ctx, userID)
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

	if err := withRollback(ctx, func(ctx context.Context) error {
		if err := e.exercise.Create(ctx, &exercise); err != nil {
			return err
		}

		for _, variant := range exercise.Variants {
			variant.ExerciseID = exercise.ID
			if err := e.variant.Create(ctx, &variant); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
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

	// Make sure we don't delete a variant assigned to a session
	variantsToDelete := []int{}
	for _, variantOld := range exerciseOld.Variants {
		if idx := slices.IndexFunc(exercise.Variants, func(v model.Variant) bool { return variantOld.ID == v.ID }); idx == -1 {
			// Variant doesn't exist anymore -> delete
			variantsToDelete = append(variantsToDelete, variantOld.ID)
		}
	}

	sessions, err := e.session.GetByVariants(ctx, variantsToDelete)
	if err != nil {
		return dto.Exercise{}, err
	}
	if len(sessions) != 0 {
		return dto.Exercise{}, fiber.NewError(fiber.StatusBadRequest, "deleted variant is assigned to a session")
	}

	if err := withRollback(ctx, func(context.Context) error {
		if err := e.exercise.Update(ctx, exercise); err != nil {
			return err
		}

		for _, variant := range exercise.Variants {
			variant.ExerciseID = exercise.ID

			if idx := slices.IndexFunc(exerciseOld.Variants, func(v model.Variant) bool { return variant.ID == v.ID }); idx == -1 {
				// Variant doesn't exist yet -> create
				if err := e.variant.Create(ctx, &variant); err != nil {
					return err
				}

				continue
			}

			// Variant already exists -> update
			if err := e.variant.Update(ctx, variant); err != nil {
				return err
			}
		}

		for _, variantToDelete := range variantsToDelete {
			if err := e.variant.Delete(ctx, variantToDelete); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
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

	session, err := e.session.GetByExercise(ctx, exercise.ID)
	if err != nil {
		return err
	}
	if session != nil {
		return fiber.NewError(fiber.StatusBadRequest, "exercise is connected to session "+session.Name)
	}

	if err := withRollback(ctx, func(ctx context.Context) error {
		if err := e.variant.DeleteByExerciseID(ctx, exercise.ID); err != nil {
			return err
		}

		if err := e.exercise.Delete(ctx, id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
