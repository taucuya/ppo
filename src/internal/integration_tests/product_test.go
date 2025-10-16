package integrationtests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/product"
	"github.com/taucuya/ppo/internal/core/structs"
	product_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/product"
)

type ProductTestFixture struct {
	t           *testing.T
	ctx         context.Context
	service     *product.Service
	productRepo *product_rep.Repository
}

func NewProductTestFixture(t *testing.T) *ProductTestFixture {
	productRepo := product_rep.New(db)
	service := product.New(productRepo)

	return &ProductTestFixture{
		t:           t,
		ctx:         context.Background(),
		service:     service,
		productRepo: productRepo,
	}
}

func (f *ProductTestFixture) createTestBrand() uuid.UUID {
	brandID := uuid.New()
	_, err := db.Exec("INSERT INTO brand (id, name, description, price_category) VALUES ($1, $2, $3, $4)",
		brandID, "Test Brand", "Test Brand Description", "premium")
	require.NoError(f.t, err)
	return brandID
}

func (f *ProductTestFixture) createTestProduct(brandID uuid.UUID) structs.Product {
	return structs.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       1000.0,
		Category:    "electronics",
		Amount:      10,
		IdBrand:     brandID,
		PicLink:     "http://example.com/pic.jpg",
		Articule:    "TEST-ART-001",
	}
}

func (f *ProductTestFixture) createAnotherTestProduct(brandID uuid.UUID) structs.Product {
	return structs.Product{
		Name:        "Another Product",
		Description: "Another Description",
		Price:       2000.0,
		Category:    "electronics",
		Amount:      5,
		IdBrand:     brandID,
		PicLink:     "http://example.com/another.jpg",
		Articule:    "TEST-ART-002",
	}
}

func (f *ProductTestFixture) createDifferentCategoryProduct(brandID uuid.UUID) structs.Product {
	return structs.Product{
		Name:        "Different Category Product",
		Description: "Different Category Description",
		Price:       1500.0,
		Category:    "clothing",
		Amount:      8,
		IdBrand:     brandID,
		PicLink:     "http://example.com/different.jpg",
		Articule:    "TEST-ART-003",
	}
}

func (f *ProductTestFixture) setupProduct() uuid.UUID {
	brandID := f.createTestBrand()
	testProduct := f.createTestProduct(brandID)
	err := f.productRepo.Create(f.ctx, testProduct)
	require.NoError(f.t, err)

	var products []struct {
		ID uuid.UUID `db:"id"`
	}
	err = db.SelectContext(f.ctx, &products, "SELECT id FROM product WHERE name = $1", testProduct.Name)
	require.NoError(f.t, err)
	require.Len(f.t, products, 1)
	return products[0].ID
}

func TestProduct_Create_AAA(t *testing.T) {
	fixture := NewProductTestFixture(t)

	tests := []struct {
		name        string
		setup       func() structs.Product
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful product creation",
			setup: func() structs.Product {
				truncateTables(t)
				brandID := fixture.createTestBrand()
				return fixture.createTestProduct(brandID)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to create product with duplicate articule",
			setup: func() structs.Product {
				truncateTables(t)
				brandID := fixture.createTestBrand()
				product1 := fixture.createTestProduct(brandID)
				err := fixture.productRepo.Create(fixture.ctx, product1)
				require.NoError(t, err)

				product2 := fixture.createAnotherTestProduct(brandID)
				product2.Articule = product1.Articule
				return product2
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, product)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProduct_GetById_AAA(t *testing.T) {
	fixture := NewProductTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully get product by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				return fixture.setupProduct()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent product by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productID := tt.setup()
			defer tt.cleanup()

			result, err := fixture.service.GetById(fixture.ctx, productID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, productID, result.Id)
			}
		})
	}
}

func TestProduct_GetByArticule_AAA(t *testing.T) {
	fixture := NewProductTestFixture(t)

	tests := []struct {
		name        string
		setup       func() string
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully get product by articule",
			setup: func() string {
				truncateTables(t)
				fixture.setupProduct()
				return "TEST-ART-001"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to get product by non-existent articule",
			setup: func() string {
				truncateTables(t)
				return "NON-EXISTENT-ART"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articule := tt.setup()
			defer tt.cleanup()

			result, err := fixture.service.GetByArticule(fixture.ctx, articule)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, articule, result.Articule)
			}
		})
	}
}

func TestProduct_GetByCategory_AAA(t *testing.T) {
	fixture := NewProductTestFixture(t)

	tests := []struct {
		name          string
		setup         func() string
		cleanup       func()
		category      string
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get products by category",
			setup: func() string {
				truncateTables(t)
				brandID := fixture.createTestBrand()

				product1 := fixture.createTestProduct(brandID)
				err := fixture.productRepo.Create(fixture.ctx, product1)
				require.NoError(t, err)

				product2 := fixture.createAnotherTestProduct(brandID)
				err = fixture.productRepo.Create(fixture.ctx, product2)
				require.NoError(t, err)

				product3 := fixture.createDifferentCategoryProduct(brandID)
				err = fixture.productRepo.Create(fixture.ctx, product3)
				require.NoError(t, err)

				return "electronics"
			},
			cleanup: func() {
				truncateTables(t)
			},
			category:      "electronics",
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "get empty products list for non-existent category",
			setup: func() string {
				truncateTables(t)
				brandID := fixture.createTestBrand()
				product := fixture.createTestProduct(brandID)
				err := fixture.productRepo.Create(fixture.ctx, product)
				require.NoError(t, err)
				return "nonexistent"
			},
			cleanup: func() {
				truncateTables(t)
			},
			category:      "nonexistent",
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := tt.setup()
			defer tt.cleanup()

			products, err := fixture.service.GetByCategory(fixture.ctx, category)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, products, tt.expectedCount)

				if tt.expectedCount > 0 {
					for _, product := range products {
						require.Equal(t, category, product.Category)
					}
				}
			}
		})
	}
}

func TestProduct_GetByBrand_AAA(t *testing.T) {
	fixture := NewProductTestFixture(t)

	tests := []struct {
		name          string
		setup         func() string
		cleanup       func()
		brand         string
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get products by brand",
			setup: func() string {
				truncateTables(t)
				brandID := fixture.createTestBrand()

				product1 := fixture.createTestProduct(brandID)
				err := fixture.productRepo.Create(fixture.ctx, product1)
				require.NoError(t, err)

				product2 := fixture.createAnotherTestProduct(brandID)
				err = fixture.productRepo.Create(fixture.ctx, product2)
				require.NoError(t, err)

				return "Test Brand"
			},
			cleanup: func() {
				truncateTables(t)
			},
			brand:         "Test Brand",
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "get empty products list for non-existent brand",
			setup: func() string {
				truncateTables(t)
				brandID := fixture.createTestBrand()
				product := fixture.createTestProduct(brandID)
				err := fixture.productRepo.Create(fixture.ctx, product)
				require.NoError(t, err)
				return "Non-existent Brand"
			},
			cleanup: func() {
				truncateTables(t)
			},
			brand:         "Non-existent Brand",
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brand := tt.setup()
			defer tt.cleanup()

			products, err := fixture.service.GetByBrand(fixture.ctx, brand)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, products, tt.expectedCount)
			}
		})
	}
}

func TestProduct_Delete_AAA(t *testing.T) {
	fixture := NewProductTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully delete product",
			setup: func() uuid.UUID {
				truncateTables(t)
				return fixture.setupProduct()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent product",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productID := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Delete(fixture.ctx, productID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := fixture.service.GetById(fixture.ctx, productID)
				require.Error(t, err)
			}
		})
	}
}
