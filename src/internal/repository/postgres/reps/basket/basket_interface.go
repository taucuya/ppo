package basket_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type BasketRepositoryInterface interface {
	Create(ctx context.Context, b structs.Basket) error
	GetBIdByUId(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (structs.Basket, error)
	GetItems(ctx context.Context, id_user uuid.UUID) ([]structs.BasketItem, error)
	AddItem(ctx context.Context, i structs.BasketItem) error
	DeleteItem(ctx context.Context, id uuid.UUID, product_id uuid.UUID) error
	UpdateItemAmount(ctx context.Context, basket_id uuid.UUID, product_id uuid.UUID, amount int) error
}
