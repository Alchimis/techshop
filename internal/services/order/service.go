package order

import (
	"context"
	"fmt"

	"github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/product"
)

type Repository interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
	GetOrderHasProductByIds(ctx context.Context, ids []int) ([]models.OrderHasProduct, error)
}

type service struct {
	repo           Repository
	productService product.Service
}

type Service interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
}

func NewService(repo Repository, productService product.Service) Service {
	return &service{
		repo:           repo,
		productService: productService,
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
	productIds := []int{}
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
			productIds = append(productIds, o.ProductId)
		} else {
			order := append(order, struct {
				Id       int
				Quantity int
			}{Id: o.OrderId, Quantity: o.ProductQuantity})
			productsIdsWithOrders[o.ProductId] = order
		}
	}
	products, err := s.productService.GetProductsByIds(ctx, productIds)
	if err != nil {
		return []models.RackWithProducts{}, err
	}
	fmt.Println("products with orders", productsIdsWithOrders)
	fmt.Println("products", products)
	return []models.RackWithProducts{}, nil //s.repo.GetOrdersByIdSortByRacks(ctx, ids)
}
