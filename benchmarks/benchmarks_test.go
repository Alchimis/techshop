package benchmarks

import (
	"context"
	"testing"

	"github.com/Alchimis/techshop/internal/app"
)

var times = 10000

func BenchmarkAllData(b *testing.B) {
	orderIds := []int{10, 11, 14, 15}
	service, err := app.New()
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < times; i++ {
		_, err = service.GetOrdersByIdSortByRacks(context.Background(), orderIds)

		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkPartialDate(b *testing.B) {
	orderIds := []int{10, 11, 14, 15}
	service, err := app.New()
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < times; i++ {
		_, err = service.GetOrdersByIdsSortedByMainRacks(context.Background(), orderIds)

		if err != nil {
			b.Error(err)
			return
		}
	}
}
