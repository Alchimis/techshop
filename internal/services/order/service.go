package order

import (
	"context"

	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services"
)

type Repository interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
}

type service struct {
	repo Repository
}

type Service interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error) {
	return []models.Order{}, services.ErrNotImplemented
}
