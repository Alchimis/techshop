package models

import (
	"fmt"
	"strings"
)

type Rack struct {
	Id    int
	Title string
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
		Id                    int
		AdditionalRacksTitles []string
	}
}

type RackWithProducts struct {
	Id       int
	Name     string
	Products []struct {
		Id              int    `json:"product_id"`
		OrderId         int    `json:"order_id"`
		Quantity        int    `json:"order_quantity"`
		Title           string `json:"product_title"`
		AdditionalRacks []struct {
			RackName *string `json:"rack_name"`
			RackId   *int    `json:"rack_id"`
		} `json:"additional_racks"`
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
			if m.RackName != nil {
				racksNames = append(racksNames, *m.RackName)
			}
		}
		if len(racksNames) > 0 {
			builder.WriteString(fmt.Sprintf("доп стелаж: %s\n", strings.Join(racksNames, ",")))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
