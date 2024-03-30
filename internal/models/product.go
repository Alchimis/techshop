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
