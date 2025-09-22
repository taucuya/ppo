package brand

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockBrandRepository)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().Create(fixture.ctx, fixture.brand).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().Create(fixture.ctx, fixture.brand).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.Create(fixture.ctx, fixture.brand)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockBrandRepository)
		expectedRet structs.Brand
		expectedErr error
	}{
		{
			name: "successful get",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.brand.Id).Return(fixture.brand, nil)
			},
			expectedRet: fixture.brand,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.brand.Id).Return(structs.Brand{}, errTest)
			},
			expectedRet: structs.Brand{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetById(fixture.ctx, fixture.brand.Id)

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

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockBrandRepository)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().Delete(fixture.ctx, fixture.brand.Id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().Delete(fixture.ctx, fixture.brand.Id).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.Delete(fixture.ctx, fixture.brand.Id)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetAllBrandsInCategory_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testBrands := []structs.Brand{
		{
			Id:            structs.GenId(),
			Name:          "Brand 1",
			Description:   "Description 1",
			PriceCategory: "Premium",
		},
		{
			Id:            structs.GenId(),
			Name:          "Brand 2",
			Description:   "Description 2",
			PriceCategory: "Premium",
		},
	}

	tests := []struct {
		name        string
		category    string
		setupMocks  func(*mock_structs.MockBrandRepository)
		expectedRet []structs.Brand
		expectedErr error
	}{
		{
			name:     "successful get all brands in category",
			category: "Premium",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().GetAllBrandsInCategory(fixture.ctx, "Premium").Return(testBrands, nil)
			},
			expectedRet: testBrands,
			expectedErr: nil,
		},
		{
			name:     "empty category",
			category: "Economy",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().GetAllBrandsInCategory(fixture.ctx, "Economy").Return([]structs.Brand{}, nil)
			},
			expectedRet: []structs.Brand{},
			expectedErr: nil,
		},
		{
			name:     "repository error",
			category: "Premium",
			setupMocks: func(mockRepo *mock_structs.MockBrandRepository) {
				mockRepo.EXPECT().GetAllBrandsInCategory(fixture.ctx, "Premium").Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetAllBrandsInCategory(fixture.ctx, tt.category)

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
