package favourites_rep

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

func (rep *Repository) Create(ctx context.Context, f structs.Favourites) error {
	id, err := rep.GetFIdByUId(ctx, f.Id)
	if err != nil {
		return err
	}
	_, err = rep.db.NamedExecContext(ctx, "insert into favourites (id_user) values ($1)", id)
	return err
}

func (rep *Repository) GetFIdByUId(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var fid uuid.UUID
	err := rep.db.GetContext(ctx, &fid, "select id from favourites where id_user = $1", id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to scan id: %w", err)
	}
	return fid, nil
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Favourites, error) {
	var f rep_structs.Favourites
	err := rep.db.GetContext(ctx, &f, "select * from favourites where id = $1", id)
	if err != nil {
		return structs.Favourites{}, fmt.Errorf("failed to scan favourites: %w", err)
	} else {
		fav := structs.Favourites{
			Id:     f.Id,
			IdUser: f.IdUser,
		}
		return fav, nil
	}
}

func (rep *Repository) GetItems(ctx context.Context, id_user uuid.UUID) ([]structs.FavouritesItem, error) {
	var items []rep_structs.FavouritesItem
	id_favourites, err := rep.GetFIdByUId(ctx, id_user)
	if err != nil {
		return nil, err
	}
	err = rep.db.SelectContext(ctx, &items, "select * from favourites_item where id_favourites = $1", id_favourites)
	if err != nil {
		return nil, err
	} else {
		var itms []structs.FavouritesItem
		for _, v := range items {
			itms = append(itms, structs.FavouritesItem{Id: v.Id,
				IdProduct:    v.IdProduct,
				IdFavourites: v.IdFavourites})
		}
		return itms, err
	}
}

func (rep *Repository) AddItem(ctx context.Context, i structs.FavouritesItem) error {
	it := rep_structs.FavouritesItem{
		IdProduct:    i.IdProduct,
		IdFavourites: i.IdFavourites,
	}

	_, err := rep.db.NamedExecContext(ctx, `insert into favourites_item (id_product, id_favourites) 
	values (:id_product, :id_favourites)`, it)
	return err
}

func (rep *Repository) DeleteItem(ctx context.Context, id uuid.UUID, product_id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from favourites_item where id_product = $1 and id_favourites = $2", product_id, id)
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
