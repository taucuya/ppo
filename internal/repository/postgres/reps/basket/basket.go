package basket_rep

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

func (rep *Repository) Create(ctx context.Context, b structs.Basket) error {
	ba := rep_structs.Basket{
		Id:     b.Id,
		IdUser: b.IdUser,
		Date:   b.Date,
	}
	_, err := rep.db.NamedExecContext(ctx, "insert into basket (id_user, date) values (:id_user, :date)", ba)
	return err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Basket, error) {
	var b structs.Basket
	err := rep.db.GetContext(ctx, &b, "select * from baket where id = $1", id)
	if err != nil {
		return structs.Basket{}, fmt.Errorf("failed to scan basket: %w", err)
	} else {
		ba := structs.Basket{
			Id:     b.Id,
			IdUser: b.IdUser,
			Date:   b.Date,
		}
		return ba, nil
	}
}

func (rep *Repository) GetItems(ctx context.Context, id_basket uuid.UUID) ([]structs.BasketItem, error) {
	var items []rep_structs.BasketItem
	err := rep.db.SelectContext(ctx, &items, "select * from basket_item where id_basket = $1", id_basket)
	if err != nil {
		return nil, err
	} else {
		var itms []structs.BasketItem
		for _, v := range items {
			itms = append(itms, structs.BasketItem{Id: v.Id,
				IdProduct: v.IdProduct,
				IdBasket:  v.IdBasket,
				Amount:    v.Amount})
		}
		return itms, err
	}
}

func (rep *Repository) AddItem(ctx context.Context, i structs.BasketItem) error {
	it := rep_structs.BasketItem{
		Id:        i.Id,
		IdProduct: i.IdProduct,
		IdBasket:  i.IdBasket,
		Amount:    i.Amount,
	}
	_, err := rep.db.NamedExecContext(ctx, `insert into basket_item (id_product, id_basket, amount) 
		values (:id_product, :id_basket, :amount)`, it)
	return err
}

func (rep *Repository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from basket_item where id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with id %v not found", id)
	}

	return nil
}

func (rep *Repository) UpdateItemAmount(ctx context.Context, id uuid.UUID, amount int) error {
	_, err := rep.db.ExecContext(ctx, "update basket_item set amount = $1 where id = $2", amount, id)
	return err
}
