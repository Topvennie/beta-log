package dto

import (
	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/utils"
)

// Session

type Session struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Exercises []SessionExercise `json:"exercises"`
}

func SessionDTO(s *model.Session) Session {
	return Session{
		ID:        s.ID,
		Name:      s.Name,
		Exercises: utils.SliceMap(s.Exercises, func(e model.SessionExercise) SessionExercise { return SessionExerciseDTO(&e) }),
	}
}

type SessionCreate struct {
	Name      string                  `json:"name" validate:"required"`
	Exercises []SessionExerciseCreate `json:"exercises" validate:"required,min=1"`
}

func (s SessionCreate) ToModel() model.Session {
	return model.Session{
		Name:      s.Name,
		Exercises: utils.SliceMap(s.Exercises, func(s SessionExerciseCreate) model.SessionExercise { return s.ToModel() }),
	}
}

type SessionUpdate struct {
	ID        int                     `json:"id" validate:"required,min=1"`
	Name      string                  `json:"name" validate:"required"`
	Exercises []SessionExerciseUpdate `json:"exercises" validate:"required,min=1"`
}

func (s SessionUpdate) ToModel() model.Session {
	return model.Session{
		ID:        s.ID,
		Name:      s.Name,
		Exercises: utils.SliceMap(s.Exercises, func(s SessionExerciseUpdate) model.SessionExercise { return s.ToModel() }),
	}
}
