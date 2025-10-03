package order_rep

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
	rep_structs "github.com/taucuya/ppo/internal/repository/postgres/structs"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful order creation",
			setupMock: func() {
				fixture.mock.ExpectQuery(`insert into "order" \(id_user, address, status\) values \(\$1, \$2, \$3\) returning id`).
					WithArgs(fixture.order.IdUser, fixture.order.Address, fixture.order.Status).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fixture.order.Id))
			},
			expectedErr: nil,
		},
		{
			name: "order creation error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`insert into "order" \(id_user, address, status\) values \(\$1, \$2, \$3\) returning id`).
					WithArgs(fixture.order.IdUser, fixture.order.Address, fixture.order.Status).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Create(fixture.ctx, fixture.order)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetById(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expected    structs.Order
		expectedErr error
	}{
		{
			name: "successful get order by id",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"date", "id_user", "address", "status", "price"}).
					AddRow(fixture.order.Date, fixture.order.IdUser, fixture.order.Address, fixture.order.Status, fixture.order.Price)
				fixture.mock.ExpectQuery(`select \* from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnRows(rows)
			},
			expected:    fixture.order,
			expectedErr: nil,
		},
		{
			name: "order not found",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.Order{},
			expectedErr: errors.New("failed to get order: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error when getting order",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnError(errTest)
			},
			expected:    structs.Order{},
			expectedErr: errors.New("failed to get order: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			order, err := fixture.repo.GetById(fixture.ctx, fixture.order.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, order)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetItems(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	items := []rep_structs.OrderItem{
		{
			Id:        uuid.New(),
			IdProduct: uuid.New(),
			IdOrder:   fixture.order.Id,
			Amount:    2,
		},
		{
			Id:        uuid.New(),
			IdProduct: uuid.New(),
			IdOrder:   fixture.order.Id,
			Amount:    1,
		},
	}

	expectedItems := []structs.OrderItem{
		{
			Id:        items[0].Id,
			IdProduct: items[0].IdProduct,
			IdOrder:   items[0].IdOrder,
			Amount:    items[0].Amount,
		},
		{
			Id:        items[1].Id,
			IdProduct: items[1].IdProduct,
			IdOrder:   items[1].IdOrder,
			Amount:    items[1].Amount,
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.OrderItem
		expectedErr error
	}{
		{
			name: "successful get items",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_product", "id_order", "amount"}).
					AddRow(items[0].Id, items[0].IdProduct, items[0].IdOrder, items[0].Amount).
					AddRow(items[1].Id, items[1].IdProduct, items[1].IdOrder, items[1].Amount)
				fixture.mock.ExpectQuery("select \\* from order_item where id_order = \\$1").
					WithArgs(fixture.order.Id).
					WillReturnRows(rows)
			},
			expected:    expectedItems,
			expectedErr: nil,
		},
		{
			name: "no items found",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_product", "id_order", "amount"})
				fixture.mock.ExpectQuery("select \\* from order_item where id_order = \\$1").
					WithArgs(fixture.order.Id).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "database error when getting items",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from order_item where id_order = \\$1").
					WithArgs(fixture.order.Id).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			items, err := fixture.repo.GetItems(fixture.ctx, fixture.order.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, items)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetFreeOrders(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	orders := []rep_structs.Order{
		{
			Id:      uuid.New(),
			Date:    time.Now(),
			IdUser:  uuid.New(),
			Address: "Address 1",
			Status:  "непринятый",
			Price:   150.0,
		},
		{
			Id:      uuid.New(),
			Date:    time.Now().Add(-time.Hour),
			IdUser:  uuid.New(),
			Address: "Address 2",
			Status:  "непринятый",
			Price:   200.0,
		},
	}

	expectedOrders := []structs.Order{
		{
			Id:      orders[0].Id,
			Date:    orders[0].Date,
			IdUser:  orders[0].IdUser,
			Address: orders[0].Address,
			Status:  orders[0].Status,
			Price:   orders[0].Price,
		},
		{
			Id:      orders[1].Id,
			Date:    orders[1].Date,
			IdUser:  orders[1].IdUser,
			Address: orders[1].Address,
			Status:  orders[1].Status,
			Price:   orders[1].Price,
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.Order
		expectedErr error
	}{
		{
			name: "successful get free orders",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "id_user", "address", "status", "price"}).
					AddRow(orders[0].Id, orders[0].Date, orders[0].IdUser, orders[0].Address, orders[0].Status, orders[0].Price).
					AddRow(orders[1].Id, orders[1].Date, orders[1].IdUser, orders[1].Address, orders[1].Status, orders[1].Price)
				fixture.mock.ExpectQuery(`select \* from "order" where status = \$1`).
					WithArgs("непринятый").
					WillReturnRows(rows)
			},
			expected:    expectedOrders,
			expectedErr: nil,
		},
		{
			name: "no free orders found",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "id_user", "address", "status", "price"})
				fixture.mock.ExpectQuery(`select \* from "order" where status = \$1`).
					WithArgs("непринятый").
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "database error when getting free orders",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "order" where status = \$1`).
					WithArgs("непринятый").
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			orders, err := fixture.repo.GetFreeOrders(fixture.ctx)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, orders)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetOrdersByUser(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	orders := []rep_structs.Order{
		{
			Id:      uuid.New(),
			Date:    time.Now(),
			IdUser:  fixture.order.IdUser,
			Address: "User Address 1",
			Status:  "принятый",
			Price:   100.0,
		},
	}

	expectedOrders := []structs.Order{
		{
			Id:      orders[0].Id,
			Date:    orders[0].Date,
			IdUser:  orders[0].IdUser,
			Address: orders[0].Address,
			Status:  orders[0].Status,
			Price:   orders[0].Price,
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.Order
		expectedErr error
	}{
		{
			name: "successful get orders by user",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "id_user", "address", "status", "price"}).
					AddRow(orders[0].Id, orders[0].Date, orders[0].IdUser, orders[0].Address, orders[0].Status, orders[0].Price)
				fixture.mock.ExpectQuery(`select \* from "orders" where id_user = \$1`).
					WithArgs(fixture.order.IdUser).
					WillReturnRows(rows)
			},
			expected:    expectedOrders,
			expectedErr: nil,
		},
		{
			name: "no orders found for user",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "date", "id_user", "address", "status", "price"})
				fixture.mock.ExpectQuery(`select \* from "orders" where id_user = \$1`).
					WithArgs(fixture.order.IdUser).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "database error when getting user orders",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "orders" where id_user = \$1`).
					WithArgs(fixture.order.IdUser).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			orders, err := fixture.repo.GetOrdersByUser(fixture.ctx, fixture.order.IdUser)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, orders)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetStatus(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expected    string
		expectedErr error
	}{
		{
			name: "successful get status",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"status"}).AddRow(fixture.order.Status)
				fixture.mock.ExpectQuery(`select status from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnRows(rows)
			},
			expected:    fixture.order.Status,
			expectedErr: nil,
		},
		{
			name: "order not found",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select status from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    "",
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "database error when getting status",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select status from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnError(errTest)
			},
			expected:    "",
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			status, err := fixture.repo.GetStatus(fixture.ctx, fixture.order.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, status)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDelete(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful delete order",
			setupMock: func() {
				fixture.mock.ExpectExec(`delete from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "order not found",
			setupMock: func() {
				fixture.mock.ExpectExec(`delete from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: errors.New("order with id " + fixture.order.Id.String() + " not found"),
		},
		{
			name: "database error when deleting order",
			setupMock: func() {
				fixture.mock.ExpectExec(`delete from "order" where id = \$1`).
					WithArgs(fixture.order.Id).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Delete(fixture.ctx, fixture.order.Id)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		orderID     uuid.UUID
		status      string
		expectedErr error
	}{
		{
			name: "successful update status",
			setupMock: func() {
				fixture.mock.ExpectExec(`update "order" set status = \$1 where id = \$2`).
					WithArgs("принятый", fixture.order.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			status:      "принятый",
			expectedErr: nil,
		},
		{
			name: "order not found",
			setupMock: func() {
				fixture.mock.ExpectExec(`update "order" set status = \$1 where id = \$2`).
					WithArgs("в обработке", fixture.order.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			status:      "в обработке",
			expectedErr: errors.New("order with id " + fixture.order.Id.String() + " not found"),
		},
		{
			name: "database error when updating status",
			setupMock: func() {
				fixture.mock.ExpectExec(`update "order" set status = \$1 where id = \$2`).
					WithArgs("на сборке", fixture.order.Id).
					WillReturnError(errTest)
			},
			status:      "на сборке",
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.UpdateStatus(fixture.ctx, fixture.order.Id, tt.status)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
