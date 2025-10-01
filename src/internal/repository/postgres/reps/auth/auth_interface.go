package auth_rep

import (
	"context"

	"github.com/google/uuid"
)

type AuthRepositoryInterface interface {
	CreateToken(ctx context.Context, id uuid.UUID, rtoken string) error
	CheckAdmin(ctx context.Context, id uuid.UUID) bool
	CheckWorker(ctx context.Context, id uuid.UUID) bool
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteToken(ctx context.Context, id uuid.UUID) error
}
