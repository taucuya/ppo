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
		Id:          p.Id,
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

func (rep *Repository) GetByName(ctx context.Context, name string) (structs.Product, error) {
	var p rep_structs.Product
	err := rep.db.GetContext(ctx, &p,
		`select * from product where name = $1`, name)
	if err != nil {
		return structs.Product{}, fmt.Errorf("failed to get product: %w", err)
	}
	pr := structs.Product{
		Id:          p.Id,
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

func (rep *Repository) GetByArticule(ctx context.Context, art string) (structs.Product, error) {
	var p rep_structs.Product
	err := rep.db.GetContext(ctx, &p,
		`select * from product where art = $1`, art)
	if err != nil {
		return structs.Product{}, fmt.Errorf("failed to get product: %w", err)
	}
	pr := structs.Product{
		Id:          p.Id,
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

func (rep *Repository) GetByCategory(ctx context.Context, category string) ([]structs.Product, error) {
	var ps []rep_structs.Product
	if err := rep.db.SelectContext(ctx, &ps, `select * from product where category = $1`, category); err != nil {
		return nil, err
	}

	var products []structs.Product

	for _, v := range ps {
		products = append(products, structs.Product{
			Id:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Category:    v.Category,
			Amount:      v.Amount,
			IdBrand:     v.IdBrand,
			PicLink:     v.PicLink,
			Articule:    v.Articule,
		})
	}
	return products, nil
}

func (rep *Repository) GetByBrand(ctx context.Context, brand string) ([]structs.Product, error) {
	var ps []rep_structs.Product
	if err := rep.db.SelectContext(ctx, &ps, `select * from product where id_brand in
	 (select id from brand where name = $1)`, brand); err != nil {
		return nil, err
	}

	var products []structs.Product

	for _, v := range ps {
		products = append(products, structs.Product{
			Id:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Category:    v.Category,
			Amount:      v.Amount,
			IdBrand:     v.IdBrand,
			PicLink:     v.PicLink,
			Articule:    v.Articule,
		})
	}
	return products, nil
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
