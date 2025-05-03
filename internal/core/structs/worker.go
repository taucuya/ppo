package structs

import "github.com/google/uuid"

type Worker struct {
	Id       uuid.UUID
	IdUser   uuid.UUID
	JobTitle string
}
