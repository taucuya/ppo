package order_rep

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

func (rep *Repository) Create(ctx context.Context, o structs.Order) error {
	var id uuid.UUID
	ord := rep_structs.Order{
		// Date:    o.Date,
		IdUser:  o.IdUser,
		Address: o.Address,
		Status:  o.Status,
		// Price:   o.Price,
	}
	err := rep.db.QueryRowContext(ctx, `
		insert into "order" (id_user, address, status) 
		values ($1, $2, $3) 
		returning id`,
		ord.IdUser, ord.Address, ord.Status).Scan(&id)
	// if err != nil {
	// 	return err
	// }
	// var basket_id uuid.UUID
	// var items []rep_structs.BasketItem
	// err = rep.db.GetContext(ctx, &basket_id, "select id from basket where id_user = $1", o.IdUser)
	// if err != nil {
	// 	return err
	// }
	// err = rep.db.SelectContext(ctx, &items, "select * from basket_item where id_basket = $1", basket_id)
	// if err != nil {
	// 	return err
	// }
	// for _, v := range items {
	// 	var am int
	// 	if err := rep.db.GetContext(ctx, &am, `select amount from product where id = $1`, v.IdProduct); err != nil {
	// 		return err
	// 	}
	// 	if am < v.Amount {
	// 		return fmt.Errorf(`not enouth products on warehouse`)
	// 	}

	// 	_, err := rep.db.ExecContext(ctx, `update product set amount = $1 where id = $2`, am - v.Amount,v.IdProduct)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// var sm float64
	// err = rep.db.GetContext(ctx, &sm, `select sum(p.price * b.amount) as total_price from basket_item b
	// 		join product p on b.id_product = p.id where b.id_basket = $1;`, basket_id)

	// if err != nil {
	// 	return err
	// }

	// _, err = rep.db.ExecContext(ctx, `update "order" set price = $1 where id = $2`, sm, id)
	// if err != nil {
	// 	return err
	// }

	// for _, item := range items {
	// 	_, err := rep.db.ExecContext(ctx,
	// 		"insert into order_item (id_product, id_order, amount) values ($1, $2, $3)",
	// 		item.IdProduct, id, item.Amount)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Order, error) {
	var o rep_structs.Order
	err := rep.db.GetContext(ctx, &o, "select * from \"order\" where id = $1", id)
	if err != nil {
		return structs.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	ord := structs.Order{
		Date:    o.Date,
		IdUser:  o.IdUser,
		Address: o.Address,
		Status:  o.Status,
		Price:   o.Price,
	}
	return ord, nil
}

func (rep *Repository) GetItems(ctx context.Context, id uuid.UUID) ([]structs.OrderItem, error) {
	var items []rep_structs.OrderItem
	err := rep.db.SelectContext(ctx, &items, "select * from order_item where id_order = $1", id)
	if err != nil {
		return nil, err
	}
	var itms []structs.OrderItem
	for _, v := range items {
		itms = append(itms, structs.OrderItem{
			Id:        v.Id,
			IdProduct: v.IdProduct,
			IdOrder:   v.IdOrder,
			Amount:    v.Amount,
		})
	}
	return itms, nil
}

func (rep *Repository) GetFreeOrders(ctx context.Context) ([]structs.Order, error) {
	var orders []rep_structs.Order
	err := rep.db.SelectContext(ctx, &orders, `select * from "order" where status = $1`, "непринятый")
	if err != nil {
		return nil, err
	}
	var ords []structs.Order
	for _, v := range orders {
		ords = append(ords, structs.Order{
			Id:      v.Id,
			Date:    v.Date,
			IdUser:  v.IdUser,
			Address: v.Address,
			Status:  v.Status,
			Price:   v.Price,
		})
	}
	return ords, nil
}

func (rep *Repository) GetStatus(ctx context.Context, id uuid.UUID) (string, error) {
	var status string
	err := rep.db.GetContext(ctx, &status, "select status from \"order\" where id = $1", id)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (rep *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from \"order\" where id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order with id %v not found", id)
	}

	return nil
}

func (rep *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	result, err := rep.db.ExecContext(ctx, "update \"order\" set status = $1 where id = $2", status, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order with id %v not found", id)
	}

	return nil
}
