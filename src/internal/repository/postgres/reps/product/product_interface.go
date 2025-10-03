package product_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type ProductRepositoryInterface interface {
	Create(ctx context.Context, p structs.Product) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Product, error)
	GetByName(ctx context.Context, name string) (structs.Product, error)
	GetByArticule(ctx context.Context, art string) (structs.Product, error)
	GetByCategory(ctx context.Context, category string) ([]structs.Product, error)
	GetByBrand(ctx context.Context, brand string) ([]structs.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
