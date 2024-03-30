package rack

import (
	"context"

	domainErrors "github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
)

type Repository interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetRacksHasProductsByProductIds(ctx context.Context, productIds []int) ([]models.RackHasProduct, error)
}

type Service interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetMainRacksByProductIds(ctx context.Context, ids []int) ([]models.MainRack, error)
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

func (s *service) GetMainRacksByProductIds(ctx context.Context, ids []int) ([]models.MainRack, error) {
	return []models.MainRack{}, domainErrors.ErrNotImplemented
}
