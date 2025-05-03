package product_rep

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	structs "github.com/taucuya/ppo/internal/core/structs"
	rep_structs "github.com/taucuya/ppo/internal/repository/postgres/structs"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (rep *Repository) Create(ctx context.Context, p structs.Product) error {
	pr := rep_structs.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Amount:      p.Amount,
		IdBrand:     p.IdBrand,
		PicLink:     p.PicLink,
		Articule:    p.Articule,
	}
	_, err := rep.db.NamedExecContext(ctx,
		`insert into product 
		(name, description, price, category, amount, id_brand, pic_link, art) 
		values 
		(:name, :description, :price, :category, :amount, :id_brand, :pic_link, :art)`,
		pr)
	return err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Product, error) {
	var p rep_structs.Product
	err := rep.db.GetContext(ctx, &p,
		`select * from product where id = $1`, id)
	if err != nil {
		return structs.Product{}, fmt.Errorf("failed to get product: %w", err)
	}
	pr := structs.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Amount:      p.Amount,
		IdBrand:     p.IdBrand,
		PicLink:     p.PicLink,
		Articule:    p.Articule,
	}
	return pr, nil
}

func (rep *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx,
		`delete from product where id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product with id %v not found", id)
	}

	return nil
}
