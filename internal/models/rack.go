package models

import (
	"fmt"
	"strings"
)

type ProductIdRackId struct {
	ProductId int
	RackId    int
}

type RackIdProductsIds struct {
	RackId      int
	ProductsIds []int
}

type RacksOfProduct struct {
	MainRackId      int
	AdditionalRacks []int
}

type Rack struct {
	Id    int    `json:"rack_name"`
	Title string `json:"rack_id"`
}

type RackWithIsMain struct {
	Id     int
	IsMain bool
}

type RackHasProduct struct {
	RackId    int
	ProductId int
	IsMain    bool
}

type MainRack struct {
	Id       int
	Title    string
	Products []struct {
		Id              int
		AdditionalRacks []Rack
	}
}

type ProductIn struct {
	Id              int    `json:"product_id"`
	OrderId         int    `json:"order_id"`
	Quantity        int    `json:"order_quantity"`
	Title           string `json:"product_title"`
	AdditionalRacks []Rack `json:"additional_racks"`
}

type RackWithProducts struct {
	Id       int
	Name     string
	Products []ProductIn
}

type RackWithProductAndAdditionalRacks struct {
	Id       int
	Name     string
	Products []struct {
		ProductOrder
		Additionalracks []Rack
	}
}

func (r RackWithProducts) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("===Стеллаж %s\n", r.Name))
	for _, p := range r.Products {
		builder.WriteString(fmt.Sprintf("%s (id=%d)\n", p.Title, p.Id))
		builder.WriteString(fmt.Sprintf("заказ %d, %d шт\n", p.OrderId, p.Quantity))
		var racksNames []string
		for _, m := range p.AdditionalRacks {
			if m.Title != "" {
				racksNames = append(racksNames, m.Title)
			}
		}
		if len(racksNames) > 0 {
			builder.WriteString(fmt.Sprintf("доп стелаж: %s\n", strings.Join(racksNames, ",")))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
