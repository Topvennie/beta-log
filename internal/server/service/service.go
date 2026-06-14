// Package service is the business logic connects the api with the internal mechanisms
package service

import (
	"context"

	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/gofiber/fiber/v3"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) withRollback(ctx fiber.Ctx, fn func(context.Context) error) error {
	return s.repo.WithRollback(ctx, fn)
}
