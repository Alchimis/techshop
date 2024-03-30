package product

import (
	"context"

	"github.com/Alchimis/techshop/internal/models"
)

type Repository interface {
	GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error)
	GetProductsByOrderId(ctx context.Context, orderId int) ([]models.SimpleProduct, error)
}

type Service interface {
	GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error)
	GetProductsByOrderId(ctx context.Context, orderId int) ([]models.SimpleProduct, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repository
}

func (s *service) GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error) {
	return s.repo.GetProductsByIds(ctx, ids)
}

func (s *service) GetProductsByOrderId(ctx context.Context, orderId int) ([]models.SimpleProduct, error) {
	return s.repo.GetProductsByOrderId(ctx, orderId)
}
