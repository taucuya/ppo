package structs

import "github.com/google/uuid"

type Worker struct {
	Id       uuid.UUID `db:"id"`
	IdUser   uuid.UUID `db:"id_user"`
	JobTitle string    `db:"job_title"`
}

type WorkersOrders struct {
	IdOrder  uuid.UUID `db:"id_order"`
	IdWorker uuid.UUID `db:"id_worker"`
}
