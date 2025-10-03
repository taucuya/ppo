package brand_rep

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

func (rep *Repository) Create(ctx context.Context, b structs.Brand) error {
	br := rep_structs.Brand{
		Id:            b.Id,
		Name:          b.Name,
		Description:   b.Description,
		PriceCategory: b.PriceCategory,
	}
	_, err := rep.db.NamedExecContext(ctx,
		"insert into brand (name, description, price_category) values (:name, :description, :price_category)",
		br)
	return err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Brand, error) {
	var b rep_structs.Brand
	err := rep.db.GetContext(ctx, &b, "select * from brand where id = $1", id)
	if err != nil {
		return structs.Brand{}, fmt.Errorf("failed to get brand: %w", err)
	}
	br := structs.Brand{
		Id:            b.Id,
		Name:          b.Name,
		Description:   b.Description,
		PriceCategory: b.PriceCategory,
	}
	return br, nil
}

func (rep *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from brand where id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("brand with id %v not found", id)
	}

	return nil
}

func (rep *Repository) GetAllBrandsInCategory(ctx context.Context, category string) ([]structs.Brand, error) {
	var brands []rep_structs.Brand
	err := rep.db.SelectContext(ctx, &brands, "select * from brand where price_category = $1 order by name", category)
	if err != nil {
		return nil, err
	}
	var br []structs.Brand
	for _, v := range brands {
		br = append(br, structs.Brand{
			Id:            v.Id,
			Name:          v.Name,
			Description:   v.Description,
			PriceCategory: v.PriceCategory,
		})
	}
	return br, nil
}
