package brand_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type BrandRepositoryInterface interface {
	Create(ctx context.Context, b structs.Brand) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Brand, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllBrandsInCategory(ctx context.Context, category string) ([]structs.Brand, error)
}
