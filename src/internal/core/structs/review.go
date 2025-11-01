package structs

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Review struct {
	Id        uuid.UUID `pg:"id"`
	IdProduct uuid.UUID `pg:"id_product"`
	IdUser    uuid.UUID `pg:"id_user"`
	Rating    int       `pg:"rating"`
	Text      string    `pg:"r_text"`
	Date      time.Time `pg:"date"`
}

var (
	ErrDuplicateReview = errors.New("duplicate review")
)
