package order

import (
	"context"

	"github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
)

type Repository interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
}

type service struct {
	repo Repository
}

type Service interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error) {
	return []models.Order{}, errors.ErrNotImplemented
}

func (s *service) GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error) {
	return s.repo.GetOrdersByIdSortByRacks(ctx, ids)
}
