package models

type Product struct {
	Id              int
	MainRackId      Rack
	Title           string
	ProductCategory string
	AdditionalRacks []Rack
}

type SimpleProduct struct {
	Id    int
	Title string
}

type ProductOrder struct {
	Id       int    `json:"product_id"`
	OrderId  int    `json:"order_id"`
	Quantity int    `json:"order_quantity"`
	Title    string `json:"product_title"`
}

type ProductWithMainRackAndAdditionalRacks struct {
	MainRack        Rack
	AdditionalRacks []Rack
}

type ProductIdAndQuantity struct {
	Id       int
	Quantity int
}

type OrderIdAndQuantity struct {
	Id       int
	Quantity int
}

type ProductIdAndOrders struct {
	ProductId int
	Orders    []OrderIdAndQuantity
}
