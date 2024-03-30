package rack

import (
	"context"

	"github.com/Alchimis/techshop/internal/models"
)

type Repository interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
}

type Service interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error) {
	return s.repo.GetRacksByIds(ctx, ids)
}
