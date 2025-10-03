package worker_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type WorkerRepositoryInterface interface {
	Create(ctx context.Context, w structs.Worker) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Worker, error)
	GetOrders(ctx context.Context, id uuid.UUID) ([]structs.Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllWorkers(ctx context.Context) ([]structs.Worker, error)
	AcceptOrder(ctx context.Context, id_order uuid.UUID, id_user uuid.UUID) error
}
