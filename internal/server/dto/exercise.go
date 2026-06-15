package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type Exercise struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Variants []string `json:"variants"`
}

func ExerciseDTO(e *model.Exercise) Exercise {
	return Exercise{
		ID:       e.ID,
		Name:     e.Name,
		Variants: e.Variants,
	}
}

type ExerciseCreate struct {
	Name     string   `json:"name" validate:"required"`
	Variants []string `json:"variants"`
}

func (e ExerciseCreate) ToModel() model.Exercise {
	return model.Exercise{
		Name:     e.Name,
		Variants: e.Variants,
	}
}

type ExerciseUpdate struct {
	ID       int      `json:"id" validate:"required,min=1"`
	Name     string   `json:"name" validate:"required"`
	Variants []string `json:"variants"`
}

func (e ExerciseUpdate) ToModel() model.Exercise {
	return model.Exercise{
		ID:       e.ID,
		Name:     e.Name,
		Variants: e.Variants,
	}
}
