package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type Gym struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	IconPath string       `json:"icon_path"`
	Source   model.Source `json:"source"`
}

func GymDTO(g *model.Gym) Gym {
	return Gym{
		ID:       g.ID,
		Name:     g.Name,
		IconPath: g.IconPath,
		Source:   g.Source,
	}
}

type GymCreate struct {
	Name     string `json:"name" validate:"required"`
	IconPath string `json:"icon_path"`
}

func (g GymCreate) ToModel() model.Gym {
	return model.Gym{
		Name:     g.Name,
		IconPath: g.IconPath,
	}
}

type GymUpdate struct {
	ID       int    `json:"id" validate:"required,min=1"`
	Name     string `json:"name" validate:"required"`
	IconPath string `json:"icon_path"`
}

func (g GymUpdate) ToModel() model.Gym {
	return model.Gym{
		ID:       g.ID,
		Name:     g.Name,
		IconPath: g.IconPath,
	}
}
