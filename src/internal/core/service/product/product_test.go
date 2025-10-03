package product

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.Build()

	tests := []struct {
		name        string
		product     structs.Product
		setupMocks  func(*mock_structs.MockProductRepository, structs.Product)
		expectedErr error
	}{
		{
			name:    "successful creation",
			product: testProduct,
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, product structs.Product) {
				mockRepo.EXPECT().Create(fixture.ctx, product).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "repository error",
			product: testProduct,
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, product structs.Product) {
				mockRepo.EXPECT().Create(fixture.ctx, product).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, tt.product)

			err := service.Create(fixture.ctx, tt.product)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockProductRepository, structs.Product)
		expectedRet structs.Product
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, product structs.Product) {
				mockRepo.EXPECT().GetById(fixture.ctx, product.Id).Return(product, nil)
			},
			expectedRet: testProduct,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, product structs.Product) {
				mockRepo.EXPECT().GetById(fixture.ctx, product.Id).Return(structs.Product{}, errTest)
			},
			expectedRet: structs.Product{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testProduct)

			ret, err := service.GetById(fixture.ctx, testProduct.Id)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestGetByArticule_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.WithArticule("UNIQUE123").Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockProductRepository, structs.Product)
		expectedRet structs.Product
		expectedErr error
	}{
		{
			name: "successful get by articule",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, product structs.Product) {
				mockRepo.EXPECT().GetByArticule(fixture.ctx, product.Articule).Return(product, nil)
			},
			expectedRet: testProduct,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, product structs.Product) {
				mockRepo.EXPECT().GetByArticule(fixture.ctx, product.Articule).Return(structs.Product{}, errTest)
			},
			expectedRet: structs.Product{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testProduct)

			ret, err := service.GetByArticule(fixture.ctx, testProduct.Articule)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestGetByCategory_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	category := "electronics"
	testProducts := []structs.Product{
		fixture.productBuilder.WithCategory(category).Build(),
		fixture.productBuilder.WithCategory(category).WithName("Product 2").Build(),
	}

	tests := []struct {
		name        string
		category    string
		setupMocks  func(*mock_structs.MockProductRepository, string, []structs.Product)
		expectedRet []structs.Product
		expectedErr error
	}{
		{
			name:     "successful get by category",
			category: category,
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, category string, products []structs.Product) {
				mockRepo.EXPECT().GetByCategory(fixture.ctx, category).Return(products, nil)
			},
			expectedRet: testProducts,
			expectedErr: nil,
		},
		{
			name:     "empty category",
			category: "nonexistent",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, category string, products []structs.Product) {
				mockRepo.EXPECT().GetByCategory(fixture.ctx, category).Return([]structs.Product{}, nil)
			},
			expectedRet: []structs.Product{},
			expectedErr: nil,
		},
		{
			name:     "repository error",
			category: category,
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, category string, products []structs.Product) {
				mockRepo.EXPECT().GetByCategory(fixture.ctx, category).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, tt.category, testProducts)

			ret, err := service.GetByCategory(fixture.ctx, tt.category)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestGetByBrand_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	brand := "TestBrand"
	brandId := structs.GenId()
	testProducts := []structs.Product{
		fixture.productBuilder.WithIdBrand(brandId).Build(),
		fixture.productBuilder.WithIdBrand(brandId).WithName("Product 2").Build(),
	}

	tests := []struct {
		name        string
		brand       string
		setupMocks  func(*mock_structs.MockProductRepository, string, []structs.Product)
		expectedRet []structs.Product
		expectedErr error
	}{
		{
			name:  "successful get by brand",
			brand: brand,
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, brand string, products []structs.Product) {
				mockRepo.EXPECT().GetByBrand(fixture.ctx, brand).Return(products, nil)
			},
			expectedRet: testProducts,
			expectedErr: nil,
		},
		{
			name:  "empty brand",
			brand: "SomeBrand",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, brand string, products []structs.Product) {
				mockRepo.EXPECT().GetByBrand(fixture.ctx, brand).Return([]structs.Product{}, nil)
			},
			expectedRet: []structs.Product{},
			expectedErr: nil,
		},
		{
			name:  "repository error",
			brand: brand,
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, brand string, products []structs.Product) {
				mockRepo.EXPECT().GetByBrand(fixture.ctx, brand).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, tt.brand, testProducts)

			ret, err := service.GetByBrand(fixture.ctx, tt.brand)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestDelete_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testProduct := fixture.productBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockProductRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, productID uuid.UUID) {
				mockRepo.EXPECT().Delete(fixture.ctx, productID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockProductRepository, productID uuid.UUID) {
				mockRepo.EXPECT().Delete(fixture.ctx, productID).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testProduct.Id)

			err := service.Delete(fixture.ctx, testProduct.Id)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}
