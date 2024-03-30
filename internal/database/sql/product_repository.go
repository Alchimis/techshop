package sql

import (
	"context"

	domainErrors "github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewProductRepository(conn *pgxpool.Pool) product.Repository {
	return productRepository{conn: conn}
}

type productRepository struct {
	conn *pgxpool.Pool
}

func (r productRepository) GetProductsByIds(ctx context.Context, ids []int) ([]models.SimpleProduct, error) {
	query := `
	SELECT id, title FROM product
	WHERE id = any ($1)
	ORDER BY id;
	`
	rows, err := r.conn.Query(ctx, query, ids)
	defer rows.Close()
	if err != nil {
		return []models.SimpleProduct{}, err
	}
	products := []models.SimpleProduct{}
	for rows.Next() {
		var product models.SimpleProduct
		if err := rows.Scan(&product.Id, &product.Title); err != nil {
			return []models.SimpleProduct{}, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r productRepository) GetProductsByOrderId(ctx context.Context, orderId int) ([]models.SimpleProduct, error) {
	return []models.SimpleProduct{}, domainErrors.ErrNotImplemented
}
