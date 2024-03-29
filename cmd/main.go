package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Alchimis/techshop/internal/app"
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

func getOrders() {
	orderIds, err := stringsToInts(os.Args)
	if err != nil {
		handleError(err)
		return
	}
	service, err := app.New()
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
	getOrders()
}
