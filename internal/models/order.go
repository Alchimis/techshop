package models

type Order struct {
	Id       int
	Products []struct {
		Product  Product
		Quantity int
	}
}
