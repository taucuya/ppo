package review

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type ReviewService interface {
	Create(ctx context.Context, r structs.Review) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Review, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ReviewsForProduct(ctx context.Context, id_product uuid.UUID) ([]structs.Review, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, r structs.Review) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Review, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ReviewsForProduct(ctx context.Context, id_product uuid.UUID) ([]structs.Review, error)
}

type Service struct {
	rep ReviewRepository
}

func New(rep ReviewRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, r structs.Review) error {
	err := s.rep.Create(ctx, r)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Review, error) {
	r, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.Review{}, err
	}
	return r, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.rep.Delete(ctx, id)
	return err
}

func (s *Service) ReviewsForProduct(ctx context.Context, id_product uuid.UUID) ([]structs.Review, error) {
	r, err := s.rep.ReviewsForProduct(ctx, id_product)
	if err != nil {
		return nil, err
	}
	return r, nil
}
