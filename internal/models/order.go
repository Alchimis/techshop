package models

type Order struct {
	Id       int
	Products []struct {
		Product  Product
		Quantity int
	}
}

type OrderHasProduct struct {
	OrderId         int
	ProductId       int
	ProductQuantity int
}
