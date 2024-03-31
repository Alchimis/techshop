package sql

import (
	"context"
	"errors"

	domainErrors "github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/product"
	"github.com/jackc/pgx/v5"
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

func (r productRepository) GetProductsByOrderId(ctx context.Context, orderId int) ([]models.OrderHasProduct, error) {
	query := `
	SELECT product_id, quantity FROM order_has_product
	WHERE order_id=$1;
	`
	rows, err := r.conn.Query(ctx, query, orderId)
	defer rows.Close()
	if err != nil {
		return []models.OrderHasProduct{}, err
	}
	var products []models.OrderHasProduct
	for rows.Next() {
		var product models.OrderHasProduct
		if err := rows.Scan(&product.ProductId, &product.ProductQuantity); err != nil {
			return []models.OrderHasProduct{}, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r productRepository) GetProductById(ctx context.Context, id int) (models.SimpleProduct, error) {
	query := `
	SELECT id, title FROM product
	WHERE id=$1; 
	`
	var product models.SimpleProduct
	err := r.conn.QueryRow(ctx, query, id).Scan(&product.Id, &product.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.SimpleProduct{}, domainErrors.ErrNotFound
		}
		return models.SimpleProduct{}, err
	}
	return product, nil
}

func (r productRepository) GetProductsIdsByOrderId(ctx context.Context, orderId int) ([]int, error) {
	query := `
	SELECT product_id FROM order_has_product
	WHERE order_id=$1;
	`
	rows, err := r.conn.Query(ctx, query, orderId)
	defer rows.Close()
	if err != nil {
		return []int{}, err
	}
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return []int{}, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
