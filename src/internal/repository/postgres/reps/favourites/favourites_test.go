package favourites_rep

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
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful favourites creation",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into favourites").
					WithArgs(fixture.favourites.IdUser).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "favourites creation error",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into favourites").
					WithArgs(fixture.favourites.IdUser).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Create(fixture.ctx, fixture.favourites)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetFIdByUId(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expectedID  uuid.UUID
		expectedErr error
	}{
		{
			name: "successful get favourites id by user id",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.favourites.Id)
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnRows(rows)
			},
			expectedID:  fixture.favourites.Id,
			expectedErr: nil,
		},
		{
			name: "favourites not found for user",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnError(sql.ErrNoRows)
			},
			expectedID:  uuid.UUID{},
			expectedErr: errors.New("failed to scan id: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error when getting favourites id",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnError(errTest)
			},
			expectedID:  uuid.UUID{},
			expectedErr: errors.New("failed to scan id: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			id, err := fixture.repo.GetFIdByUId(fixture.ctx, fixture.favourites.IdUser)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expectedID, id)
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
		expected    structs.Favourites
		expectedErr error
	}{
		{
			name: "successful get favourites by id",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "id_user"}).
					AddRow(fixture.favourites.Id, fixture.favourites.IdUser)
				fixture.mock.ExpectQuery("select \\* from favourites where id = \\$1").
					WithArgs(fixture.favourites.Id).
					WillReturnRows(rows)
			},
			expected:    fixture.favourites,
			expectedErr: nil,
		},
		{
			name: "favourites not found",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from favourites where id = \\$1").
					WithArgs(fixture.favourites.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.Favourites{},
			expectedErr: errors.New("failed to scan favourites: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error when getting favourites",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from favourites where id = \\$1").
					WithArgs(fixture.favourites.Id).
					WillReturnError(errTest)
			},
			expected:    structs.Favourites{},
			expectedErr: errors.New("failed to scan favourites: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			favourites, err := fixture.repo.GetById(fixture.ctx, fixture.favourites.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, favourites)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetItems(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	items := []rep_structs.FavouritesItem{
		{
			Id:           uuid.New(),
			IdProduct:    uuid.New(),
			IdFavourites: fixture.favourites.Id,
		},
		{
			Id:           uuid.New(),
			IdProduct:    uuid.New(),
			IdFavourites: fixture.favourites.Id,
		},
	}

	expectedItems := []structs.FavouritesItem{
		{
			Id:           items[0].Id,
			IdProduct:    items[0].IdProduct,
			IdFavourites: items[0].IdFavourites,
		},
		{
			Id:           items[1].Id,
			IdProduct:    items[1].IdProduct,
			IdFavourites: items[1].IdFavourites,
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.FavouritesItem
		expectedErr error
	}{
		{
			name: "successful get items",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.favourites.Id)
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnRows(rows)

				itemRows := sqlmock.NewRows([]string{"id", "id_product", "id_favourites"}).
					AddRow(items[0].Id, items[0].IdProduct, items[0].IdFavourites).
					AddRow(items[1].Id, items[1].IdProduct, items[1].IdFavourites)
				fixture.mock.ExpectQuery("select \\* from favourites_item where id_favourites = \\$1").
					WithArgs(fixture.favourites.Id).
					WillReturnRows(itemRows)
			},
			expected:    expectedItems,
			expectedErr: nil,
		},
		{
			name: "error getting favourites id",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errors.New("failed to scan id: " + errTest.Error()),
		},
		{
			name: "error getting items",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.favourites.Id)
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnRows(rows)

				fixture.mock.ExpectQuery("select \\* from favourites_item where id_favourites = \\$1").
					WithArgs(fixture.favourites.Id).
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errTest,
		},
		{
			name: "no items found",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.favourites.Id)
				fixture.mock.ExpectQuery("select id from favourites where id_user = \\$1").
					WithArgs(fixture.favourites.IdUser).
					WillReturnRows(rows)

				itemRows := sqlmock.NewRows([]string{"id", "id_product", "id_favourites"})
				fixture.mock.ExpectQuery("select \\* from favourites_item where id_favourites = \\$1").
					WithArgs(fixture.favourites.Id).
					WillReturnRows(itemRows)
			},
			expected:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			items, err := fixture.repo.GetItems(fixture.ctx, fixture.favourites.IdUser)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, items)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestAddItem(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful add item",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into favourites_item").
					WithArgs(fixture.favouritesItem.IdProduct, fixture.favouritesItem.IdFavourites).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "error adding item",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into favourites_item").
					WithArgs(fixture.favouritesItem.IdProduct, fixture.favouritesItem.IdFavourites).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.AddItem(fixture.ctx, fixture.favouritesItem)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		expectedErr error
	}{
		{
			name: "successful delete item",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from favourites_item where id_product = \\$1 and id_favourites = \\$2").
					WithArgs(fixture.favouritesItem.IdProduct, fixture.favourites.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "item not found",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from favourites_item where id_product = \\$1 and id_favourites = \\$2").
					WithArgs(fixture.favouritesItem.IdProduct, fixture.favourites.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: errors.New("item with id " + fixture.favourites.Id.String() + " not found"),
		},
		{
			name: "database error when deleting item",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from favourites_item where id_product = \\$1 and id_favourites = \\$2").
					WithArgs(fixture.favouritesItem.IdProduct, fixture.favourites.Id).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "error getting rows affected",
			setupMock: func() {
				result := sqlmock.NewErrorResult(errTest)
				fixture.mock.ExpectExec("delete from favourites_item where id_product = \\$1 and id_favourites = \\$2").
					WithArgs(fixture.favouritesItem.IdProduct, fixture.favourites.Id).
					WillReturnResult(result)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.DeleteItem(fixture.ctx, fixture.favourites.Id, fixture.favouritesItem.IdProduct)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
