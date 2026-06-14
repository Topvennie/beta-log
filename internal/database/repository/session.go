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
	session, err := s.repo.queries(ctx).SessionGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get session with id %d | %w", id, err)
	}

	return model.SessionModel(session), nil
}

func (s *Session) GetAllByUserID(ctx context.Context, userID int) ([]*model.Session, error) {
	sessions, err := s.repo.queries(ctx).SessionGetAll(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("get all sessions for user %d | %w", userID, err)
	}

	return utils.SliceMap(sessions, model.SessionModel), nil
}

func (s *Session) Create(ctx context.Context, session *model.Session) error {
	id, err := s.repo.queries(ctx).SessionCreate(ctx, sqlc.SessionCreateParams{
		UserID:   int32(session.UserID),
		Name:     session.Name,
		Active:   session.Active,
		Position: toPgInt4(session.Position),
	})
	if err != nil {
		return fmt.Errorf("create session %+v | %w", *session, err)
	}

	session.ID = int(id)

	return nil
}

func (s *Session) Update(ctx context.Context, session model.Session) error {
	err := s.repo.queries(ctx).SessionUpdate(ctx, sqlc.SessionUpdateParams{
		ID:       int32(session.ID),
		Name:     session.Name,
		Active:   session.Active,
		Position: toPgInt4(session.Position),
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
