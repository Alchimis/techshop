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
	SELECT rack_id, is_main FROM rack_has_product
	WHERE product_id=$1
	`
	rows, err := r.conn.Query(ctx, query, productId)
	defer rows.Close()
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
	SELECT rack_id, is_main FROM rack_has_product
	WHERE product_id=$1
	`
	rows, err := r.conn.Query(ctx, query, productId)
	defer rows.Close()
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

func (r rackRepository) GetRacksIdsByProductsIds(ctx context.Context, productsIds []int) ([]int, error) {
	query := `
	SELECT rack_id FROM rack_has_product
	WHERE product_id = any ($1);
	`
	rows, err := r.conn.Query(ctx, query, productsIds)
	defer rows.Close()
	if err != nil {
		return []int{}, err
	}
	var ids []int
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return []int{}, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r rackRepository) GetMainRacksIdsByProductsIds(ctx context.Context, productsIds []int) ([]int, error) {
	query := `
	SELECT rack_id FROM rack_has_product
	WHERE is_main=true AND product_id= any ($1);
	`
	rows, err := r.conn.Query(ctx, query, productsIds)
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

func (r rackRepository) GetMainRacksCountByProductsIds(ctx context.Context, productsIds []int) (int, error) {
	query := `
	SELECT COUNT(*) FROM rack_has_product
	WHERE product_id = any ($1) AND is_main=true;
	`
	var id int
	err := r.conn.QueryRow(ctx, query, productsIds).Scan(&id)
	return id, err
}

func (r rackRepository) GetRacksHasProductsByProductsIdsSplitByIsMain(ctx context.Context, productsIds []int) ([]models.RackIdProductsIds, []models.RackIdProductsIds, error) {
	query := `
	SELECT rack_id, product_id, is_main FROM rack_has_product
	WHERE product_id = any($1)
	ORDER BY rack_id;
	`
	rows, err := r.conn.Query(ctx, query, productsIds)
	defer rows.Close()
	if err != nil {
		return []models.RackIdProductsIds{}, []models.RackIdProductsIds{}, err
	}
	var (
		mainRacks       []models.RackIdProductsIds
		additionalRacks []models.RackIdProductsIds
		actualId        int   = -1
		additionalIds   []int = []int{}
		mainIds         []int = []int{}
	)
	for rows.Next() {
		var rack models.ProductIdRackId
		var isMain bool
		if err := rows.Scan(&rack.RackId, &rack.ProductId, &isMain); err != nil {
			return []models.RackIdProductsIds{}, []models.RackIdProductsIds{}, err
		}
		if actualId != rack.RackId {
			if len(mainIds) != 0 {
				mainRacks = append(mainRacks, models.RackIdProductsIds{
					RackId:      actualId,
					ProductsIds: mainIds,
				})
				mainIds = []int{}
			}
			if len(additionalIds) != 0 {
				additionalRacks = append(additionalRacks, models.RackIdProductsIds{
					RackId:      actualId,
					ProductsIds: additionalIds,
				})
				additionalIds = []int{}
			}
			actualId = rack.RackId
		}
		if isMain {
			mainIds = append(mainIds, rack.ProductId)
		} else {
			additionalIds = append(additionalIds, rack.ProductId)
		}
	}
	if len(mainIds) != 0 {
		mainRacks = append(mainRacks, models.RackIdProductsIds{
			RackId:      actualId,
			ProductsIds: mainIds,
		})
	}
	if len(additionalIds) != 0 {
		additionalRacks = append(additionalRacks, models.RackIdProductsIds{
			RackId:      actualId,
			ProductsIds: additionalIds,
		})
	}
	return mainRacks, additionalRacks, nil
}
