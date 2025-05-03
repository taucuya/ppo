package brand

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type BrandService interface {
	Create(ctx context.Context, b structs.Brand) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Brand, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllBrands(ctx context.Context) ([]structs.Brand, error)
}

type BrandRepository interface {
	Create(ctx context.Context, b structs.Brand) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Brand, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllBrands(ctx context.Context) ([]structs.Brand, error)
}

type Service struct {
	rep BrandRepository
}

func New(rep BrandRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, b structs.Brand) error {
	err := s.rep.Create(ctx, b)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Brand, error) {
	b, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.Brand{}, err
	}
	return b, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.rep.Delete(ctx, id)
	return err
}

func (s *Service) GetAllBrands(ctx context.Context) ([]structs.Brand, error) {
	arr, err := s.rep.GetAllBrands(ctx)
	if err != nil {
		return nil, err
	}
	return arr, nil
}
