package worker

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type WorkerService interface {
	Create(ctx context.Context, w structs.Worker) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Worker, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllWorkers(ctx context.Context) ([]structs.Worker, error)
	GetOrders(ctx context.Context, id uuid.UUID) ([]structs.Order, error)
	AcceptOrder(ctx context.Context, id_order uuid.UUID, id_worker uuid.UUID) error
}

type WorkerRepository interface {
	Create(ctx context.Context, w structs.Worker) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Worker, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllWorkers(ctx context.Context) ([]structs.Worker, error)
	GetOrders(ctx context.Context, id uuid.UUID) ([]structs.Order, error)
	AcceptOrder(ctx context.Context, id_order uuid.UUID, id_worker uuid.UUID) error
}

type Service struct {
	rep WorkerRepository
}

func New(rep WorkerRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, w structs.Worker) error {
	err := s.rep.Create(ctx, w)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Worker, error) {
	w, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.Worker{}, err
	}
	return w, nil
}

func (s *Service) GetOrders(ctx context.Context, id uuid.UUID) ([]structs.Order, error) {
	order, err := s.rep.GetOrders(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.rep.Delete(ctx, id)
	return err
}

func (s *Service) GetAllWorkers(ctx context.Context) ([]structs.Worker, error) {
	workers, err := s.rep.GetAllWorkers(ctx)
	if err != nil {
		return nil, err
	}
	return workers, nil
}

func (s *Service) AcceptOrder(ctx context.Context, id_order uuid.UUID, id uuid.UUID) error {
	err := s.rep.AcceptOrder(ctx, id_order, id)
	return err
}
