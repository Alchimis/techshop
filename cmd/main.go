package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Alchimis/techshop/internal/app"
	"github.com/Alchimis/techshop/internal/utils"
)

func handleError(err error) {
	fmt.Println(err)
}

func main() {
	service, err := app.New()
	if err != nil {
		handleError(err)
		return
	}
	orderIds, err := utils.StringsToInts(os.Args)
	if err != nil {
		handleError(err)
		return
	}
	order, err := service.GetOrdersByIdSortByRacks(context.Background(), orderIds)

	if err != nil {
		handleError(err)
		return
	}
	for _, o := range order {
		fmt.Println(o)
	}
}
