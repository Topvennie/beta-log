package service

import (
	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Gym struct {
	gym repository.Gym
}

func NewGym() *Gym {
	return &Gym{
		gym: *repository.NewGym(),
	}
}

func (g *Gym) GetAll(ctx fiber.Ctx) ([]dto.Gym, error) {
	userID, err := getID(ctx)
	if err != nil {
		return nil, err
	}

	gyms, err := g.gym.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.SliceMap(gyms, dto.GymDTO), nil
}

func (g *Gym) Create(ctx fiber.Ctx, gymSave dto.GymCreate) (dto.Gym, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Gym{}, err
	}

	gym := gymSave.ToModel()
	gym.UserID = userID
	gym.ExternalID = uuid.NewString()
	gym.Source = model.SourceManual

	if err := g.gym.Create(ctx, &gym); err != nil {
		return dto.Gym{}, err
	}

	return dto.GymDTO(&gym), nil
}

func (g *Gym) Update(ctx fiber.Ctx, gymSave dto.GymUpdate) (dto.Gym, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Gym{}, err
	}

	oldGym, err := g.gym.Get(ctx, gymSave.ID)
	if err != nil {
		return dto.Gym{}, err
	}
	if oldGym == nil || oldGym.UserID != userID {
		return dto.Gym{}, fiber.ErrNotFound
	}
	if oldGym.Source != model.SourceManual {
		return dto.Gym{}, fiber.NewError(fiber.StatusBadRequest, "only manual gyms can be updated")
	}

	gym := gymSave.ToModel()
	gym.UserID = userID
	gym.ExternalID = oldGym.ExternalID
	gym.Source = oldGym.Source

	if err := g.gym.Update(ctx, gym); err != nil {
		return dto.Gym{}, err
	}

	return dto.GymDTO(&gym), nil
}
