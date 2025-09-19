package structs

import "github.com/google/uuid"

type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	Category    string
	Amount      int
	IdBrand     uuid.UUID
	PicLink     string
	Articule    string
}
