package sql

import (
	"context"

	domainErrors "github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/rack"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NerRackRepository(conn *pgxpool.Pool) rack.Repository {
	return rackRepository{conn: conn}
}

type rackRepository struct {
	conn *pgxpool.Pool
}

func (r rackRepository) GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error) {
	query := `
	SELECT id, name FROM rack
	WHERE id = any ($1);
	`
	rows, err := r.conn.Query(ctx, query, ids)
	defer rows.Close()
	if err != nil {
		return []models.Rack{}, nil
	}
	var racks []models.Rack
	for rows.Next() {
		var rack models.Rack
		if err := rows.Scan(&rack.Id, &rack.Title); err != nil {
			return []models.Rack{}, err
		}
		racks = append(racks, rack)
	}
	return racks, nil
}

func (r rackRepository) GetRacksHasProductsByProductIds(ctx context.Context, productIds []int) ([]models.RackHasProduct, error) {
	query := `
	SELECT rack_id, product_id, is_main FROM rack_has_product
	WHERE product_id = any($1)
	`
	rows, err := r.conn.Query(ctx, query, productIds)
	defer rows.Close()
	if err != nil {
		return []models.RackHasProduct{}, err
	}
	var rackHasProducts []models.RackHasProduct
	for rows.Next() {
		var rack models.RackHasProduct
		if err := rows.Scan(&rack.RackId, &rack.ProductId, &rack.IsMain); err != nil {
			return []models.RackHasProduct{}, err
		}
		rackHasProducts = append(rackHasProducts, rack)
	}
	return rackHasProducts, nil
}

func (r rackRepository) GetRackById(ctx context.Context, id int) (models.Rack, error) {
	query := `
	SELECT id, name FROM rack
	WHERE id=$1;
	`
	var rack models.Rack
	err := r.conn.QueryRow(ctx, query, id).Scan(&rack.Id, &rack.Title)
	if err != nil {
		if err != pgx.ErrNoRows {
			return models.Rack{}, domainErrors.ErrNotFound
		}
		return models.Rack{}, err
	}
	return rack, nil
}

func (r rackRepository) GetRacksByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error) {
	query := `
	SELECT rack_id, if_main FROM rack_has_product
	WHERE product_id=$1
	`
	rows, err := r.conn.Query(ctx, query, productId)
	rows.Close()
	if err != nil {
		return []models.RackWithIsMain{}, err
	}
	var racks []models.RackWithIsMain
	for rows.Next() {
		var rack models.RackWithIsMain
		err := rows.Scan(&rack.Id, &rack.IsMain)
		if err != nil {
			return []models.RackWithIsMain{}, err
		}
		racks = append(racks, rack)
	}
	return racks, nil
}

func (r rackRepository) GetAllRackOfProduct(ctx context.Context, productId int) (models.RacksOfProduct, error) {
	query := `
	SELECT rack_id, if_main FROM rack_has_product
	WHERE product_id=$1
	`
	rows, err := r.conn.Query(ctx, query, productId)
	rows.Close()
	if err != nil {
		return models.RacksOfProduct{}, err
	}
	var racks models.RacksOfProduct
	for rows.Next() {
		var (
			id     int
			isMain bool
		)
		err := rows.Scan(&id, &isMain)
		if err != nil {
			return models.RacksOfProduct{}, err
		}
		if isMain {
			racks.MainRackId = id
		} else {
			racks.AdditionalRacks = append(racks.AdditionalRacks, id)
		}
	}
	return racks, nil
}
