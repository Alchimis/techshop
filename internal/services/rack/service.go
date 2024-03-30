package rack

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alchimis/techshop/internal/models"
)

type Repository interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetRacksHasProductsByProductIds(ctx context.Context, productIds []int) ([]models.RackHasProduct, error)
	GetRackById(ctx context.Context, id int) (models.Rack, error)
}

type Service interface {
	GetRacksByIds(ctx context.Context, ids []int) ([]models.Rack, error)
	GetMainRacksByProductIds(ctx context.Context, ids []int) ([]models.MainRack, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
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
					Title: newRack.Title,
				}
			} else {
				mainRack.Products = append(mainRack.Products, r.ProductId)
				mainRacks[r.RackId] = mainRack
			}
		}
		newRack, err := s.repo.GetRackById(ctx, r.RackId)
		if err != nil {
			return []models.MainRack{}, err
		}
		products[r.ProductId] = append(products[r.ProductId], newRack)
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
			for _, ad := range additionalRacks {
				pp.AdditionalRacks = append(pp.AdditionalRacks, ad)
			}
			p = append(p, pp)
		}
		racks = append(racks, models.MainRack{
			Id:       k,
			Title:    v.Title,
			Products: p,
		})
	}
	fmt.Println("racksHasProducts ", racksHasProducts)
	return racks, nil
}
