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

func main() {
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
	orders, err := service.GetOrdersById(context.Background(), orderIds)
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(orders)
}
