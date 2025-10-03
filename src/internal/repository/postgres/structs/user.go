package structs

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Date_of_birth time.Time `db:"date_of_birth"`
	Mail          string    `db:"mail"`
	Password      string    `db:"password"`
	Phone         string    `db:"phone"`
	Address       string    `db:"address"`
	Status        string    `db:"status"`
	Role          string    `db:"role"`
}
