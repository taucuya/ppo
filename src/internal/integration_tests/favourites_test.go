package integrationtests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/favourites"
	"github.com/taucuya/ppo/internal/core/structs"
	favourites_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/favourites"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

type FavouritesTestFixture struct {
	t              *testing.T
	ctx            context.Context
	service        *favourites.Service
	favouritesRepo *favourites_rep.Repository
	userRepo       *user_rep.Repository
}

func NewFavouritesTestFixture(t *testing.T) *FavouritesTestFixture {
	favouritesRepo := favourites_rep.New(db)
	userRepo := user_rep.New(db)

	service := favourites.New(favouritesRepo)

	return &FavouritesTestFixture{
		t:              t,
		ctx:            context.Background(),
		service:        service,
		favouritesRepo: favouritesRepo,
		userRepo:       userRepo,
	}
}

func (f *FavouritesTestFixture) createTestUser() structs.User {
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	return structs.User{
		Name:          "Test User",
		Date_of_birth: dob,
		Mail:          "test@example.com",
		Password:      "password123",
		Phone:         "89016475843",
		Address:       "123 Test St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *FavouritesTestFixture) createTestFavourites(userID uuid.UUID) structs.Favourites {
	return structs.Favourites{
		IdUser: userID,
	}
}

func (f *FavouritesTestFixture) createTestFavouritesItem(favouritesID, productID uuid.UUID) structs.FavouritesItem {
	return structs.FavouritesItem{
		IdProduct:    productID,
		IdFavourites: favouritesID,
	}
}

func (f *FavouritesTestFixture) createTestProduct() uuid.UUID {
	productID := uuid.New()
	brandID := uuid.New()
	_, err := db.Exec("INSERT INTO brand (id, name) VALUES ($1, $2)", brandID, "Test Brand")
	require.NoError(f.t, err)
	_, err = db.Exec("INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		productID, "Test Product", "Test Description", 1000, brandID, 10, "TEST-ART-123")
	require.NoError(f.t, err)
	return productID
}

func (f *FavouritesTestFixture) setupUserWithFavourites() (uuid.UUID, uuid.UUID) {
	testUser := f.createTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)

	testFavourites := f.createTestFavourites(userID)
	err = f.favouritesRepo.Create(f.ctx, testFavourites)
	require.NoError(f.t, err)

	favouritesID, err := f.favouritesRepo.GetFIdByUId(f.ctx, userID)
	require.NoError(f.t, err)

	return userID, favouritesID
}

func TestFavourites_Create_AAA(t *testing.T) {
	fixture := NewFavouritesTestFixture(t)

	tests := []struct {
		name        string
		setup       func() structs.Favourites
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful favourites creation",
			setup: func() structs.Favourites {
				truncateTables(t)
				testUser := fixture.createTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
				return fixture.createTestFavourites(userID)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to create favourites for non-existent user",
			setup: func() structs.Favourites {
				truncateTables(t)
				return fixture.createTestFavourites(uuid.New())
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			favourites := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, favourites)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFavourites_GetById_AAA(t *testing.T) {
	fixture := NewFavouritesTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully get favourites by user id",
			setup: func() uuid.UUID {
				truncateTables(t)
				userID, _ := fixture.setupUserWithFavourites()
				return userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to get favourites for non-existent user",
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
			userID := tt.setup()
			defer tt.cleanup()

			result, err := fixture.service.GetById(fixture.ctx, userID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, userID, result.IdUser)
			}
		})
	}
}

func TestFavourites_GetItems_AAA(t *testing.T) {
	fixture := NewFavouritesTestFixture(t)

	tests := []struct {
		name          string
		setup         func() uuid.UUID
		cleanup       func()
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get empty favourites items",
			setup: func() uuid.UUID {
				truncateTables(t)
				userID, _ := fixture.setupUserWithFavourites()
				return userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
		{
			name: "successfully get favourites with items",
			setup: func() uuid.UUID {
				truncateTables(t)
				userID, favouritesID := fixture.setupUserWithFavourites()

				brandID := uuid.New()
				product1 := uuid.New()
				product2 := uuid.New()

				_, err := db.Exec(`INSERT INTO brand (id, name) VALUES ($1, $2)`, brandID, "Test Brand")
				require.NoError(t, err)

				_, err = db.Exec(`INSERT INTO product (id, name, description, price, id_brand, amount, art) 
					VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)`,
					product1, "Product 1", "Desc 1", 1000, brandID, 10, "ART-001",
					product2, "Product 2", "Desc 2", 2000, brandID, 5, "ART-002")
				require.NoError(t, err)

				item1 := fixture.createTestFavouritesItem(favouritesID, product1)
				item2 := fixture.createTestFavouritesItem(favouritesID, product2)

				err = fixture.favouritesRepo.AddItem(fixture.ctx, item1)
				require.NoError(t, err)
				err = fixture.favouritesRepo.AddItem(fixture.ctx, item2)
				require.NoError(t, err)

				return userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "fail to get items for non-existent user",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			defer tt.cleanup()

			items, err := fixture.service.GetItems(fixture.ctx, userID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, items, tt.expectedCount)
			}
		})
	}
}

func TestFavourites_AddItem_AAA(t *testing.T) {
	fixture := NewFavouritesTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.FavouritesItem, uuid.UUID)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully add item to favourites",
			setup: func() (structs.FavouritesItem, uuid.UUID) {
				truncateTables(t)
				userID, favouritesID := fixture.setupUserWithFavourites()
				productID := fixture.createTestProduct()
				item := fixture.createTestFavouritesItem(favouritesID, productID)
				return item, userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to add item to non-existent user favourites",
			setup: func() (structs.FavouritesItem, uuid.UUID) {
				truncateTables(t)
				productID := fixture.createTestProduct()
				item := fixture.createTestFavouritesItem(uuid.New(), productID)
				return item, uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, userID := tt.setup()
			defer tt.cleanup()

			err := fixture.service.AddItem(fixture.ctx, item, userID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				items, err := fixture.service.GetItems(fixture.ctx, userID)
				require.NoError(t, err)
				require.Len(t, items, 1)
				require.Equal(t, item.IdProduct, items[0].IdProduct)
			}
		})
	}
}

func TestFavourites_DeleteItem_AAA(t *testing.T) {
	fixture := NewFavouritesTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (uuid.UUID, uuid.UUID)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "fail to delete non-existent item",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				userID, _ := fixture.setupUserWithFavourites()
				return userID, uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
		{
			name: "fail to delete item from non-existent user",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				return uuid.New(), uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, productID := tt.setup()
			defer tt.cleanup()

			err := fixture.service.DeleteItem(fixture.ctx, userID, productID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				items, err := fixture.service.GetItems(fixture.ctx, userID)
				require.NoError(t, err)
				require.Len(t, items, 0)
			}
		})
	}
}
