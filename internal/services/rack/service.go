package rack

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/utils"
)

type Repository interface {
	GetRacksByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error)
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetRacksHasProductsByProductIds(ctx context.Context, productIds []int) ([]models.RackHasProduct, error)
	GetRackById(ctx context.Context, id int) (models.Rack, error)
	GetAllRackOfProduct(ctx context.Context, productId int) (models.RacksOfProduct, error)
	GetRacksIdsByProductsIds(ctx context.Context, productsIds []int) ([]int, error)
	GetMainRacksIdsByProductsIds(ctx context.Context, productsIds []int) ([]int, error)
	GetMainRacksCountByProductsIds(ctx context.Context, productsIds []int) (int, error)
	GetRacksHasProductsByProductsIdsSplitByIsMain(ctx context.Context, productsIds []int) ([]models.RackIdProductsIds, []models.RackIdProductsIds, error)
}

type Service interface {
	GetRackById(ctx context.Context, id int) (models.Rack, error)
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetMainRacksByProductIds(ctx context.Context, ids []int) ([]models.MainRack, error)
	GetRacksByProductId(ctx context.Context, productId int) ([]models.RackHasProduct, error)
	GetRacksByProductsIds(ctx context.Context, productsIds []int) ([]models.Rack, error)
	GetRacksWithProductByProductId(ctx context.Context, productId int) (models.ProductWithMainRackAndAdditionalRacks, error)
	GetRacksHasProductByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error)
	GetAllRackOfProduct(ctx context.Context, productId int) (models.RacksOfProduct, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetRackById(ctx context.Context, id int) (models.Rack, error) {
	return s.repo.GetRackById(ctx, id)
}

func (s *service) GetAllRackOfProduct(ctx context.Context, productId int) (models.RacksOfProduct, error) {
	return s.repo.GetAllRackOfProduct(ctx, productId)
}

func (s *service) GetRacksWithProductByProductId(ctx context.Context, productId int) (models.ProductWithMainRackAndAdditionalRacks, error) {
	racks, err := s.GetRacksByProductId(ctx, productId)
	if err != nil {
		return models.ProductWithMainRackAndAdditionalRacks{}, nil
	}
	var product models.ProductWithMainRackAndAdditionalRacks
	for _, rack := range racks {
		r, err := s.repo.GetRackById(ctx, rack.RackId)
		if err != nil {
			return models.ProductWithMainRackAndAdditionalRacks{}, err
		}
		if rack.IsMain {
			product.MainRack = models.Rack{
				Id:    rack.RackId,
				Title: r.Title,
			}
		} else {
			product.AdditionalRacks = append(product.AdditionalRacks, models.Rack{
				Id:    rack.RackId,
				Title: r.Title,
			})
		}
	}
	return product, nil
}

func (s *service) GetRacksByProductId(ctx context.Context, productId int) ([]models.RackHasProduct, error) {
	rm, err := utils.GetRacks(ctx)
	if err != nil {
		return []models.RackHasProduct{}, err
	}
	racksHasProducts, err := s.repo.GetRacksByProductId(ctx, productId)
	if err != nil {
		return []models.RackHasProduct{}, err
	}
	var racks []models.RackHasProduct
	for _, rackHasProducts := range racksHasProducts {
		rack, ok := rm[rackHasProducts.Id]
		if !ok {
			rack, err := s.repo.GetRackById(ctx, rackHasProducts.Id)
			if err != nil {
				return []models.RackHasProduct{}, err
			}
			rm[rackHasProducts.Id] = rack
		}
		racks = append(racks, models.RackHasProduct{
			RackId:    rack.Id,
			ProductId: rackHasProducts.Id,
			IsMain:    rackHasProducts.IsMain,
		})
	}
	return racks, nil
}

func (s *service) GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error) {
	return s.repo.GetRacksByIds(ctx, ids)
}

