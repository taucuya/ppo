package structs

import (
	"errors"

	"github.com/google/uuid"
)

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

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrReviewNotFound    = errors.New("review not found")
	ErrDuplicateArticule = errors.New("duplicate articule")
)
