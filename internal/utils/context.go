package utils

import (
	"context"

	"errors"

	"github.com/Alchimis/techshop/internal/models"
)

type ContextContent string

const (
	PRODUCTS ContextContent = "PRODUCTS"
	RACKS    ContextContent = "RACKS"
)

var (
	ErrKeyNotFound error = errors.New("key not found in context")
)

func SetupContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, PRODUCTS, make(map[int]models.SimpleProduct))
	ctx = context.WithValue(ctx, RACKS, make(map[int]models.Rack))
	return ctx
}

func get[T any](ctx context.Context, c ContextContent) (T, error) {
	v, ok := ctx.Value(c).(T)
	if !ok {
		return v, ErrKeyNotFound
	}
	return v, nil
}

func GetProducts(ctx context.Context) (map[int]models.SimpleProduct, error) {
	return get[map[int]models.SimpleProduct](ctx, PRODUCTS)
}

func GetRacks(ctx context.Context) (map[int]models.Rack, error) {
	return get[map[int]models.Rack](ctx, RACKS)
}
