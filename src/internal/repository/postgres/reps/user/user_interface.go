package user_rep

import (
	"context"

	"github.com/google/uuid"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, u structs.User) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (structs.User, error)
	GetByMail(ctx context.Context, mail string) (structs.User, error)
	GetAllUsers(ctx context.Context) ([]structs.User, error)
	GetByPhone(ctx context.Context, phone string) (structs.User, error)
}
