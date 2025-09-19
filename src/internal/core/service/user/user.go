package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
)

type UserService interface {
	Create(ctx context.Context, u structs.User) error
	GetById(ctx context.Context, id uuid.UUID) (structs.User, error)
	GetByMail(ctx context.Context, mail string) (structs.User, error)
	GetAllUsers(ctx context.Context) ([]structs.User, error)
	GetByPhone(ctx context.Context, phone string) (structs.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, u structs.User) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (structs.User, error)
	GetByMail(ctx context.Context, mail string) (structs.User, error)
	GetAllUsers(ctx context.Context) ([]structs.User, error)
	GetByPhone(ctx context.Context, phone string) (structs.User, error)
}

type UsrBasket interface {
	Create(ctx context.Context, u structs.Basket) error
}

type UsrFavourites interface {
	Create(ctx context.Context, u structs.Favourites) error
}

type Service struct {
	rep UserRepository
	bsk UsrBasket
	fav UsrFavourites
}

func New(rep UserRepository, bsk UsrBasket, fav UsrFavourites) *Service {
	return &Service{rep: rep, bsk: bsk, fav: fav}
}

func (s *Service) Create(ctx context.Context, u structs.User) error {
	id, err := s.rep.Create(ctx, u)
	if err != nil {
		return err
	}
	basket := structs.Basket{
		IdUser: id,
		Date:   time.Now(),
	}
	err = s.bsk.Create(ctx, basket)
	if err != nil {
		return err
	}
	favourites := structs.Favourites{
		IdUser: id,
	}
	err = s.fav.Create(ctx, favourites)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (structs.User, error) {
	u, err := s.rep.GetById(ctx, id)
	if err != nil {
		return structs.User{}, err
	}
	return u, nil
}

func (s *Service) GetByMail(ctx context.Context, mail string) (structs.User, error) {
	u, err := s.rep.GetByMail(ctx, mail)
	if err != nil {
		return structs.User{}, err
	}
	return u, nil
}

func (s *Service) GetAllUsers(ctx context.Context) ([]structs.User, error) {
	u, err := s.rep.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetByPhone(ctx context.Context, phone string) (structs.User, error) {
	u, err := s.rep.GetByPhone(ctx, phone)
	if err != nil {
		return structs.User{}, err
	}
	return u, nil
}
