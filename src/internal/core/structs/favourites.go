package structs

import (
	"errors"

	"github.com/google/uuid"
)

type FavouritesItem struct {
	Id           uuid.UUID
	IdProduct    uuid.UUID
	IdFavourites uuid.UUID
}

type Favourites struct {
	Id     uuid.UUID
	IdUser uuid.UUID
}

var (
	ErrFavouritesNotFound = errors.New("favourites not found")
	ErrItemNotFound       = errors.New("item not found")
	ErrDuplicateItem      = errors.New("item already exists")
)
