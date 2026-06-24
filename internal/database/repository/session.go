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

type Session struct {
	repo Repository
}

func (r *Repository) NewSession() *Session {
	return &Session{
		repo: *r,
	}
}

func (s *Session) Get(ctx context.Context, id int) (*model.Session, error) {
	rows, err := s.repo.queries(ctx).SessionGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get session with id %d | %w", id, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	session := model.SessionModel(rows[0].Session)

	for _, row := range rows {
		if !row.SessionExercisesView.ID.Valid {
			continue
		}

		sessionExercise := model.SessionExerciseViewModel(row.SessionExercisesView)
		sessionExercise.Exercise = *model.ExerciseViewModel(row.ExercisesView)
		if row.VariantsView.ID.Valid {
			sessionExercise.Variant = *model.VariantViewModel(row.VariantsView)
		}

		session.Exercises = append(session.Exercises, *sessionExercise)
	}

	return session, nil
}

func (s *Session) GetAllByUser(ctx context.Context, userID int) ([]*model.Session, error) {
	rows, err := s.repo.queries(ctx).SessionGetAll(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all sessions for user %d | %w", userID, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	sessionMap := map[int]*model.Session{}

	for _, row := range rows {
		session, ok := sessionMap[int(row.Session.ID)]
		if !ok {
			session = model.SessionModel(row.Session)
		}

		if row.SessionExercisesView.ID.Valid {
			sessionExercise := model.SessionExerciseViewModel(row.SessionExercisesView)
			sessionExercise.Exercise = *model.ExerciseViewModel(row.ExercisesView)
			if row.VariantsView.ID.Valid {
				sessionExercise.Variant = *model.VariantViewModel(row.VariantsView)
			}

			session.Exercises = append(session.Exercises, *sessionExercise)
		}

		sessionMap[session.ID] = session
	}

	return utils.MapValues(sessionMap), nil
}

func (s *Session) GetByExercise(ctx context.Context, exerciseID int) (*model.Session, error) {
	rows, err := s.repo.queries(ctx).SessionGetByExercise(ctx, toInt(exerciseID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get session by exercise %d | %w", exerciseID, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	session := model.SessionModel(rows[0].Session)

	for _, row := range rows {
		if !row.SessionExercisesView.ID.Valid {
			continue
		}

		sessionExercise := model.SessionExerciseViewModel(row.SessionExercisesView)
		sessionExercise.Exercise = *model.ExerciseViewModel(row.ExercisesView)
		if row.VariantsView.ID.Valid {
			sessionExercise.Variant = *model.VariantViewModel(row.VariantsView)
		}

		session.Exercises = append(session.Exercises, *sessionExercise)
	}

	return session, nil
}

func (s *Session) GetByVariants(ctx context.Context, variantIDs []int) ([]*model.Session, error) {
	rows, err := s.repo.queries(ctx).SessionGetByVariants(ctx, utils.SliceMap(variantIDs, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get session by variants %v | %w", variantIDs, err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	sessionMap := map[int]*model.Session{}

	for _, row := range rows {
		session, ok := sessionMap[int(row.Session.ID)]
		if !ok {
			session = model.SessionModel(row.Session)
		}

		if row.SessionExercisesView.ID.Valid {
			sessionExercise := model.SessionExerciseViewModel(row.SessionExercisesView)
			sessionExercise.Exercise = *model.ExerciseViewModel(row.ExercisesView)
			if row.VariantsView.ID.Valid {
				sessionExercise.Variant = *model.VariantViewModel(row.VariantsView)
			}

			session.Exercises = append(session.Exercises, *sessionExercise)
		}

		sessionMap[session.ID] = session
	}

	return utils.MapValues(sessionMap), nil
}

func (s *Session) Create(ctx context.Context, session *model.Session) error {
	id, err := s.repo.queries(ctx).SessionCreate(ctx, sqlc.SessionCreateParams{
		UserID: int32(session.UserID),
		Name:   session.Name,
	})
	if err != nil {
		return fmt.Errorf("create session %+v | %w", *session, err)
	}

	session.ID = int(id)

	return nil
}

func (s *Session) Update(ctx context.Context, session model.Session) error {
	err := s.repo.queries(ctx).SessionUpdate(ctx, sqlc.SessionUpdateParams{
		ID:   int32(session.ID),
		Name: session.Name,
	})
	if err != nil {
		return fmt.Errorf("update session %+v | %w", session, err)
	}

	return nil
}

func (s *Session) Delete(ctx context.Context, id int) error {
	err := s.repo.queries(ctx).SessionDelete(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("delete session with id %d | %w", id, err)
	}

	return nil
}
