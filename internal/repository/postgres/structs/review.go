package structs

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	Id        uuid.UUID `db:"id"`
	IdProduct uuid.UUID `db:"id_product"`
	IdUser    uuid.UUID `db:"id_user"`
	Rating    int       `db:"rating"`
	Text      string    `db:"r_text"`
	Date      time.Time `db:"date"`
}
