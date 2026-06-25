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

type User struct {
	repo Repository
}

func (r *Repository) NewUser() *User {
	return &User{
		repo: *r,
	}
}

func (u *User) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with id %d | %w", id, err)
	}

	return model.UserModel(user), nil
}

func (u *User) GetByUID(ctx context.Context, uid string) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with uid %s | %w", uid, err)
	}

	return model.UserModel(user), nil
}

func (u *User) GetAll(ctx context.Context) ([]*model.User, error) {
	users, err := u.repo.queries(ctx).UserGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all users %w", err)
	}

	return utils.SliceMap(users, model.UserModel), nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	id, err := u.repo.queries(ctx).UserCreate(ctx, sqlc.UserCreateParams{
		Uid:  user.UID,
		Name: user.Name,
	})
	if err != nil {
		return fmt.Errorf("create user %+v | %w", *user, err)
	}

	user.ID = int(id)

	return nil
}
