package structs

import (
	"github.com/google/uuid"
)

type FavouritesItem struct {
	Id           uuid.UUID `db:"id"`
	IdProduct    uuid.UUID `db:"id_product"`
	IdFavourites uuid.UUID `db:"id_favourites"`
}

type Favourites struct {
	Id     uuid.UUID `db:"id"`
	IdUser uuid.UUID `db:"id_user"`
}
