package basket_rep

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
	rep_structs "github.com/taucuya/ppo/internal/repository/postgres/structs"
)

func TestCreate(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		basket      structs.Basket
		expectedErr error
	}{
		{
			name: "successful basket creation",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into basket \\(id_user, date\\) values \\(\\?, \\?\\)").
					WithArgs(fixture.basket.IdUser, fixture.basket.Date).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			basket:      fixture.basket,
			expectedErr: nil,
		},
		{
			name: "basket creation error",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into basket \\(id_user, date\\) values \\(\\?, \\?\\)").
					WithArgs(fixture.basket.IdUser, fixture.basket.Date).
					WillReturnError(errTest)
			},
			basket:      fixture.basket,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Create(fixture.ctx, tt.basket)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetBIdByUId(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		userID      uuid.UUID
		expectedID  uuid.UUID
		expectedErr error
	}{
		{
			name: "successful get basket ID by user ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.basket.Id)
				fixture.mock.ExpectQuery("select id from basket where id_user = \\$1").
					WithArgs(fixture.basket.IdUser).
					WillReturnRows(rows)
			},
			userID:      fixture.basket.IdUser,
			expectedID:  fixture.basket.Id,
			expectedErr: nil,
		},
		{
			name: "basket not found for user",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from basket where id_user = \\$1").
					WithArgs(fixture.basket.IdUser).
					WillReturnError(sql.ErrNoRows)
			},
			userID:      fixture.basket.IdUser,
			expectedID:  uuid.UUID{},
			expectedErr: errors.New("failed to scan id: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error when getting basket ID",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from basket where id_user = \\$1").
					WithArgs(fixture.basket.IdUser).
					WillReturnError(errTest)
			},
			userID:      fixture.basket.IdUser,
			expectedID:  uuid.UUID{},
			expectedErr: errors.New("failed to scan id: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			id, err := fixture.repo.GetBIdByUId(fixture.ctx, tt.userID)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expectedID, id)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetById(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		basketID    uuid.UUID
		expected    structs.Basket
		expectedErr error
	}{
		{
			name: "successful get basket by ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_user", "date"}).
					AddRow(fixture.basket.Id, fixture.basket.IdUser, fixture.basket.Date)
				fixture.mock.ExpectQuery("select \\* from basket where id = \\$1").
					WithArgs(fixture.basket.Id).
					WillReturnRows(rows)
			},
			basketID:    fixture.basket.Id,
			expected:    fixture.basket,
			expectedErr: nil,
		},
		{
			name: "basket not found",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from basket where id = \\$1").
					WithArgs(fixture.basket.Id).
					WillReturnError(sql.ErrNoRows)
			},
			basketID:    fixture.basket.Id,
			expected:    structs.Basket{},
			expectedErr: errors.New("failed to scan basket: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error when getting basket",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from basket where id = \\$1").
					WithArgs(fixture.basket.Id).
					WillReturnError(errTest)
			},
			basketID:    fixture.basket.Id,
			expected:    structs.Basket{},
			expectedErr: errors.New("failed to scan basket: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			basket, err := fixture.repo.GetById(fixture.ctx, tt.basketID)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, basket)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetItems(t *testing.T) {
	fixture := NewTestFixture(t)

	items := []rep_structs.BasketItem{
		{
			Id:        uuid.New(),
			IdProduct: fixture.basketItem.IdProduct,
			IdBasket:  fixture.basket.Id,
			Amount:    1,
		},
		{
			Id:        uuid.New(),
			IdProduct: uuid.New(),
			IdBasket:  fixture.basket.Id,
			Amount:    3,
		},
	}

	expectedItems := []structs.BasketItem{
		{
			Id:        items[0].Id,
			IdProduct: items[0].IdProduct,
			IdBasket:  items[0].IdBasket,
			Amount:    items[0].Amount,
		},
		{
			Id:        items[1].Id,
			IdProduct: items[1].IdProduct,
			IdBasket:  items[1].IdBasket,
			Amount:    items[1].Amount,
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.BasketItem
		expectedErr error
	}{
		{
			name: "successful get items",
			setupMock: func() {
				basketRows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.basket.Id)
				fixture.mock.ExpectQuery("select id from basket where id_user = \\$1").
					WithArgs(fixture.basket.IdUser).
					WillReturnRows(basketRows)

				itemRows := sqlmock.NewRows([]string{"id", "id_product", "id_basket", "amount"}).
					AddRow(items[0].Id, items[0].IdProduct, items[0].IdBasket, items[0].Amount).
					AddRow(items[1].Id, items[1].IdProduct, items[1].IdBasket, items[1].Amount)
				fixture.mock.ExpectQuery("select \\* from basket_item where id_basket = \\$1").
					WithArgs(fixture.basket.Id).
					WillReturnRows(itemRows)
			},
			expected:    expectedItems,
			expectedErr: nil,
		},
		{
			name: "error getting basket ID",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from basket where id_user = \\$1").
					WithArgs(fixture.basket.IdUser).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errors.New("failed to scan id: " + errTest.Error()),
		},
		{
			name: "error getting items",
			setupMock: func() {
				basketRows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.basket.Id)
				fixture.mock.ExpectQuery("select id from basket where id_user = \\$1").
					WithArgs(fixture.basket.IdUser).
					WillReturnRows(basketRows)

				fixture.mock.ExpectQuery("select \\* from basket_item where id_basket = \\$1").
					WithArgs(fixture.basket.Id).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			items, err := fixture.repo.GetItems(fixture.ctx, fixture.basket.IdUser)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, items)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestAddItem(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		item        structs.BasketItem
		expectedErr error
	}{
		{
			name: "successful add new item",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from basket_item where id_basket = \$1 and id_product = \$2`).
					WithArgs(fixture.basketItem.IdBasket, fixture.basketItem.IdProduct).
					WillReturnError(sql.ErrNoRows)

				fixture.mock.ExpectExec(`insert into basket_item`).
					WithArgs(fixture.basketItem.IdProduct, fixture.basketItem.IdBasket, fixture.basketItem.Amount).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "successful update existing item amount",
			setupMock: func() {
				existingItem := rep_structs.BasketItem{
					Id:        uuid.New(),
					IdProduct: fixture.basketItem.IdProduct,
					IdBasket:  fixture.basket.Id,
					Amount:    1,
				}

				rows := sqlmock.NewRows([]string{"id", "id_product", "id_basket", "amount"}).
					AddRow(existingItem.Id, existingItem.IdProduct, existingItem.IdBasket, existingItem.Amount)
				fixture.mock.ExpectQuery(`select \* from basket_item where id_basket = \$1 and id_product = \$2`).
					WithArgs(fixture.basketItem.IdBasket, fixture.basketItem.IdProduct).
					WillReturnRows(rows)

				fixture.mock.ExpectExec("update basket_item set amount = \\$1 where id_product = \\$2 and id_basket = \\$3").
					WithArgs(3, fixture.basketItem.IdProduct, fixture.basketItem.IdBasket).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "error inserting new item",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from basket_item where id_basket = \$1 and id_product = \$2`).
					WithArgs(fixture.basketItem.IdBasket, fixture.basketItem.IdProduct).
					WillReturnError(sql.ErrNoRows)

				fixture.mock.ExpectExec(`insert into basket_item`).
					WithArgs(fixture.basketItem.IdProduct, fixture.basketItem.IdBasket, fixture.basketItem.Amount).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.AddItem(fixture.ctx, fixture.basketItem)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDeleteItem(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		basketID    uuid.UUID
		productID   uuid.UUID
		expectedErr error
	}{
		{
			name: "successful delete item",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from basket_item where id_product = \\$1 and id_basket = \\$2").
					WithArgs(fixture.basketItem.IdProduct, fixture.basket.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			basketID:    fixture.basket.Id,
			productID:   fixture.basketItem.IdProduct,
			expectedErr: nil,
		},
		{
			name: "item not found",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from basket_item where id_product = \\$1 and id_basket = \\$2").
					WithArgs(fixture.basketItem.IdProduct, fixture.basket.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			basketID:    fixture.basket.Id,
			productID:   fixture.basketItem.IdProduct,
			expectedErr: errors.New("item with id " + fixture.basket.Id.String() + " not found"),
		},
		{
			name: "database error when deleting item",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from basket_item where id_product = \\$1 and id_basket = \\$2").
					WithArgs(fixture.basketItem.IdProduct, fixture.basket.Id).
					WillReturnError(errTest)
			},
			basketID:    fixture.basket.Id,
			productID:   fixture.basketItem.IdProduct,
			expectedErr: errTest,
		},
		{
			name: "error getting rows affected",
			setupMock: func() {
				result := sqlmock.NewErrorResult(errTest)
				fixture.mock.ExpectExec("delete from basket_item where id_product = \\$1 and id_basket = \\$2").
					WithArgs(fixture.basketItem.IdProduct, fixture.basket.Id).
					WillReturnResult(result)
			},
			basketID:    fixture.basket.Id,
			productID:   fixture.basketItem.IdProduct,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.DeleteItem(fixture.ctx, tt.basketID, tt.productID)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestUpdateItemAmount(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		basketID    uuid.UUID
		productID   uuid.UUID
		amount      int
		expectedErr error
	}{
		{
			name: "successful update item amount",
			setupMock: func() {
				fixture.mock.ExpectExec("update basket_item set amount = \\$1 where id_product = \\$2 and id_basket = \\$3").
					WithArgs(5, fixture.basketItem.IdProduct, fixture.basket.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			basketID:    fixture.basket.Id,
			productID:   fixture.basketItem.IdProduct,
			amount:      5,
			expectedErr: nil,
		},
		{
			name: "database error when updating item amount",
			setupMock: func() {
				fixture.mock.ExpectExec("update basket_item set amount = \\$1 where id_product = \\$2 and id_basket = \\$3").
					WithArgs(5, fixture.basketItem.IdProduct, fixture.basket.Id).
					WillReturnError(errTest)
			},
			basketID:    fixture.basket.Id,
			productID:   fixture.basketItem.IdProduct,
			amount:      5,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.UpdateItemAmount(fixture.ctx, tt.basketID, tt.productID, tt.amount)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
