package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type OrderService interface {
	Create(ctx context.Context, o structs.Order) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Order, error)
	GetItems(ctx context.Context, id uuid.UUID) ([]structs.OrderItems, error)
	GetStatus(ctx context.Context, id uuid.UUID) (string, error)
	ChangeOrderStatus(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type OrderRepository interface {
	Create(ctx context.Context, o structs.Order) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Order, error)
	GetItems(ctx context.Context, id uuid.UUID) ([]structs.OrderItems, error)
	GetStatus(ctx context.Context, id uuid.UUID) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type Service struct {
	rep OrderRepository
}

func New(rep OrderRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, o structs.Order) error {
	err := s.rep.Create(ctx, o)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Order, error) {
	o, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.Order{}, err
	}
	return o, nil
}

func (s *Service) GetItems(ctx context.Context, id uuid.UUID) ([]structs.OrderItems, error) {
	items, err := s.rep.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) GetStatus(ctx context.Context, id uuid.UUID) (string, error) {
	str, err := s.rep.GetStatus(ctx, id)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (s *Service) ChangeOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	stat, err := s.rep.GetStatus(ctx, id)
	if err != nil {
		return err
	}
	if stat == status {
		return nil
	} else {
		_, err := s.rep.GetById(ctx, id)
		if err != nil {
			return err
		} else {
			err := s.rep.UpdateStatus(ctx, id, status)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.rep.Delete(ctx, id)
	return err
}
