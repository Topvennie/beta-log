package service

import (
	"context"

	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/gofiber/fiber/v3"
)

type User struct {
	service Service

	user repository.User
}

func (s *Service) NewUser() *User {
	return &User{
		service: *s,
		user:    *s.repo.NewUser(),
	}
}

func (u *User) GetByID(ctx context.Context, id int) (dto.User, error) {
	user, err := u.user.GetByID(ctx, id)
	if err != nil {
		return dto.User{}, err
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(*user), nil
}

func (u *User) GetByUID(ctx context.Context, uid string) (dto.User, error) {
	user, err := u.user.GetByUID(ctx, uid)
	if err != nil {
		return dto.User{}, err
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(*user), nil
}

func (u *User) Save(ctx context.Context, userSave dto.User) (dto.User, error) {
	user := userSave.ToModel()
	if err := u.user.Create(ctx, &user); err != nil {
		return dto.User{}, err
	}

	return dto.UserDTO(user), nil
}
