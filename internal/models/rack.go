package models

type Rack struct {
	Id    int
	Title string
}

type RackWithProducts struct {
	Id       int
	Products struct {
		Id                 int
		OrderId            int
		Title              string
		AdditionalRacksIds []int
	}
}
