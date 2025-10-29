package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/brand"
	"github.com/taucuya/ppo/internal/core/structs"
	brand_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/brand"
)

type BrandTestFixture struct {
	t         *testing.T
	ctx       context.Context
	service   *brand.Service
	brandRepo *brand_rep.Repository
	testID    string
}

func NewBrandTestFixture(t *testing.T) *BrandTestFixture {
	brandRepo := brand_rep.New(db)
	service := brand.New(brandRepo)

	testID := uuid.New().String()[:8]

	return &BrandTestFixture{
		t:         t,
		ctx:       context.Background(),
		service:   service,
		brandRepo: brandRepo,
		testID:    testID,
	}
}

func (f *BrandTestFixture) generateTestBrand() structs.Brand {
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%s-%d", f.testID, timestamp)

	return structs.Brand{
		Name:          fmt.Sprintf("Test Brand %s", uniqueID),
		Description:   "Test Description",
		PriceCategory: "premium",
	}
}

func (f *BrandTestFixture) generateAnotherTestBrand() structs.Brand {
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%s-%d", f.testID, timestamp)

	return structs.Brand{
		Name:          fmt.Sprintf("Another Brand %s", uniqueID),
		Description:   "Another Description",
		PriceCategory: "budget",
	}
}

func (f *BrandTestFixture) cleanupBrandData(brandID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM product WHERE id_brand = $1", brandID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM brand WHERE id = $1", brandID)
}

func (f *BrandTestFixture) createBrandForTest() (uuid.UUID, structs.Brand) {
	testBrand := f.generateTestBrand()
	err := f.brandRepo.Create(f.ctx, testBrand)
	require.NoError(f.t, err)

	var brandID uuid.UUID
	err = db.GetContext(f.ctx, &brandID, "SELECT id FROM brand WHERE name = $1", testBrand.Name)
	require.NoError(f.t, err)

	return brandID, testBrand
}

func TestBrand_Create_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)

	tests := []struct {
		name        string
		setup       func() structs.Brand
		expectedErr bool
	}{
		{
			name: "successful brand creation",
			setup: func() structs.Brand {
				return fixture.generateTestBrand()
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brand := tt.setup()

			err := fixture.service.Create(fixture.ctx, brand)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var createdBrandID uuid.UUID
				err = db.GetContext(fixture.ctx, &createdBrandID,
					"SELECT id FROM brand WHERE name = $1", brand.Name)
				require.NoError(t, err)
				defer fixture.cleanupBrandData(createdBrandID)
			}
		})
	}
}

func TestBrand_GetById_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully get brand by id",
			setup: func() uuid.UUID {
				brandID, _ := fixture.createBrandForTest()
				return brandID
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent brand by id",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brandID := tt.setup()
			if brandID != uuid.Nil {
				defer fixture.cleanupBrandData(brandID)
			}

			result, err := fixture.service.GetById(fixture.ctx, brandID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, brandID, result.Id)
			}
		})
	}
}

func TestBrand_Delete_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully delete brand",
			setup: func() uuid.UUID {
				brandID, _ := fixture.createBrandForTest()
				return brandID
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent brand",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brandID := tt.setup()
			// Не очищаем здесь, т.к. тест сам удаляет бренд

			err := fixture.service.Delete(fixture.ctx, brandID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Проверяем что бренд действительно удален
				_, err := fixture.service.GetById(fixture.ctx, brandID)
				require.Error(t, err)
			}
		})
	}
}

func TestBrand_GetAllBrandsInCategory_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)

	tests := []struct {
		name          string
		setup         func() (string, []uuid.UUID)
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get brands in budget category",
			setup: func() (string, []uuid.UUID) {
				brandIDs := []uuid.UUID{}

				premiumBrand := fixture.generateTestBrand()
				premiumBrand.PriceCategory = "premium"
				err := fixture.brandRepo.Create(fixture.ctx, premiumBrand)
				require.NoError(t, err)
				var premiumBrandID uuid.UUID
				err = db.GetContext(fixture.ctx, &premiumBrandID, "SELECT id FROM brand WHERE name = $1", premiumBrand.Name)
				require.NoError(t, err)
				brandIDs = append(brandIDs, premiumBrandID)

				budgetBrand := fixture.generateAnotherTestBrand()
				budgetBrand.PriceCategory = "budget"
				err = fixture.brandRepo.Create(fixture.ctx, budgetBrand)
				require.NoError(t, err)
				var budgetBrandID uuid.UUID
				err = db.GetContext(fixture.ctx, &budgetBrandID, "SELECT id FROM brand WHERE name = $1", budgetBrand.Name)
				require.NoError(t, err)
				brandIDs = append(brandIDs, budgetBrandID)

				return "budget", brandIDs
			},
			expectedCount: 1,
			expectedErr:   false,
		},
		{
			name: "get empty brands list for non-existent category",
			setup: func() (string, []uuid.UUID) {
				brandIDs := []uuid.UUID{}

				brand := fixture.generateTestBrand()
				brand.PriceCategory = "premium"
				err := fixture.brandRepo.Create(fixture.ctx, brand)
				require.NoError(t, err)
				var brandID uuid.UUID
				err = db.GetContext(fixture.ctx, &brandID, "SELECT id FROM brand WHERE name = $1", brand.Name)
				require.NoError(t, err)
				brandIDs = append(brandIDs, brandID)

				return "nonexistent", brandIDs
			},
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, brandIDs := tt.setup()

			// Очищаем созданные бренды после теста
			defer func() {
				for _, brandID := range brandIDs {
					fixture.cleanupBrandData(brandID)
				}
			}()

			brands, err := fixture.service.GetAllBrandsInCategory(fixture.ctx, category)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, brands, tt.expectedCount)

				if tt.expectedCount > 0 {
					for _, brand := range brands {
						require.Equal(t, category, brand.PriceCategory)
					}
				}
			}
		})
	}
}
