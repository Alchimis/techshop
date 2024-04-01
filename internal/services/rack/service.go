package rack

import (
	"context"
	"errors"

	"github.com/Alchimis/techshop/internal/models"
	"github.com/Alchimis/techshop/internal/utils"
)

type Repository interface {
	GetRacksByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error)
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetRacksHasProductsByProductIds(ctx context.Context, productIds []int) ([]models.RackHasProduct, error)
	GetRackById(ctx context.Context, id int) (models.Rack, error)
}

type Service interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetMainRacksByProductIds(ctx context.Context, ids []int) ([]models.MainRack, error)
	GetRacksByProductId(ctx context.Context, productId int) ([]models.RackHasProduct, error)
	GetRacksWithProductByProductId(ctx context.Context, productId int) (models.ProductWithMainRackAndAdditionalRacks, error)
	GetRacksHasProductByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
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
	racksHasProducts, err := s.repo.GetRacksHasProductsByProductIds(ctx, ids)
	if err != nil {
		return []models.MainRack{}, err
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
				newRack, err := s.repo.GetRackById(ctx, r.RackId)
				if err != nil {
					return []models.MainRack{}, err
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
			newRack, err := s.repo.GetRackById(ctx, r.RackId)
			if err != nil {
				return []models.MainRack{}, err
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

func (s *service) GetRacksHasProductByProductId(ctx context.Context, productId int) ([]models.RackWithIsMain, error) {
	return s.repo.GetRacksByProductId(ctx, productId)
}