func (s *service) GetMainRacksByProductIds(ctx context.Context, ids []int) ([]models.MainRack, error) {
	var racksMap map[int]models.Rack
	{
		racks, err := s.GetRacksByProductsIds(ctx, ids)
		if err != nil {
			return []models.MainRack{}, err
		}
		racksMap = make(map[int]models.Rack, len(racks))
		for _, rack := range racks {
			racksMap[rack.Id] = rack
		}
	}

	mainRacksProducts, additionalRacks, err := s.repo.GetRacksHasProductsByProductsIdsSplitByIsMain(ctx, ids)
	if err != nil {
		return []models.MainRack{}, err
	}

	var mainRacksProduct map[int][]int
	{
		mainRacksProduct = make(map[int][]int, len(mainRacksProducts))
		for _, rack := range mainRacksProducts {
			mainRacksProduct[rack.RackId] = append(mainRacksProduct[rack.RackId], rack.ProductsIds...)
		}
	}
	additionalRacksOfProduct := make(map[int][]models.Rack, len(ids))
	{
		for _, r := range additionalRacks {
			rack, exists := racksMap[r.RackId]
			if !exists {
				rack, err := s.repo.GetRackById(ctx, r.RackId)
				if err != nil {
					return []models.MainRack{}, err
				}
				racksMap[r.RackId] = rack
			}
			for _, pId := range r.ProductsIds {
				additionalRacksOfProduct[pId] = append(additionalRacksOfProduct[pId], rack)
			}
		}
	}
	var racks []models.MainRack
	for rackId := range mainRacksProduct {
		rr := racksMap[rackId]
		var k models.MainRack = models.MainRack{
			Id:    rr.Id,
			Title: rr.Title,
		}
		for _, productId := range mainRacksProduct[rackId] {
			additionalRacks := additionalRacksOfProduct[productId]
			k.Products = append(k.Products, struct {
				Id              int
				AdditionalRacks []models.Rack
			}{
				Id:              productId,
				AdditionalRacks: additionalRacks,
			})
		}
		racks = append(racks, k)
	}
	return racks, nil
}

func (s *service) GetRacksHasProductByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error) {
	return s.repo.GetRacksByProductId(ctx, productId)
}

func (s *service) GetRacksByProductsIds(ctx context.Context, productsIds []int) ([]models.Rack, error) {
	racksIds, err := s.repo.GetRacksIdsByProductsIds(ctx, productsIds)
	if err != nil {
		return []models.Rack{}, err
	}
	return s.repo.GetRacksByIds(ctx, racksIds)
}

func GetMainRacksByProductIds(s *service, ctx context.Context, ids []int) ([]models.MainRack, error) {
	racksHasProducts, err := s.repo.GetRacksHasProductsByProductIds(ctx, ids)
	if err != nil {
		return []models.MainRack{}, err
	}
	var racksMap map[int]models.Rack
	{
		racks, err := s.GetRacksByProductsIds(ctx, ids)
		if err != nil {
			return []models.MainRack{}, err
		}
		racksMap = make(map[int]models.Rack, len(racks))
		for _, rack := range racks {
			racksMap[rack.Id] = rack
		}
	}
	products := make(map[int][]models.Rack)
	mainRacks := make(map[int]struct {
		Title    string
		Products []int
	})
	for _, r := range racksHasProducts {
		if r.IsMain {
			mainRack, ok := mainRacks[r.RackId]
			if !ok {
				newRack, exists := racksMap[r.RackId]
				if !exists {
					return []models.MainRack{}, fmt.Errorf("rack service: rack (id=%v) not found in rack map", r.RackId)
				}
				mainRacks[r.RackId] = struct {
					Title    string
					Products []int
				}{
					Title:    newRack.Title,
					Products: []int{r.ProductId},
				}
			} else {
				mainRack.Products = append(mainRack.Products, r.ProductId)
				mainRacks[r.RackId] = mainRack
			}
			_, ok = products[r.ProductId]
			if !ok {
				products[r.ProductId] = []models.Rack{}
			}
		} else {
			newRack, exists := racksMap[r.RackId]
			if !exists {
				return []models.MainRack{}, fmt.Errorf("rack service: rack (id=%v) not found in rack map", r.RackId)
			}
			products[r.ProductId] = append(products[r.ProductId], newRack)
		}
	}
	var racks []models.MainRack
	for k, v := range mainRacks {
		p := []struct {
			Id              int
			AdditionalRacks []models.Rack
		}{}
		for _, i := range v.Products {
			var pp struct {
				Id              int
				AdditionalRacks []models.Rack
			}
			pp.Id = i
			additionalRacks, ok := products[i]
			if !ok {
				return []models.MainRack{}, errors.New("additional rack not found")
			}
			pp.AdditionalRacks = append(pp.AdditionalRacks, additionalRacks...)
			p = append(p, pp)
		}
		racks = append(racks, models.MainRack{
			Id:       k,
			Title:    v.Title,
			Products: p,
		})
	}
	return racks, nil
}
