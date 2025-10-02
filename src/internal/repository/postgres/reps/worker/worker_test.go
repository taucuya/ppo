package worker_rep

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into worker`).
					WithArgs(testWorker.IdUser, testWorker.JobTitle).
					WillReturnResult(sqlmock.NewResult(1, 1))

				fixture.mock.ExpectExec(`update "user"`).
					WithArgs("работник склада", testWorker.IdUser).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "error inserting worker",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into worker`).
					WithArgs(testWorker.IdUser, testWorker.JobTitle).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "error updating user role",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into worker`).
					WithArgs(testWorker.IdUser, testWorker.JobTitle).
					WillReturnResult(sqlmock.NewResult(1, 1))

				fixture.mock.ExpectExec(`update "user"`).
					WithArgs("работник склада", testWorker.IdUser).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Create(fixture.ctx, testWorker)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetById(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMock   func()
		expected    structs.Worker
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_user", "job_title"}).
					AddRow(testWorker.Id, testWorker.IdUser, testWorker.JobTitle)
				fixture.mock.ExpectQuery(`select \* from worker where id = \$1`).
					WithArgs(testWorker.Id).
					WillReturnRows(rows)
			},
			expected:    testWorker,
			expectedErr: nil,
		},
		{
			name: "worker not found",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from worker where id = \$1`).
					WithArgs(testWorker.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.Worker{},
			expectedErr: errors.New("failed to get worker: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from worker where id = \$1`).
					WithArgs(testWorker.Id).
					WillReturnError(errTest)
			},
			expected:    structs.Worker{},
			expectedErr: errors.New("failed to get worker: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetById(fixture.ctx, testWorker.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetOrders(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testWorker := fixture.workerMother.ValidWorker()
	testOrders := fixture.workerMother.OrdersList()

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.Order
		expectedErr error
	}{
		{
			name: "successful get orders",
			setupMock: func() {
				workerRows := sqlmock.NewRows([]string{"id"}).AddRow(testWorker.Id)
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1`).
					WithArgs(testWorker.IdUser).
					WillReturnRows(workerRows)

				orderRows := sqlmock.NewRows([]string{"id", "date", "id_user", "address", "status", "price"})
				for _, order := range testOrders {
					orderRows.AddRow(order.Id, order.Date, order.IdUser, order.Address, order.Status, order.Price)
				}
				fixture.mock.ExpectQuery(`select \* from "order" where id in`).
					WithArgs(testWorker.Id).
					WillReturnRows(orderRows)
			},
			expected:    testOrders,
			expectedErr: nil,
		},
		{
			name: "worker not found",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1`).
					WithArgs(testWorker.IdUser).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "error getting orders",
			setupMock: func() {
				workerRows := sqlmock.NewRows([]string{"id"}).AddRow(testWorker.Id)
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1`).
					WithArgs(testWorker.IdUser).
					WillReturnRows(workerRows)

				fixture.mock.ExpectQuery(`select \* from "order" where id in`).
					WithArgs(testWorker.Id).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errors.New("failed to get orders: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetOrders(fixture.ctx, testWorker.IdUser)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDelete(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMock: func() {
				fixture.mock.ExpectExec(`delete from worker where id = \$1`).
					WithArgs(testWorker.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "worker not found",
			setupMock: func() {
				fixture.mock.ExpectExec(`delete from worker where id = \$1`).
					WithArgs(testWorker.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: errors.New("worker with id " + testWorker.Id.String() + " not found"),
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectExec(`delete from worker where id = \$1`).
					WithArgs(testWorker.Id).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Delete(fixture.ctx, testWorker.Id)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetAllWorkers(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testWorkers := fixture.workerMother.WorkersList()

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.Worker
		expectedErr error
	}{
		{
			name: "successful get all workers",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_user", "job_title"})
				for _, worker := range testWorkers {
					rows.AddRow(worker.Id, worker.IdUser, worker.JobTitle)
				}
				fixture.mock.ExpectQuery(`select \* from worker order by id`).
					WillReturnRows(rows)
			},
			expected:    testWorkers,
			expectedErr: nil,
		},
		{
			name: "empty list",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_user", "job_title"})
				fixture.mock.ExpectQuery(`select \* from worker order by id`).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from worker order by id`).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetAllWorkers(fixture.ctx)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestAcceptOrder(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testOrder := fixture.workerMother.ValidOrder()
	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful accept order",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into order_worker \(id_order, id_worker\) values \(\$1, \$2\)`).
					WithArgs(testOrder.Id, testWorker.IdUser).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into order_worker \(id_order, id_worker\) values \(\$1, \$2\)`).
					WithArgs(testOrder.Id, testWorker.IdUser).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.AcceptOrder(fixture.ctx, testOrder.Id, testWorker.IdUser)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
