package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type ProductService interface {
	Create(ctx context.Context, p structs.Product) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Product, error)
	GetByCategory(ctx context.Context, category string) ([]structs.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductRepository interface {
	Create(ctx context.Context, p structs.Product) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Product, error)
	GetByCategory(ctx context.Context, category string) ([]structs.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
type Service struct {
	rep ProductRepository
}

func New(rep ProductRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, p structs.Product) error {
	err := s.rep.Create(ctx, p)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Product, error) {
	p, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.Product{}, err
	}
	return p, nil
}

func (s *Service) GetByCategory(ctx context.Context, category string) ([]structs.Product, error) {
	p, err := s.rep.GetByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.rep.Delete(ctx, id)
	return err
}
