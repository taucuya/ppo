package brand_rep

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
		expectedErr error
	}{
		{
			name: "successful brand creation",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into brand \(name, description, price_category\) values \(\?, \?, \?\)`).
					WithArgs(fixture.brand.Name, fixture.brand.Description, fixture.brand.PriceCategory).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "brand creation error",
			setupMock: func() {
				fixture.mock.ExpectExec(`insert into brand \(name, description, price_category\) values \(\?, \?, \?\)`).
					WithArgs(fixture.brand.Name, fixture.brand.Description, fixture.brand.PriceCategory).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Create(fixture.ctx, fixture.brand)

			fixture.AssertError(err, tt.expectedErr)
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
		expected    structs.Brand
		expectedErr error
	}{
		{
			name: "successful get brand by ID",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price_category"}).
					AddRow(fixture.brand.Id, fixture.brand.Name, fixture.brand.Description, fixture.brand.PriceCategory)
				fixture.mock.ExpectQuery("select \\* from brand where id = \\$1").
					WithArgs(fixture.brand.Id).
					WillReturnRows(rows)
			},
			expected:    fixture.brand,
			expectedErr: nil,
		},
		{
			name: "brand not found",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from brand where id = \\$1").
					WithArgs(fixture.brand.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.Brand{},
			expectedErr: errors.New("failed to get brand: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error when getting brand",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from brand where id = \\$1").
					WithArgs(fixture.brand.Id).
					WillReturnError(errTest)
			},
			expected:    structs.Brand{},
			expectedErr: errors.New("failed to get brand: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			brand, err := fixture.repo.GetById(fixture.ctx, fixture.brand.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, brand)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDelete(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		brandID     uuid.UUID
		expectedErr error
	}{
		{
			name: "successful brand deletion",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from brand where id = \\$1").
					WithArgs(fixture.brand.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "brand not found",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from brand where id = \\$1").
					WithArgs(fixture.brand.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: errors.New("brand with id " + fixture.brand.Id.String() + " not found"),
		},
		{
			name: "database error when deleting brand",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from brand where id = \\$1").
					WithArgs(fixture.brand.Id).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.Delete(fixture.ctx, fixture.brand.Id)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetAllBrandsInCategory(t *testing.T) {
	fixture := NewTestFixture(t)

	brands := []rep_structs.Brand{
		{
			Id:            uuid.New(),
			Name:          "Brand 1",
			Description:   "Description 1",
			PriceCategory: "premium",
		},
		{
			Id:            uuid.New(),
			Name:          "Brand 2",
			Description:   "Description 2",
			PriceCategory: "premium",
		},
	}

	expectedBrands := []structs.Brand{
		{
			Id:            brands[0].Id,
			Name:          brands[0].Name,
			Description:   brands[0].Description,
			PriceCategory: brands[0].PriceCategory,
		},
		{
			Id:            brands[1].Id,
			Name:          brands[1].Name,
			Description:   brands[1].Description,
			PriceCategory: brands[1].PriceCategory,
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		category    string
		expected    []structs.Brand
		expectedErr error
	}{
		{
			name: "successful get all brands in category",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price_category"}).
					AddRow(brands[0].Id, brands[0].Name, brands[0].Description, brands[0].PriceCategory).
					AddRow(brands[1].Id, brands[1].Name, brands[1].Description, brands[1].PriceCategory)
				fixture.mock.ExpectQuery("select \\* from brand where price_category = \\$1 order by name").
					WithArgs("premium").
					WillReturnRows(rows)
			},
			category:    "premium",
			expected:    expectedBrands,
			expectedErr: nil,
		},
		{
			name: "no brands found in category",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price_category"})
				fixture.mock.ExpectQuery("select \\* from brand where price_category = \\$1 order by name").
					WithArgs("budget").
					WillReturnRows(rows)
			},
			category:    "budget",
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "database error when getting brands",
			setupMock: func() {
				fixture.mock.ExpectQuery("select \\* from brand where price_category = \\$1 order by name").
					WithArgs("premium").
					WillReturnError(errTest)
			},
			category:    "premium",
			expected:    nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			brands, err := fixture.repo.GetAllBrandsInCategory(fixture.ctx, tt.category)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, brands)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
