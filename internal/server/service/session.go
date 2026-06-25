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

type Session struct {
	service Service

	exercise        repository.Exercise
	session         repository.Session
	sessionExercise repository.SessionExercise
	variant         repository.Variant
}

func (s *Service) NewSession() *Session {
	return &Session{
		service:         *s,
		exercise:        *s.repo.NewExercise(),
		session:         *s.repo.NewSession(),
		sessionExercise: *s.repo.NewSessionExercise(),
		variant:         *s.repo.NewVariant(),
	}
}

func (s *Session) GetAll(ctx fiber.Ctx, userID int) ([]dto.Session, error) {
	sessions, err := s.session.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return utils.SliceMap(sessions, dto.SessionDTO), nil
}

func (s *Session) Create(ctx fiber.Ctx, sessionCreate dto.SessionCreate) (dto.Session, error) {
	userID, err := getID(ctx)
	if err != nil {
		return dto.Session{}, err
	}

	session := sessionCreate.ToModel()
	session.UserID = userID

	// Check if all used exercises are owned by the user
	exercises, err := s.exercise.GetByIDs(ctx, utils.SliceMap(session.Exercises, func(s model.SessionExercise) int { return s.ExerciseID }))
	if err != nil {
		return dto.Session{}, err
	}
	if idx := slices.IndexFunc(exercises, func(e *model.Exercise) bool { return e.UserID != userID }); idx != -1 {
		return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "unknown exercise")
	}

	// Check if all variants do exist for that exercise
	for _, sessionExercise := range session.Exercises {
		// Find the actual exercise
		exercise, ok := utils.SliceFind(exercises, func(e *model.Exercise) bool { return e.ID == sessionExercise.ExerciseID })
		if !ok {
			return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "unknown exercise")
		}

		// Check the variant
		if sessionExercise.VariantID == 0 {
			continue
		}
		if _, ok = utils.SliceFind(exercise.Variants, func(v model.Variant) bool { return v.ID == sessionExercise.VariantID }); !ok {
			return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "unknown variant")
		}
	}

	// Verify exercises position
	for i := range len(session.Exercises) {
		if _, ok := utils.SliceFind(session.Exercises, func(e model.SessionExercise) bool { return e.Position == i+1 }); !ok {
			return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "positions need to increment by one starting with one")
		}
	}

	if err := s.service.withRollback(ctx, func(ctx context.Context) error {
		if err := s.session.Create(ctx, &session); err != nil {
			return err
		}

		for i := range session.Exercises {
			session.Exercises[i].SessionID = session.ID
			if err := s.sessionExercise.Create(ctx, &session.Exercises[i]); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return dto.Session{}, err
	}

	// Refetch data

	newSession, err := s.session.Get(ctx, session.ID)
	if err != nil {
		return dto.Session{}, err
	}

	return dto.SessionDTO(newSession), nil
}

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

	session := sessionUpdate.ToModel()
	session.UserID = userID

	// Check if all used exercises are owned by the user
	exercises, err := s.exercise.GetByIDs(ctx, utils.SliceMap(session.Exercises, func(s model.SessionExercise) int { return s.ExerciseID }))
	if err != nil {
		return dto.Session{}, err
	}
	if idx := slices.IndexFunc(exercises, func(e *model.Exercise) bool { return e.UserID != userID }); idx != -1 {
		return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "unknown exercise")
	}

	// Check if all variants do exist for that exercise
	for _, sessionExercise := range session.Exercises {
		// Find the actual exercise
		exercise, ok := utils.SliceFind(exercises, func(e *model.Exercise) bool { return e.ID == sessionExercise.ExerciseID })
		if !ok {
			return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "unknown exercise")
		}

		// Check the variant
		if sessionExercise.VariantID == 0 {
			continue
		}
		if _, ok = utils.SliceFind(exercise.Variants, func(v model.Variant) bool { return v.ID == sessionExercise.VariantID }); !ok {
			return dto.Session{}, fiber.NewError(fiber.StatusBadRequest, "unknown variant")
		}
	}

	if err := s.service.withRollback(ctx, func(ctx context.Context) error {
		if err := s.session.Update(ctx, session); err != nil {
			return err
		}

		if err := s.sessionExercise.DeleteBySession(ctx, session.ID); err != nil {
			return err
		}

		for i := range session.Exercises {
			session.Exercises[i].SessionID = session.ID
			if err := s.sessionExercise.Create(ctx, &session.Exercises[i]); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return dto.Session{}, err
	}

	// Refetch data

	newSession, err := s.session.Get(ctx, session.ID)
	if err != nil {
		return dto.Session{}, err
	}

	return dto.SessionDTO(newSession), nil
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
