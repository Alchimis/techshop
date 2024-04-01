package product

import (
	"context"

	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/utils"
)

type Repository interface {
	GetProductById(ctx context.Context, id int) (models.SimpleProduct, error)
	GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error)
	GetProductsByOrderId(ctx context.Context, orderId int) ([]models.OrderHasProduct, error)
	GetProductsIdsByOrderId(ctx context.Context, orderId int) ([]int, error)
	GetProductIdAndQuantityByOrderId(ctx context.Context, orderId int) ([]models.ProductIdAndQuantity, error)
}

type Service interface {
	GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error)
	GetProductsByOrderId(ctx context.Context, orderId int) ([]models.OrderHasProduct, error)
	GetProductsOrdersByOrderId(ctx context.Context, orderId int) ([]models.ProductOrder, error)
	GetProductIdAndQuantityByOrderId(ctx context.Context, orderId int) ([]models.ProductIdAndQuantity, error)
	GetProductById(ctx context.Context, id int) (models.SimpleProduct, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repository
}

func (s *service) GetProductIdAndQuantityByOrderId(ctx context.Context, orderId int) ([]models.ProductIdAndQuantity, error) {
	return s.repo.GetProductIdAndQuantityByOrderId(ctx, orderId)
}

func (s *service) GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error) {
	return s.repo.GetProductsByIds(ctx, ids)
}

func (s *service) GetProductById(ctx context.Context, id int) (models.SimpleProduct, error) {
	return s.repo.GetProductById(ctx, id)
}

func (s *service) GetProductsByOrderId(ctx context.Context, orderId int) ([]models.OrderHasProduct, error) {
	return s.repo.GetProductsByOrderId(ctx, orderId)
}

func (s *service) GetProductsOrdersByOrderId(ctx context.Context, orderId int) ([]models.ProductOrder, error) {
	pm, err := utils.GetProducts(ctx)
	if err != nil {
		return []models.ProductOrder{}, err
	}
	orders, err := s.repo.GetProductsByOrderId(ctx, orderId)
	if err != nil {
		return []models.ProductOrder{}, err
	}
	var productOrders []models.ProductOrder
	for _, order := range orders {
		product, ok := pm[orderId]
		if !ok {
			product, err = s.repo.GetProductById(ctx, order.ProductId)
			if err != nil {
				return []models.ProductOrder{}, err
			}
			pm[orderId] = product
		}
		productOrders = append(productOrders, models.ProductOrder{
			Id:       order.ProductId,
			OrderId:  order.OrderId,
			Quantity: order.ProductQuantity,
			Title:    product.Title,
		})
	}
	return productOrders, nil
}
