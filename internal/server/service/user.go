// Package service is the business logic connects the api with the internal mechanisms
package service

import (
	"context"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/gofiber/fiber/v3"
)

type User struct {
	setting repository.Setting
	user    repository.User
}

func NewUser() *User {
	return &User{
		user: *repository.NewUser(),
	}
}

func (u *User) GetMe(ctx fiber.Ctx) (dto.User, error) {
	id, err := getID(ctx)
	if err != nil {
		return dto.User{}, err
	}

	user, err := u.user.GetByID(ctx, id)
	if err != nil {
		return dto.User{}, err
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(user), nil
}

func (u *User) GetByUID(ctx fiber.Ctx, uid string) (dto.User, error) {
	user, err := u.user.GetByUID(ctx, uid)
	if err != nil {
		return dto.User{}, err
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(user), nil
}

func (u *User) Create(ctx fiber.Ctx, userSave dto.User) (dto.User, error) {
	user := userSave.ToModel()

	if err := withRollback(ctx, func(ctx context.Context) error {
		if err := u.user.Create(ctx, &user); err != nil {
			return err
		}

		if err := u.setting.Create(ctx, &model.Setting{UserID: user.ID}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return dto.User{}, err
	}

	return dto.UserDTO(&user), nil
}
