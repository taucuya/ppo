package review_rep

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

func (rep *Repository) Create(ctx context.Context, r structs.Review) error {
	rw := rep_structs.Review{
		IdProduct: r.IdProduct,
		IdUser:    r.IdUser,
		Rating:    r.Rating,
		Text:      r.Text,
		Date:      r.Date,
	}
	_, err := rep.db.NamedExecContext(ctx,
		`insert into review (id_product, id_user, rating, r_text, date)
		 values (:id_product, :id_user, :rating, :r_text, :date)`,
		rw)

	return err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Review, error) {
	var r rep_structs.Review
	err := rep.db.GetContext(ctx, &r, "select * from review where id = $1", id)
	if err != nil {
		return structs.Review{}, fmt.Errorf("failed to get review: %w", err)
	}
	rw := structs.Review{
		IdProduct: r.IdProduct,
		IdUser:    r.IdUser,
		Rating:    r.Rating,
		Text:      r.Text,
		Date:      r.Date,
	}
	return rw, nil
}

func (rep *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from review where id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("review with id %v not found", id)
	}

	return nil
}

func (rep *Repository) ReviewsForProduct(ctx context.Context, id_product uuid.UUID) ([]structs.Review, error) {
	var rws []rep_structs.Review
	err := rep.db.SelectContext(ctx, &rws,
		"select * from review where id_product = $1 order by date desc",
		id_product)
	if err != nil {
		return nil, err
	}
	var r []structs.Review
	for _, v := range rws {
		r = append(r, structs.Review{
			Id:        v.Id,
			IdProduct: v.IdProduct,
			IdUser:    v.IdUser,
			Rating:    v.Rating,
			Text:      v.Text,
			Date:      v.Date,
		})
	}
	return r, nil
}
