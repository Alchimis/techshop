package sql

import (
	"context"

	"github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) order.Repository {
	return orderRepository{conn: conn}
}

func (r orderRepository) GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error) {
	return []models.Order{}, errors.ErrNotImplemented
}

func (r orderRepository) GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error) {
	return []models.RackWithProducts{}, errors.ErrNotImplemented
}
