package basket

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type BasketService interface {
	Create(ctx context.Context, b structs.Basket) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Basket, error)
	GetItems(ctx context.Context, id_basket uuid.UUID) ([]structs.BasketItem, error)
	AddItem(ctx context.Context, i structs.BasketItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	UpdateItemAmount(ctx context.Context, id uuid.UUID, amount int) error
}

type BasketRepository interface {
	Create(ctx context.Context, b structs.Basket) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Basket, error)
	GetItems(ctx context.Context, id_basket uuid.UUID) ([]structs.BasketItem, error)
	AddItem(ctx context.Context, i structs.BasketItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	UpdateItemAmount(ctx context.Context, id uuid.UUID, amount int) error
}

type Service struct {
	rep BasketRepository
}

func New(rep BasketRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, b structs.Basket) error {
	err := s.rep.Create(ctx, b)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Basket, error) {
	b, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.Basket{}, err
	}
	return b, nil
}

func (s *Service) GetItems(ctx context.Context, id_basket uuid.UUID) ([]structs.BasketItem, error) {
	arr, err := s.rep.GetItems(ctx, id_basket)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (s *Service) AddItem(ctx context.Context, i structs.BasketItem) error {
	// id := structs.GenId()
	err := s.rep.AddItem(ctx, i)
	return err
}

func (s *Service) DeleteItem(ctx context.Context, id uuid.UUID) error {
	err := s.rep.DeleteItem(ctx, id)
	return err
}

func (s *Service) UpdateItemAmount(ctx context.Context, id uuid.UUID, amount int) error {
	err := s.rep.UpdateItemAmount(ctx, id, amount)
	return err
}
