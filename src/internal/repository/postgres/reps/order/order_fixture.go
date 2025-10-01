package order_rep

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type TestFixture struct {
	t      *testing.T
	db     *sql.DB
	sqlxDB *sqlx.DB
	mock   sqlmock.Sqlmock
	repo   *Repository
	ctx    context.Context
	order  structs.Order
}

func NewTestFixture(t *testing.T) *TestFixture {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	order := structs.Order{
		Date:    time.Now(),
		IdUser:  uuid.New(),
		Address: "Test Address",
		Status:  "непринятый",
		Price:   100.0,
	}

	return &TestFixture{
		t:      t,
		db:     db,
		sqlxDB: sqlxDB,
		mock:   mock,
		repo:   New(sqlxDB),
		ctx:    context.Background(),
		order:  order,
	}
}

func (f *TestFixture) AssertError(actual, expected error) {
	if expected == nil {
		assert.NoError(f.t, actual)
	} else {
		assert.EqualError(f.t, actual, expected.Error())
	}
}

func (f *TestFixture) Cleanup() {
	f.db.Close()
}
