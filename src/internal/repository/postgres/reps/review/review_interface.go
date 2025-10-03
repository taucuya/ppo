package review_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type ReviewRepositoryInterface interface {
	Create(ctx context.Context, r structs.Review) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Review, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ReviewsForProduct(ctx context.Context, id_product uuid.UUID) ([]structs.Review, error)
}
