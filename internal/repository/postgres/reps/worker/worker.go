package worker_rep

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

func (rep *Repository) Create(ctx context.Context, w structs.Worker) error {
	wr := rep_structs.Worker{
		Id:       w.Id,
		IdUser:   w.IdUser,
		JobTitle: w.JobTitle,
	}
	_, err := rep.db.NamedExecContext(ctx,
		"insert into worker (id_user, job_title) values (:id_user, :job_title)",
		wr)
	return err
}

func (rep *Repository) GetById(ctx context.Context, id uuid.UUID) (structs.Worker, error) {
	var w rep_structs.Worker
	err := rep.db.GetContext(ctx, &w, "select * from worker where id = $1", id)
	if err != nil {
		return structs.Worker{}, fmt.Errorf("failed to get worker: %w", err)
	}
	wr := structs.Worker{
		Id:       w.Id,
		IdUser:   w.IdUser,
		JobTitle: w.JobTitle,
	}
	return wr, nil
}

func (rep *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := rep.db.ExecContext(ctx, "delete from worker where id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("worker with id %v not found", id)
	}

	return nil
}

func (rep *Repository) GetAllWorkers(ctx context.Context) ([]structs.Worker, error) {
	var wrs []rep_structs.Worker
	err := rep.db.SelectContext(ctx, &wrs, "select * from worker order by id")
	if err != nil {
		return nil, err
	}

	var w []structs.Worker
	for _, v := range wrs {
		w = append(w, structs.Worker{
			Id:       v.Id,
			IdUser:   v.IdUser,
			JobTitle: v.JobTitle,
		})
	}
	return w, nil
}

func (rep *Repository) AcceptOrder(ctx context.Context, id_order uuid.UUID, id_worker uuid.UUID) error {
	_, err := rep.db.ExecContext(ctx, `update "order" set id_worker = $1 where id = $2`, id_worker, id_order)
	if err != nil {
		return err
	}
	// в курсовой это делает триггер
	_, err = rep.db.ExecContext(ctx, `update "order" set status = $1 where id = $2`, "принятый", id_order)
	return err
}
