package product_rep

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(structs.Product)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(product structs.Product) {
				fixture.mock.ExpectExec(`insert into product`).
					WithArgs(
						product.Name,
						product.Description,
						product.Price,
						product.Category,
						product.Amount,
						product.IdBrand,
						product.PicLink,
						product.Articule,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(product structs.Product) {
				fixture.mock.ExpectExec(`insert into product`).
					WithArgs(
						product.Name,
						product.Description,
						product.Price,
						product.Category,
						product.Amount,
						product.IdBrand,
						product.PicLink,
						product.Articule,
					).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(testProduct)

			err := fixture.repo.Create(fixture.ctx, testProduct)
			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetById(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(structs.Product)
		expectedRet structs.Product
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMocks: func(product structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"}).
					AddRow(product.Id, product.Name, product.Description, product.Price, product.Category, product.Amount, product.IdBrand, product.PicLink, product.Articule)
				fixture.mock.ExpectQuery(`select \* from product where id = \$1`).
					WithArgs(product.Id).
					WillReturnRows(rows)
			},
			expectedRet: testProduct,
			expectedErr: nil,
		},
		{
			name: "product not found",
			setupMocks: func(product structs.Product) {
				fixture.mock.ExpectQuery(`select \* from product where id = \$1`).
					WithArgs(product.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRet: structs.Product{},
			expectedErr: errors.New("failed to get product: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error",
			setupMocks: func(product structs.Product) {
				fixture.mock.ExpectQuery(`select \* from product where id = \$1`).
					WithArgs(product.Id).
					WillReturnError(errTest)
			},
			expectedRet: structs.Product{},
			expectedErr: errors.New("failed to get product: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(testProduct)

			ret, err := fixture.repo.GetById(fixture.ctx, testProduct.Id)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetByName(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.WithName("UniqueProduct").Build()

	tests := []struct {
		name        string
		setupMocks  func(structs.Product)
		expectedRet structs.Product
		expectedErr error
	}{
		{
			name: "successful get by name",
			setupMocks: func(product structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"}).
					AddRow(product.Id, product.Name, product.Description, product.Price, product.Category, product.Amount, product.IdBrand, product.PicLink, product.Articule)
				fixture.mock.ExpectQuery(`select \* from product where name = \$1`).
					WithArgs(product.Name).
					WillReturnRows(rows)
			},
			expectedRet: testProduct,
			expectedErr: nil,
		},
		{
			name: "product not found by name",
			setupMocks: func(product structs.Product) {
				fixture.mock.ExpectQuery(`select \* from product where name = \$1`).
					WithArgs(product.Name).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRet: structs.Product{},
			expectedErr: errors.New("failed to get product: " + sql.ErrNoRows.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(testProduct)

			ret, err := fixture.repo.GetByName(fixture.ctx, testProduct.Name)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetByArticule(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.WithArticule("UNIQUE123").Build()

	tests := []struct {
		name        string
		setupMocks  func(structs.Product)
		expectedRet structs.Product
		expectedErr error
	}{
		{
			name: "successful get by articule",
			setupMocks: func(product structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"}).
					AddRow(product.Id, product.Name, product.Description, product.Price, product.Category, product.Amount, product.IdBrand, product.PicLink, product.Articule)
				fixture.mock.ExpectQuery(`select \* from product where art = \$1`).
					WithArgs(product.Articule).
					WillReturnRows(rows)
			},
			expectedRet: testProduct,
			expectedErr: nil,
		},
		{
			name: "product not found by articule",
			setupMocks: func(product structs.Product) {
				fixture.mock.ExpectQuery(`select \* from product where art = \$1`).
					WithArgs(product.Articule).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRet: structs.Product{},
			expectedErr: errors.New("failed to get product: " + sql.ErrNoRows.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(testProduct)

			ret, err := fixture.repo.GetByArticule(fixture.ctx, testProduct.Articule)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetByCategory(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	category := "electronics"
	testProducts := []structs.Product{
		fixture.productBuilder.WithCategory(category).Build(),
		fixture.productBuilder.WithCategory(category).WithName("Product 2").Build(),
	}

	tests := []struct {
		name        string
		category    string
		setupMocks  func(string, []structs.Product)
		expectedRet []structs.Product
		expectedErr error
	}{
		{
			name:     "successful get by category",
			category: category,
			setupMocks: func(category string, products []structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"})
				for _, p := range products {
					rows.AddRow(p.Id, p.Name, p.Description, p.Price, p.Category, p.Amount, p.IdBrand, p.PicLink, p.Articule)
				}
				fixture.mock.ExpectQuery(`select \* from product where category = \$1`).
					WithArgs(category).
					WillReturnRows(rows)
			},
			expectedRet: testProducts,
			expectedErr: nil,
		},
		{
			name:     "empty category",
			category: "nonexistent",
			setupMocks: func(category string, products []structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"})
				fixture.mock.ExpectQuery(`select \* from product where category = \$1`).
					WithArgs(category).
					WillReturnRows(rows)
			},
			expectedRet: nil,
			expectedErr: nil,
		},
		{
			name:     "database error",
			category: category,
			setupMocks: func(category string, products []structs.Product) {
				fixture.mock.ExpectQuery(`select \* from product where category = \$1`).
					WithArgs(category).
					WillReturnError(errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(tt.category, testProducts)

			ret, err := fixture.repo.GetByCategory(fixture.ctx, tt.category)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetByBrand(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	brand := "TestBrand"
	testProducts := []structs.Product{
		fixture.productBuilder.Build(),
		fixture.productBuilder.WithName("Product 2").Build(),
	}

	tests := []struct {
		name        string
		brand       string
		setupMocks  func(string, []structs.Product)
		expectedRet []structs.Product
		expectedErr error
	}{
		{
			name:  "successful get by brand",
			brand: brand,
			setupMocks: func(brand string, products []structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"})
				for _, p := range products {
					rows.AddRow(p.Id, p.Name, p.Description, p.Price, p.Category, p.Amount, p.IdBrand, p.PicLink, p.Articule)
				}
				fixture.mock.ExpectQuery(`select \* from product where id_brand in`).
					WithArgs(brand).
					WillReturnRows(rows)
			},
			expectedRet: testProducts,
			expectedErr: nil,
		},
		{
			name:  "empty brand",
			brand: "SomeBrand",
			setupMocks: func(brand string, products []structs.Product) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category", "amount", "id_brand", "pic_link", "art"})
				fixture.mock.ExpectQuery(`select \* from product where id_brand in`).
					WithArgs(brand).
					WillReturnRows(rows)
			},
			expectedRet: nil,
			expectedErr: nil,
		},
		{
			name:  "database error",
			brand: brand,
			setupMocks: func(brand string, products []structs.Product) {
				fixture.mock.ExpectQuery(`select \* from product where id_brand in`).
					WithArgs(brand).
					WillReturnError(errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(tt.brand, testProducts)

			ret, err := fixture.repo.GetByBrand(fixture.ctx, tt.brand)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDelete(t *testing.T) {
	t.Parallel()
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(productID uuid.UUID) {
				fixture.mock.ExpectExec(`delete from product where id = \$1`).
					WithArgs(productID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "product not found",
			setupMocks: func(productID uuid.UUID) {
				fixture.mock.ExpectExec(`delete from product where id = \$1`).
					WithArgs(productID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: errors.New("product with id " + testProduct.Id.String() + " not found"),
		},
		{
			name: "database error",
			setupMocks: func(productID uuid.UUID) {
				fixture.mock.ExpectExec(`delete from product where id = \$1`).
					WithArgs(productID).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(testProduct.Id)

			err := fixture.repo.Delete(fixture.ctx, testProduct.Id)
			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
