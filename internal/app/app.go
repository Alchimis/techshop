package app

import (
	"fmt"

	"github.com/Alchimis/techshop/internal/config"
	"github.com/Alchimis/techshop/internal/services/order"
)

func New() order.Service {
	config := config.NewConfig()
	connString := fmt.Sprintf("host=%s port=5432 ")
}
