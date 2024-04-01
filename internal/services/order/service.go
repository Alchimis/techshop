package order

import (
	"context"
	"fmt"

	errs "errors"

	"github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/product"
	"github.com/Alchimis/techshop/internal/services/rack"
	"github.com/Alchimis/techshop/internal/utils"
)

type Repository interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
	GetOrderHasProductByIds(ctx context.Context, ids []int) ([]models.OrderHasProduct, error)
}

type service struct {
	repo           Repository
	productService product.Service
	rackService    rack.Service
}

type Service interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
	GetOrdersByIdsSortedByMainRacks(ctx context.Context, ind []int) ([]models.RackWithProducts, error)
}

func NewService(repo Repository, productService product.Service, rackService rack.Service) Service {
	return &service{
		repo:           repo,
		productService: productService,
		rackService:    rackService,
	}
}

func (s *service) GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error) {
	return []models.Order{}, errors.ErrNotImplemented
}

func (s *service) GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error) {
	ordersHasProducts, err := s.repo.GetOrderHasProductByIds(ctx, ids)
	if err != nil {
		return []models.RackWithProducts{}, err
	}
	productsIdsWithOrders := map[int][]struct {
		Id       int
		Quantity int
	}{}
	productIds := []int{}
	for _, o := range ordersHasProducts {
		order, ok := productsIdsWithOrders[o.ProductId]
		if !ok {
			productsIdsWithOrders[o.ProductId] = []struct {
				Id       int
				Quantity int
			}{{
				Id:       o.OrderId,
				Quantity: o.ProductQuantity,
			}}
			productIds = append(productIds, o.ProductId)
		} else {
			order := append(order, struct {
				Id       int
				Quantity int
			}{Id: o.OrderId, Quantity: o.ProductQuantity})
			productsIdsWithOrders[o.ProductId] = order
		}
	}
	var products map[int]string = make(map[int]string)
	{
		p, err := s.productService.GetProductsByIds(ctx, productIds)
		if err != nil {
			return []models.RackWithProducts{}, err
		}
		for _, pp := range p {
			products[pp.Id] = pp.Title
		}
	}
	racks, err := s.rackService.GetMainRacksByProductIds(ctx, productIds)
	if err != nil {
		return []models.RackWithProducts{}, err
	}
	var racksWithProducts []models.RackWithProducts
	for _, rack := range racks {
		var p []models.ProductIn
		for _, productInRack := range rack.Products {
			title, ok := products[productInRack.Id]
			if !ok {
				return []models.RackWithProducts{}, errs.New(fmt.Sprintf("title of product (id=%v) not found", productInRack.Id))
			}
			orders, ok := productsIdsWithOrders[productInRack.Id]
			if !ok {
				return []models.RackWithProducts{}, errs.New(fmt.Sprintf("order of product (id=%v) not found", productInRack.Id))
			}
			additionalRacks := utils.Map(productInRack.AdditionalRacks, func(r models.Rack) struct {
				RackName *string `json:"rack_name"`
				RackId   *int    `json:"rack_id"`
			} {
				return struct {
					RackName *string `json:"rack_name"`
					RackId   *int    `json:"rack_id"`
				}{
					RackName: &r.Title,
					RackId:   &r.Id,
				}
			})
			for _, order := range orders {
				var pp models.ProductIn = models.ProductIn{
					Id:              productInRack.Id,
					Title:           title,
					OrderId:         order.Id,
					Quantity:        order.Quantity,
					AdditionalRacks: additionalRacks,
				}
				p = append(p, pp)
			}
		}
		var racksWithProduct models.RackWithProducts = models.RackWithProducts{
			Id:       rack.Id,
			Name:     rack.Title,
			Products: p,
		}
		racksWithProducts = append(racksWithProducts, racksWithProduct)
	}
	return racksWithProducts, nil //s.repo.GetOrdersByIdSortByRacks(ctx, ids)
}

var mapAdditionalRacks = func(r models.Rack) struct {
	RackName *string `json:"rack_name"`
	RackId   *int    `json:"rack_id"`
} {
	return struct {
		RackName *string `json:"rack_name"`
		RackId   *int    `json:"rack_id"`
	}{
		RackName: &r.Title,
		RackId:   &r.Id,
	}
}

func (s *service) GetOrdersByIdsSortedByMainRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error) {
	mainRacks := make(map[int]models.RackWithProducts)
	ctx = utils.SetupContext(ctx)
	for _, orderId := range ids {
		//productOrders, err := s.productService.GetProductsOrdersByOrderId(ctx, orderId)
		productsQuantity, err := s.productService.GetProductIdAndQuantityByOrderId(ctx, orderId)
		if err != nil {
			return []models.RackWithProducts{}, err
		}
		for _, productOrder := range productsQuantity {
			product, err := s.productService.GetProductById(ctx, productOrder.Id)
			if err != nil {
				return []models.RackWithProducts{}, err
			}
			//racksOfProduct, err := s.rackService.GetRacksHasProductByProductId(ctx, productOrder.Id)
			if err != nil {
				return []models.RackWithProducts{}, err
			}
			racks, err := s.rackService.GetRacksWithProductByProductId(ctx, productOrder.Id)
			if err != nil {
				return []models.RackWithProducts{}, err
			}

			if mainRack, ok := mainRacks[racks.MainRack.Id]; !ok {
				mainRacks[racks.MainRack.Id] = models.RackWithProducts{
					Id:   racks.MainRack.Id,
					Name: racks.MainRack.Title,
					Products: []models.ProductIn{
						{
							Id:              productOrder.Id,
							OrderId:         orderId,
							Quantity:        productOrder.Quantity,
							Title:           product.Title,
							AdditionalRacks: utils.Map(racks.AdditionalRacks, mapAdditionalRacks),
						},
					},
				}
			} else {
				mainRack.Products = append(mainRack.Products, models.ProductIn{
					Id:              productOrder.Id,
					OrderId:         orderId,
					Quantity:        productOrder.Quantity,
					Title:           product.Title,
					AdditionalRacks: utils.Map(racks.AdditionalRacks, mapAdditionalRacks),
				})
				mainRacks[racks.MainRack.Id] = mainRack
			}
		}
	}
	var r []models.RackWithProducts
	for _, v := range mainRacks {
		r = append(r, v)
	}
	return r, nil
}
