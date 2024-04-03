package order

import (
	"context"
	"fmt"

	domainErrors "github.com/Alchimis/techshop/internal/errors"
	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/services/product"
	"github.com/Alchimis/techshop/internal/services/rack"
	"github.com/Alchimis/techshop/internal/utils"
)

type Repository interface {
	GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error)
	GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error)
	GetOrderHasProductByIds(ctx context.Context, ids []int) ([]models.OrderHasProduct, error)
	GetOrderHasProductByOrdersIdsGroupByProduct(ctx context.Context, ids []int) ([]models.ProductIdAndOrders, error)
	HelloWorld(ctx context.Context) error
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
	HelloWorld(ctx context.Context) error
}

func NewService(repo Repository, productService product.Service, rackService rack.Service) Service {
	return &service{
		repo:           repo,
		productService: productService,
		rackService:    rackService,
	}
}

func (s *service) HelloWorld(ctx context.Context) error {
	return s.repo.HelloWorld(ctx)
}

func (s *service) GetOrdersById(ctx context.Context, ids []int) ([]models.Order, error) {
	return []models.Order{}, domainErrors.ErrNotImplemented
}

func (s *service) GetOrdersByIdSortByRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error) {
	ordersHasProducts, err := s.repo.GetOrderHasProductByOrdersIdsGroupByProduct(ctx, ids) //s.repo.GetOrderHasProductByIds(ctx, ids)

	if err != nil {
		return []models.RackWithProducts{}, err
	}
	productsIdsWithOrders := make(map[int][]models.OrderIdAndQuantity, len(ordersHasProducts))
	productIds := []int{}
	for _, o := range ordersHasProducts {
		productsIdsWithOrders[o.ProductId] = append(productsIdsWithOrders[o.ProductId], o.Orders...)
		productIds = append(productIds, o.ProductId)
	}
	var products map[int]string
	{
		p, err := s.productService.GetProductsByIds(ctx, productIds)
		products = make(map[int]string, len(p))
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
				return []models.RackWithProducts{}, fmt.Errorf("title of product (id=%v) not found", productInRack.Id)
			}
			orders, ok := productsIdsWithOrders[productInRack.Id]
			if !ok {
				return []models.RackWithProducts{}, fmt.Errorf("order of product (id=%v) not found", productInRack.Id)
			}
			for _, order := range orders {
				var pp models.ProductIn = models.ProductIn{
					Id:              productInRack.Id,
					Title:           title,
					OrderId:         order.Id,
					Quantity:        order.Quantity,
					AdditionalRacks: productInRack.AdditionalRacks,
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
	return racksWithProducts, nil
}

func (s *service) GetOrdersByIdsSortedByMainRacks(ctx context.Context, ids []int) ([]models.RackWithProducts, error) {
	mainRacks := make(map[int]models.RackWithProducts)
	racksMap := make(map[int]models.Rack)
	productsMap := make(map[int]struct {
		product         models.SimpleProduct
		mainRack        models.Rack
		additionalRacks []models.Rack
	})
	for _, orderId := range ids {
		productsQuantity, err := s.productService.GetProductIdAndQuantityByOrderId(ctx, orderId)
		if err != nil {
			return []models.RackWithProducts{}, err
		}
		for _, productOrder := range productsQuantity {
			p, ok := productsMap[productOrder.Id]
			if !ok {
				product, err := s.productService.GetProductById(ctx, productOrder.Id)
				if err != nil {
					return []models.RackWithProducts{}, err
				}
				racksOfProduct, err := s.rackService.GetAllRackOfProduct(ctx, productOrder.Id)
				if err != nil {
					return []models.RackWithProducts{}, err
				}
				var (
					r models.Rack
				)
				if r, ok = racksMap[racksOfProduct.MainRackId]; !ok {
					r, err = s.rackService.GetRackById(ctx, racksOfProduct.MainRackId)
					if err != nil {
						return []models.RackWithProducts{}, err
					}
					racksMap[racksOfProduct.MainRackId] = r
				}

				additionalRacks, err := utils.MapWithError(racksOfProduct.AdditionalRacks, func(id int) (models.Rack, error) {
					var (
						r models.Rack
					)
					if r, ok = racksMap[racksOfProduct.MainRackId]; !ok {
						r, err = s.rackService.GetRackById(ctx, racksOfProduct.MainRackId)
						if err != nil {
							return models.Rack{}, err
						}
						racksMap[racksOfProduct.MainRackId] = r
					}
					return r, nil
				})

				if err != nil {
					return []models.RackWithProducts{}, err
				}
				p = struct {
					product         models.SimpleProduct
					mainRack        models.Rack
					additionalRacks []models.Rack
				}{
					product:         product,
					mainRack:        r,
					additionalRacks: additionalRacks,
				}
				productsMap[productOrder.Id] = p
			}
			mainRack, ok := mainRacks[p.mainRack.Id]
			if !ok {
				mainRacks[p.mainRack.Id] = models.RackWithProducts{
					Id:   p.mainRack.Id,
					Name: p.mainRack.Title,
					Products: []models.ProductIn{{
						Id:              p.product.Id,
						Title:           p.product.Title,
						Quantity:        productOrder.Quantity,
						OrderId:         productOrder.Id,
						AdditionalRacks: p.additionalRacks,
					}},
				}
			} else {
				mainRack.Products = append(mainRack.Products, models.ProductIn{
					Id:              p.product.Id,
					Title:           p.product.Title,
					Quantity:        productOrder.Quantity,
					OrderId:         productOrder.Id,
					AdditionalRacks: p.additionalRacks,
				})
				mainRacks[p.mainRack.Id] = mainRack
			}
		}
	}
	var r []models.RackWithProducts
	for _, v := range mainRacks {
		r = append(r, v)
	}
	return r, nil
}
