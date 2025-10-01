package favourites_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type FavouritesRepositoryInterface interface {
	Create(ctx context.Context, f structs.Favourites) error
	GetFIdByUId(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (structs.Favourites, error)
	GetItems(ctx context.Context, id_user uuid.UUID) ([]structs.FavouritesItem, error)
	AddItem(ctx context.Context, i structs.FavouritesItem) error
	DeleteItem(ctx context.Context, id uuid.UUID, product_id uuid.UUID) error
}
