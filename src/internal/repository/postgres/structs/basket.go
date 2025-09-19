package structs

import (
	"time"

	"github.com/google/uuid"
)

type BasketItem struct {
	Id        uuid.UUID `db:"id"`
	IdProduct uuid.UUID `db:"id_product"`
	IdBasket  uuid.UUID `db:"id_basket"`
	Amount    int       `db:"amount"`
}

type Basket struct {
	Id     uuid.UUID `db:"id"`
	IdUser uuid.UUID `db:"id_user"`
	Date   time.Time `db:"date"`
}
