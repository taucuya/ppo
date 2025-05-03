package structs

import "github.com/google/uuid"

type Brand struct {
	Id            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	PriceCategory string    `db:"price_category"`
}
