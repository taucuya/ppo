package structs

import "github.com/google/uuid"

type Brand struct {
	Id            uuid.UUID
	Name          string
	Description   string
	PriceCategory string
}
