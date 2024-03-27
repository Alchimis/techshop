package sql

import (
	"context"
	"encoding/json"

	errs "errors"

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
	query := `
	SELECT 
	p.main_rack_id,
	p.main_rack_name,
	json_agg(
		json_build_object(
			'order_id', order_has_product.order_id,
			'product_id',p.id, 
			'product_title',p.title, 
			'order_quantity',order_has_product.quantity,   
			'additional_racks',p.racks)) products
FROM client_order
LEFT JOIN order_has_product 
	ON client_order.id=order_has_product.order_id
LEFT JOIN (
	SELECT pr.id, 
		pr.title, 
		pr.main_rack_id,
		pr.main_rack_name,
		array_agg(json_build_object('rack_id', rack.id, 'rack_name', rack.name)) racks 
	FROM 
		(SELECT product.id, product.title, product.main_rack_id, rack.name as main_rack_name  FROM product
		 LEFT JOIN rack ON product.main_rack_id=rack.id
		) pr
	LEFT JOIN rack_has_product 
		ON rack_has_product.product_id=pr.id 
	LEFT JOIN rack 
		ON rack.id=rack_has_product.rack_id
	GROUP BY pr.id, pr.title, pr.main_rack_id, pr.main_rack_name) p
	ON order_has_product.product_id=p.id
WHERE client_order.id = any ($1)
GROUP BY p.main_rack_id, p.main_rack_name;
	`
	rows, err := r.conn.Query(ctx, query, ids)
	defer rows.Close()
	if err != nil {
		return []models.RackWithProducts{}, errs.Join(errs.New("database sql: error with query execution: "), err)
	}
	var racks []models.RackWithProducts
	for rows.Next() {
		var (
			rack     models.RackWithProducts
			jsonBody *string
		)
		rows.Scan(&rack.Id, &rack.Name, &jsonBody)
		if jsonBody == nil {
			return []models.RackWithProducts{}, errs.New("json body was empty")
		}

		err := json.Unmarshal([]byte(*jsonBody), &rack.Products)
		if err != nil {
			return racks, err
		}
		racks = append(racks, rack)
	}
	return racks, nil
}
