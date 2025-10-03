package favourites

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type FavouritesService interface {
	Create(ctx context.Context, f structs.Favourites) error
	GetById(ctx context.Context, id uuid.UUID) (structs.Favourites, error)
	GetItems(ctx context.Context, id_Favourites uuid.UUID) ([]structs.FavouritesItem, error)
	AddItem(ctx context.Context, i structs.FavouritesItem) error
	DeleteItem(ctx context.Context, id uuid.UUID, product_id uuid.UUID) error
}

type FavouritesRepository interface {
	Create(ctx context.Context, f structs.Favourites) error
	GetFIdByUId(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (structs.Favourites, error)
	GetItems(ctx context.Context, id_favourites uuid.UUID) ([]structs.FavouritesItem, error)
	AddItem(ctx context.Context, i structs.FavouritesItem) error
	DeleteItem(ctx context.Context, id uuid.UUID, product_id uuid.UUID) error
}

type Service struct {
	rep FavouritesRepository
}

func New(rep FavouritesRepository) *Service {
	return &Service{rep: rep}
}

func (s *Service) Create(ctx context.Context, f structs.Favourites) error {
	err := s.rep.Create(ctx, f)
	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.Favourites, error) {
	bid, err := s.rep.GetFIdByUId(ctx, id)
	if err != nil {
		return structs.Favourites{}, err
	}
	b, err := s.rep.GetById(ctx, bid)
	if err != nil {
		return structs.Favourites{}, err
	}
	return b, nil
}

func (s *Service) GetItems(ctx context.Context, id_user uuid.UUID) ([]structs.FavouritesItem, error) {
	arr, err := s.rep.GetItems(ctx, id_user)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (s *Service) AddItem(ctx context.Context, i structs.FavouritesItem, id uuid.UUID) error {
	fid, err := s.rep.GetFIdByUId(ctx, id)
	if err != nil {
		return err
	}

	i.IdFavourites = fid
	err = s.rep.AddItem(ctx, i)
	return err
}

func (s *Service) DeleteItem(ctx context.Context, id uuid.UUID, id_item uuid.UUID) error {
	fid, err := s.rep.GetFIdByUId(ctx, id)
	if err != nil {
		return err
	}
	err = s.rep.DeleteItem(ctx, fid, id_item)
	return err
}
