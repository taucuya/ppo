package structs

import (
	"time"

	"github.com/google/uuid"
)

type BasketItem struct {
	Id        uuid.UUID
	IdProduct uuid.UUID
	IdBasket  uuid.UUID
	Amount    int
}

type Basket struct {
	Id     uuid.UUID
	IdUser uuid.UUID
	Date   time.Time
}
