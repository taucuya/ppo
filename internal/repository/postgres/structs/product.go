package structs

import "github.com/google/uuid"

type Product struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Category    string    `db:"category"`
	Amount      int       `db:"amount"`
	IdBrand     uuid.UUID `db:"id_brand"`
	PicLink     string    `db:"pic_link"`
	Articule    string    `db:"art"`
}
