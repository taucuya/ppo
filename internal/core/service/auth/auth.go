package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/structs"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(ctx context.Context, u structs.User) error
	LogIn(ctx context.Context, mail string, password string) (atoken string, rtoken string, err error)
	LogOut(ctx context.Context, id uuid.UUID) error
	VerifyAToken(ctx context.Context, token string) error
	VerifyRToken(ctx context.Context, token string) (uuid.UUID, error)
	RefreshToken(ctx context.Context, atoken string, rtoken string) (string, error)
}

type AuthRepository interface {
	CreateToken(ctx context.Context, id uuid.UUID, rtoken string) error
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteToken(ctx context.Context, id uuid.UUID) error
}

type AuthUser interface {
	Create(ctx context.Context, u structs.User) error
	GetByMail(ctx context.Context, mail string) (structs.User, error)
}

type AuthProvider interface {
	GenToken(ctx context.Context, id uuid.UUID) (atoken string, rtoken string, err error)
	VerifyToken(ctx context.Context, token string) (bool, error)
	RefreshToken(ctx context.Context, atoken string, rtoken string) (string, error)
}

type Service struct {
	prov AuthProvider
	rep  AuthRepository
	usr  AuthUser
}

func New(prov AuthProvider, rep AuthRepository, usr AuthUser) *Service {
	return &Service{prov: prov, rep: rep, usr: usr}
}

func (s *Service) SignIn(ctx context.Context, u structs.User) error {
	return s.usr.Create(ctx, u)

}

func (s *Service) LogIn(ctx context.Context, mail string, password string) (atoken string, rtoken string, err error) {
	u, err := s.usr.GetByMail(ctx, mail)
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", "", err
	}
	atoken, rtoken, err = s.prov.GenToken(ctx, u.Id)
	if err != nil {
		return "", "", err
	}
	err = s.rep.CreateToken(ctx, u.Id, rtoken)
	if err != nil {
		return "", "", err
	}
	return atoken, rtoken, err
}

func (s *Service) LogOut(ctx context.Context, rtoken string) error {
	id, _, err := s.VerifyRToken(ctx, rtoken)
	if err != nil {
		return err
	}
	err = s.rep.DeleteToken(ctx, id)
	return err
}

func (s *Service) VerifyAToken(ctx context.Context, token string) (bool, error) {
	expired, err := s.prov.VerifyToken(ctx, token)
	if err != nil {
		return expired, err
	}
	return expired, nil
}

func (s *Service) VerifyRToken(ctx context.Context, token string) (uuid.UUID, bool, error) {
	expired, err := s.prov.VerifyToken(ctx, token)
	if err != nil {
		return uuid.UUID{}, expired, err
	}
	id, err := s.rep.VerifyToken(ctx, token)
	if err != nil {
		return uuid.UUID{}, expired, err
	}
	return id, expired, nil
}

func (s *Service) RefreshToken(ctx context.Context, atoken string, rtoken string) (string, error) {
	token, err := s.prov.RefreshToken(ctx, atoken, rtoken)
	if err != nil {
		return "", err
	}
	return token, nil
}
