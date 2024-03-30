package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Alchimis/techshop/internal/app"
	"github.com/Alchimis/techshop/internal/services/order"
)

func stringsToInts(strings []string) ([]int, error) {
	var ints []int
	for _, s := range strings[1:] {
		i, err := strconv.Atoi(s)
		if err != nil {
			return []int{}, errors.Join(errors.New("cmd .stringsToInt(): can't parse int: "), err)
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func handleError(err error) {
	fmt.Println(err)
}

func getOrders(service order.Service) {
	orderIds, err := stringsToInts(os.Args)
	if err != nil {
		handleError(err)
		return
	}
	orders, err := service.GetOrdersByIdSortByRacks(context.Background(), orderIds)
	if err != nil {
		handleError(err)
		return
	}
	for _, o := range orders {
		fmt.Println(o)
	}
}

func main() {
	service, err := app.New()
	if err != nil {
		handleError(err)
		return
	}
	startedAt := time.Now()
	for i := 0; i < 10000; i++ {
		getOrders(service)
	}
	fmt.Println(time.Since(startedAt))
}
