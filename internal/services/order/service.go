package order

import (
	"context"
	"fmt"

	"github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
)

type Repository interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
	GetOrderHasProductByIds(ctx context.Context, ids []int) ([]models.OrderHasProduct, error)
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
	ordersHasProducts, err := s.repo.GetOrderHasProductByIds(ctx, ids)
	if err != nil {
		return []models.RackWithProducts{}, err
	}
	productsIdsWithOrders := map[int][]struct {
		Id       int
		Quantity int
	}{}
	for _, o := range ordersHasProducts {
		order, ok := productsIdsWithOrders[o.ProductId]
		if !ok {
			productsIdsWithOrders[o.ProductId] = []struct {
				Id       int
				Quantity int
			}{{
				Id:       o.OrderId,
				Quantity: o.ProductQuantity,
			}}
		} else {
			order := append(order, struct {
				Id       int
				Quantity int
			}{Id: o.OrderId, Quantity: o.ProductQuantity})
			productsIdsWithOrders[o.ProductId] = order
		}
	}
	fmt.Println(productsIdsWithOrders)
	return []models.RackWithProducts{}, nil //s.repo.GetOrdersByIdSortByRacks(ctx, ids)
}
