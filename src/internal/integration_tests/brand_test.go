package integrationtests

import (
	"context"
	"testing"

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
}

func NewBrandTestFixture(t *testing.T) *BrandTestFixture {
	brandRepo := brand_rep.New(db)
	service := brand.New(brandRepo)

	return &BrandTestFixture{
		t:         t,
		ctx:       context.Background(),
		service:   service,
		brandRepo: brandRepo,
	}
}

func (f *BrandTestFixture) createTestBrand() structs.Brand {
	return structs.Brand{
		Name:          "Test Brand",
		Description:   "Test Description",
		PriceCategory: "premium",
	}
}

func (f *BrandTestFixture) createAnotherTestBrand() structs.Brand {
	return structs.Brand{
		Name:          "Another Brand",
		Description:   "Another Description",
		PriceCategory: "budget",
	}
}

func (f *BrandTestFixture) assertBrandEqual(expected, actual structs.Brand) {
	require.Equal(f.t, expected.Name, actual.Name)
	require.Equal(f.t, expected.Description, actual.Description)
	require.Equal(f.t, expected.PriceCategory, actual.PriceCategory)
}

func TestBrand_Create_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)
	testBrand := fixture.createTestBrand()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		brand       structs.Brand
		expectedErr bool
	}{
		{
			name: "successful brand creation",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			brand:       testBrand,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, tt.brand)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBrand_GetById_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully get brand by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				brand := fixture.createTestBrand()
				err := fixture.brandRepo.Create(fixture.ctx, brand)
				require.NoError(t, err)

				// Получаем ID созданного бренда через репозиторий
				var brands []struct {
					ID uuid.UUID `db:"id"`
				}
				err = db.SelectContext(fixture.ctx, &brands, "SELECT id FROM brand WHERE name = $1", brand.Name)
				require.NoError(t, err)
				require.Len(t, brands, 1)
				return brands[0].ID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent brand by id",
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
			brandID := tt.setup()
			defer tt.cleanup()

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
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully delete brand",
			setup: func() uuid.UUID {
				truncateTables(t)
				brand := fixture.createTestBrand()
				err := fixture.brandRepo.Create(fixture.ctx, brand)
				require.NoError(t, err)

				var brands []struct {
					ID uuid.UUID `db:"id"`
				}
				err = db.SelectContext(fixture.ctx, &brands, "SELECT id FROM brand WHERE name = $1", brand.Name)
				require.NoError(t, err)
				require.Len(t, brands, 1)
				return brands[0].ID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent brand",
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
			brandID := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Delete(fixture.ctx, brandID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := fixture.service.GetById(fixture.ctx, brandID)
				require.Error(t, err)
			}
		})
	}
}

func TestBrand_GetAllBrandsInCategory_AAA(t *testing.T) {
	fixture := NewBrandTestFixture(t)
	testBrand1 := fixture.createTestBrand()
	testBrand2 := fixture.createAnotherTestBrand()

	tests := []struct {
		name          string
		setup         func() string
		cleanup       func()
		category      string
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get brands in premium category",
			setup: func() string {
				truncateTables(t)
				err := fixture.brandRepo.Create(fixture.ctx, testBrand1)
				require.NoError(t, err)
				err = fixture.brandRepo.Create(fixture.ctx, testBrand2)
				require.NoError(t, err)
				return "premium"
			},
			cleanup: func() {
				truncateTables(t)
			},
			category:      "premium",
			expectedCount: 1,
			expectedErr:   false,
		},
		{
			name: "successfully get brands in budget category",
			setup: func() string {
				truncateTables(t)
				err := fixture.brandRepo.Create(fixture.ctx, testBrand1)
				require.NoError(t, err)
				err = fixture.brandRepo.Create(fixture.ctx, testBrand2)
				require.NoError(t, err)
				return "budget"
			},
			cleanup: func() {
				truncateTables(t)
			},
			category:      "budget",
			expectedCount: 1,
			expectedErr:   false,
		},
		{
			name: "get empty brands list for non-existent category",
			setup: func() string {
				truncateTables(t)
				err := fixture.brandRepo.Create(fixture.ctx, testBrand1)
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
