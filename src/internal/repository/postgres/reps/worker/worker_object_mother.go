package worker_rep

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

// WorkerMother для создания тестовых данных по паттерну Object Mother
type WorkerMother struct{}

func NewWorkerMother() *WorkerMother {
	return &WorkerMother{}
}

func (m *WorkerMother) ValidWorker() structs.Worker {
	return structs.Worker{
		Id:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		IdUser:   uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		JobTitle: "Developer",
	}
}

func (m *WorkerMother) AnotherWorker() structs.Worker {
	return structs.Worker{
		Id:       uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		IdUser:   uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		JobTitle: "Manager",
	}
}

func (m *WorkerMother) WorkersList() []structs.Worker {
	return []structs.Worker{
		m.ValidWorker(),
		m.AnotherWorker(),
	}
}

func (m *WorkerMother) ValidOrder() structs.Order {
	return structs.Order{
		Id:      uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		Date:    time.Now(),
		IdUser:  uuid.MustParse("66666666-6666-6666-6666-666666666666"),
		Address: "Test Address",
		Status:  "pending",
		Price:   100.0,
	}
}

func (m *WorkerMother) AnotherOrder() structs.Order {
	return structs.Order{
		Id:      uuid.MustParse("77777777-7777-7777-7777-777777777777"),
		Date:    time.Now().Add(-24 * time.Hour),
		IdUser:  uuid.MustParse("88888888-8888-8888-8888-888888888888"),
		Address: "Another Address",
		Status:  "in_progress",
		Price:   200.0,
	}
}

func (m *WorkerMother) OrdersList() []structs.Order {
	return []structs.Order{
		m.ValidOrder(),
		m.AnotherOrder(),
	}
}

// TestFixture для управления зависимостями тестов
type TestFixture struct {
	t            *testing.T
	db           *sql.DB
	sqlxDB       *sqlx.DB
	mock         sqlmock.Sqlmock
	repo         *Repository
	ctx          context.Context
	workerMother *WorkerMother
}

func NewTestFixture(t *testing.T) *TestFixture {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return &TestFixture{
		t:            t,
		db:           db,
		sqlxDB:       sqlxDB,
		mock:         mock,
		repo:         New(sqlxDB),
		ctx:          context.Background(),
		workerMother: NewWorkerMother(),
	}
}

func (f *TestFixture) Cleanup() {
	f.db.Close()
}

func (f *TestFixture) AssertError(actual, expected error) {
	if expected == nil {
		assert.NoError(f.t, actual)
	} else {
		assert.EqualError(f.t, actual, expected.Error())
	}
}
