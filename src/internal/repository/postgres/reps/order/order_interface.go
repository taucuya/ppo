package order_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, o structs.Order) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Order, error)
	GetItems(ctx context.Context, id uuid.UUID) ([]structs.OrderItem, error)
	GetFreeOrders(ctx context.Context) ([]structs.Order, error)
	GetOrdersByUser(ctx context.Context, id uuid.UUID) ([]structs.Order, error)
	GetStatus(ctx context.Context, id uuid.UUID) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
