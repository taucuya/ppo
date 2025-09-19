package structs

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	Id        uuid.UUID `db:"id"`
	IdProduct uuid.UUID `db:"id_product"`
	IdOrder   uuid.UUID `db:"id_order"`
	Amount    int       `db:"amount"`
}

type Order struct {
	Id      uuid.UUID `db:"id"`
	Date    time.Time `db:"date"`
	IdUser  uuid.UUID `db:"id_user"`
	Address string    `db:"address"`
	Status  string    `db:"status"`
	Price   float64   `db:"price"`
}
