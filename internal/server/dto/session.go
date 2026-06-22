package dto

import (
	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/utils"
)

// Session

type Session struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Active    bool              `json:"active"`
	Position  int               `json:"position,omitzero"`
	Exercises []SessionExercise `json:"exercises"`
}

func SessionDTO(s *model.Session) Session {
	return Session{
		ID:        s.ID,
		Name:      s.Name,
		Active:    s.Active,
		Position:  s.Position,
		Exercises: []SessionExercise{},
	}
}

func SessionDTOPopulated(session *model.Session, sessionExercises []*model.SessionExercise, exercises []*model.Exercise) Session {
	exerciseMap := make(map[int]*model.Exercise)
	for _, e := range exercises {
		exerciseMap[e.ID] = e
	}

	s := SessionDTO(session)

	for _, sessionExercise := range sessionExercises {
		e, ok := exerciseMap[sessionExercise.ExerciseID]
		if !ok {
			continue
		}

		s.Exercises = append(s.Exercises, SessionExerciseDTO(sessionExercise, e))
	}

	return s
}

type SessionCreate struct {
	Name      string                  `json:"name" validate:"required"`
	Active    *bool                   `json:"active" validate:"required"`
	Position  int                     `json:"position"`
	Exercises []SessionExerciseCreate `json:"exercises" validate:"required,min=1"`
}

func (s SessionCreate) ToModel() model.Session {
	return model.Session{
		Name:      s.Name,
		Active:    s.Active != nil && *s.Active,
		Position:  s.Position,
		Exercises: utils.SliceMap(s.Exercises, func(s SessionExerciseCreate) model.SessionExercise { return s.ToModel() }),
	}
}

type SessionUpdate struct {
	ID        int                     `json:"id" validate:"required,min=1"`
	Name      string                  `json:"name" validate:"required"`
	Active    *bool                   `json:"active" validate:"required"`
	Position  int                     `json:"position"`
	Exercises []SessionExerciseUpdate `json:"exercises" validate:"required,min=1"`
}

func (s SessionUpdate) ToModel() model.Session {
	return model.Session{
		ID:        s.ID,
		Name:      s.Name,
		Active:    s.Active != nil && *s.Active,
		Position:  s.Position,
		Exercises: utils.SliceMap(s.Exercises, func(s SessionExerciseUpdate) model.SessionExercise { return s.ToModel() }),
	}
}
