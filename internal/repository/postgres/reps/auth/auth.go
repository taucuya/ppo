package auth_rep

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (rep *Repository) CreateToken(ctx context.Context, id uuid.UUID, rtoken string) error {
	_, err := rep.db.ExecContext(ctx, "insert into token (rtoken) values ($1)", rtoken)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (rep *Repository) VerifyToken(ctx context.Context, token string) (uuid.UUID, error) {
	var id uuid.UUID
	err := rep.db.GetContext(ctx, &id, "select id from token where rtoken = $1", token)
	if err != nil {
		return uuid.UUID{}, err
	} else {
		return id, nil
	}
}

func (rep *Repository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from token where id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token with id %v not found", id)
	}

	return nil
}
