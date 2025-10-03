package basket_rep

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

var (
	errTest = errors.New("test error")
)

type TestFixture struct {
	t          *testing.T
	ctx        context.Context
	db         *sqlx.DB
	mock       sqlmock.Sqlmock
	repo       *Repository
	basket     structs.Basket
	basketItem structs.BasketItem
}

func NewTestFixture(t *testing.T) *TestFixture {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	basketID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()

	return &TestFixture{
		t:    t,
		ctx:  context.Background(),
		db:   sqlxDB,
		mock: mock,
		repo: New(sqlxDB),
		basket: structs.Basket{
			Id:     basketID,
			IdUser: userID,
			Date:   time.Now(),
		},
		basketItem: structs.BasketItem{
			Id:        uuid.New(),
			IdProduct: productID,
			IdBasket:  basketID,
			Amount:    2,
		},
	}
}

func (f *TestFixture) Cleanup() {
	f.db.Close()
}

func (f *TestFixture) AssertError(err error, expectedErr error) {
	if expectedErr != nil {
		if err == nil {
			f.t.Errorf("Expected error %v, got nil", expectedErr)
			return
		} else if err.Error() != expectedErr.Error() {
			f.t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	} else if err != nil {
		f.t.Errorf("Expected error nil, got %v", err)
		return
	}
}
