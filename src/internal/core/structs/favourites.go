package structs

import (
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
