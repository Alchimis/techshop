package app

import (
	"context"
	"fmt"

	"github.com/Alchimis/techshop/internal/config"
	"github.com/Alchimis/techshop/internal/database/sql"
	"github.com/Alchimis/techshop/internal/services/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (order.Service, error) {
	config := config.NewConfig()
	connString := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName)
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	orderRepo := sql.NewRepository(conn)
	return order.NewService(orderRepo), nil
}
