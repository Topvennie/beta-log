package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/sqlc"
	"github.com/Topvennie/beta-log/pkg/utils"
)

type SessionExercise struct {
	repo Repository
}

func (r *Repository) NewSessionExercise() *SessionExercise {
	return &SessionExercise{
		repo: *r,
	}
}

func (s *SessionExercise) Get(ctx context.Context, id int) (*model.SessionExercise, error) {
	sessionExercise, err := s.repo.queries(ctx).SessionExerciseGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get session exercise with id %d | %w", id, err)
	}

	return model.SessionExerciseModel(sessionExercise), nil
}

func (s *SessionExercise) GetBySession(ctx context.Context, sessionID int) ([]*model.SessionExercise, error) {
	sessionExercises, err := s.repo.queries(ctx).SessionExerciseGetBySession(ctx, int32(sessionID))
	if err != nil {
		return nil, fmt.Errorf("get all session exercises for session %d | %w", sessionID, err)
	}

	return utils.SliceMap(sessionExercises, model.SessionExerciseModel), nil
}

func (s *SessionExercise) Create(ctx context.Context, sessionExercise *model.SessionExercise) error {
	id, err := s.repo.queries(ctx).SessionExerciseCreate(ctx, sqlc.SessionExerciseCreateParams{
		SessionID:  int32(sessionExercise.SessionID),
		ExerciseID: int32(sessionExercise.ExerciseID),
		VariantID:  toInt(sessionExercise.VariantID),
		Position:   int32(sessionExercise.Position),
		Sets:       int32(sessionExercise.Sets),
		Reps:       toInt(sessionExercise.Reps),
		Weight:     toInt(sessionExercise.Weight),
		DurationS:  toInt(sessionExercise.DurationS),
	})
	if err != nil {
		return fmt.Errorf("create session exercise %+v | %w", *sessionExercise, err)
	}

	sessionExercise.ID = int(id)

	return nil
}

func (s *SessionExercise) Update(ctx context.Context, sessionExercise model.SessionExercise) error {
	if err := s.repo.queries(ctx).SessionExerciseUpdate(ctx, sqlc.SessionExerciseUpdateParams{
		ID:        int32(sessionExercise.ID),
		VariantID: toInt(sessionExercise.VariantID),
		Position:  int32(sessionExercise.Position),
		Sets:      int32(sessionExercise.Sets),
		Reps:      toInt(sessionExercise.Reps),
		Weight:    toInt(sessionExercise.Weight),
		DurationS: toInt(sessionExercise.DurationS),
	}); err != nil {
		return fmt.Errorf("update session exercise %+v | %w", sessionExercise, err)
	}

	return nil
}

func (s *SessionExercise) Delete(ctx context.Context, id int) error {
	if err := s.repo.queries(ctx).SessionExerciseDelete(ctx, int32(id)); err != nil {
		return fmt.Errorf("delete session exercise with id %d | %w", id, err)
	}

	return nil
}

func (s *SessionExercise) DeleteBySession(ctx context.Context, sessionID int) error {
	if err := s.repo.queries(ctx).SessionExerciseDeleteBySession(ctx, int32(sessionID)); err != nil {
		return fmt.Errorf("delete session exercises by session id %d | %w", sessionID, err)
	}

	return nil
}
