package structs

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	Id        uuid.UUID
	IdProduct uuid.UUID
	IdOrder   uuid.UUID
	Amount    int
}

type Order struct {
	Id      uuid.UUID
	Date    time.Time
	IdUser  uuid.UUID
	Address string
	Status  string
	Price   float64
}
