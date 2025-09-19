package basket_rep

import (
	"context"
	"database/sql"
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

func (rep *Repository) GetBIdByUId(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var bid uuid.UUID
	err := rep.db.GetContext(ctx, &bid, "select id from basket where id_user = $1", id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to scan id: %w", err)
	}
	return bid, nil
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Basket, error) {
	var b rep_structs.Basket
	err := rep.db.GetContext(ctx, &b, "select * from basket where id = $1", id)
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

func (rep *Repository) GetItems(ctx context.Context, id_user uuid.UUID) ([]structs.BasketItem, error) {
	var items []rep_structs.BasketItem
	id_basket, err := rep.GetBIdByUId(ctx, id_user)
	if err != nil {
		return nil, err
	}
	err = rep.db.SelectContext(ctx, &items, "select * from basket_item where id_basket = $1", id_basket)
	if err != nil {
		fmt.Println(err)
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
		IdProduct: i.IdProduct,
		IdBasket:  i.IdBasket,
		Amount:    i.Amount,
	}

	var item rep_structs.BasketItem

	err := rep.db.GetContext(ctx, &item, `select * from basket_item where id_basket = $1
	 and id_product = $2`, i.IdBasket, i.IdProduct)
	if err != sql.ErrNoRows {
		err = rep.UpdateItemAmount(ctx, i.IdBasket, i.IdProduct, item.Amount+it.Amount)
	} else {
		_, err = rep.db.NamedExecContext(ctx, `insert into basket_item (id_product, id_basket, amount) 
		values (:id_product, :id_basket, :amount)`, it)
	}
	return err
}

func (rep *Repository) DeleteItem(ctx context.Context, id uuid.UUID, product_id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from basket_item where id_product = $1 and id_basket = $2", product_id, id)
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

func (rep *Repository) UpdateItemAmount(ctx context.Context, basket_id uuid.UUID, product_id uuid.UUID, amount int) error {
	_, err := rep.db.ExecContext(ctx, "update basket_item set amount = $1 where id_product = $2 and id_basket = $3", amount, product_id, basket_id)
	return err
}
