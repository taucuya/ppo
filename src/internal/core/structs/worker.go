package structs

import (
	"errors"

	"github.com/google/uuid"
)

type Worker struct {
	Id       uuid.UUID
	IdUser   uuid.UUID
	JobTitle string
}

type WorkersOrders struct {
	IdOrder  uuid.UUID
	IdWorker uuid.UUID
}

var (
	ErrWorkerNotFound       = errors.New("worker not found")
	ErrDuplicateWorker      = errors.New("duplicate worker")
	ErrOrderAlreadyAccepted = errors.New("order already accepted")
)
