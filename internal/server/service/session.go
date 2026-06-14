package service

import (
	"context"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type Session struct {
	service         Service
	session         repository.Session
	sessionExercise repository.SessionExercise
	exercise        repository.Exercise
}

func (s *Service) NewSession() *Session {
	return &Session{
		service:         *s,
		session:         *s.repo.NewSession(),
		sessionExercise: *s.repo.NewSessionExercise(),
		exercise:        *s.repo.NewExercise(),
	}
}

func (s *Session) GetAll(ctx fiber.Ctx, userID int) ([]dto.Session, error) {
	sessions, err := s.session.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.Session, 0, len(sessions))
	for _, session := range sessions {
		se, err := s.sessionExercise.GetBySession(ctx, session.ID)
		if err != nil {
			return nil, err
		}

		exercises, err := s.populateExercises(ctx, se)
		if err != nil {
			return nil, err
		}

		result = append(result, dto.SessionDTOPopulated(session, se, exercises))
	}

	return result, nil
}

func (s *Session) Create(ctx fiber.Ctx, sessionCreate dto.SessionCreate) (dto.Session, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Session{}, err
	}

	sessionModel := sessionCreate.ToModel()
	sessionModel.UserID = userID

	if err := s.service.withRollback(ctx, func(ctx context.Context) error {
		if err := s.session.Create(ctx, &sessionModel); err != nil {
			return err
		}

		for i := range sessionModel.Exercises {
			sessionModel.Exercises[i].SessionID = sessionModel.ID
			if err := s.sessionExercise.Create(ctx, &sessionModel.Exercises[i]); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return dto.Session{}, err
	}

	// Refetch data

	sessionExercises, err := s.sessionExercise.GetBySession(ctx, sessionModel.ID)
	if err != nil {
		return dto.Session{}, err
	}

	exercises, err := s.populateExercises(ctx, sessionExercises)
	if err != nil {
		return dto.Session{}, err
	}

	return dto.SessionDTOPopulated(&sessionModel, sessionExercises, exercises), nil
}

// nolint:gocognit // It's fine
func (s *Session) Update(ctx fiber.Ctx, sessionUpdate dto.SessionUpdate) (dto.Session, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Session{}, err
	}

	sessionOld, err := s.session.Get(ctx, sessionUpdate.ID)
	if err != nil {
		return dto.Session{}, err
	}
	if sessionOld == nil {
		return dto.Session{}, fiber.ErrNotFound
	}
	if sessionOld.UserID != userID {
		return dto.Session{}, fiber.ErrNotFound
	}

	sessionModel := sessionUpdate.ToModel()
	sessionModel.UserID = userID

	if err := s.service.withRollback(ctx, func(ctx context.Context) error {
		if err := s.session.Update(ctx, sessionModel); err != nil {
			return err
		}

		existingSE, err := s.sessionExercise.GetBySession(ctx, sessionModel.ID)
		if err != nil {
			return err
		}

		incomingMap := make(map[int]bool)
		for _, se := range sessionModel.Exercises {
			if se.ID != 0 {
				incomingMap[se.ID] = true
			}
		}

		for _, se := range existingSE {
			if ok := incomingMap[se.ID]; !ok {
				if err := s.sessionExercise.Delete(ctx, se.ID); err != nil {
					return err
				}
			}
		}

		for i := range sessionModel.Exercises {
			sessionModel.Exercises[i].SessionID = sessionModel.ID

			if sessionModel.Exercises[i].ID != 0 {
				if err := s.sessionExercise.Update(ctx, sessionModel.Exercises[i]); err != nil {
					return err
				}
			} else {
				if err := s.sessionExercise.Create(ctx, &sessionModel.Exercises[i]); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return dto.Session{}, err
	}

	// Refetch data

	sessionExercises, err := s.sessionExercise.GetBySession(ctx, sessionModel.ID)
	if err != nil {
		return dto.Session{}, err
	}

	exercises, err := s.populateExercises(ctx, sessionExercises)
	if err != nil {
		return dto.Session{}, err
	}

	return dto.SessionDTOPopulated(&sessionModel, sessionExercises, exercises), nil
}

func (s *Session) Delete(ctx fiber.Ctx, id int) error {
	userID, err := getID(ctx)
	if err != nil {
		return err
	}

	session, err := s.session.Get(ctx, id)
	if err != nil {
		return err
	}
	if session == nil {
		return fiber.ErrNotFound
	}
	if session.UserID != userID {
		return fiber.ErrNotFound
	}

	return s.service.withRollback(ctx, func(ctx context.Context) error {
		exercises, err := s.sessionExercise.GetBySession(ctx, id)
		if err != nil {
			return err
		}

		for _, se := range exercises {
			if err := s.sessionExercise.Delete(ctx, se.ID); err != nil {
				return err
			}
		}

		return s.session.Delete(ctx, id)
	})
}

func (s *Session) populateExercises(ctx context.Context, sessionExercises []*model.SessionExercise) ([]*model.Exercise, error) {
	if len(sessionExercises) == 0 {
		return []*model.Exercise{}, nil
	}

	return s.exercise.GetByIDs(ctx, utils.SliceMap(sessionExercises, func(se *model.SessionExercise) int { return se.ExerciseID }))
}
