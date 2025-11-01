package structs

import (
	"errors"

	"github.com/google/uuid"
)

type Brand struct {
	Id            uuid.UUID
	Name          string
	Description   string
	PriceCategory string
}

var (
	ErrBrandNotFound = errors.New("brand not found")
	ErrNotFound      = errors.New("not found")
	ErrNoRows        = errors.New("no rows in result set")
)
