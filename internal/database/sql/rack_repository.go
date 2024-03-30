package sql

import (
	"context"

	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/rack"
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
