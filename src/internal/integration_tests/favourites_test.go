package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

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
	testID         string
}

func NewFavouritesTestFixture(t *testing.T) *FavouritesTestFixture {
	favouritesRepo := favourites_rep.New(db)
	userRepo := user_rep.New(db)
	service := favourites.New(favouritesRepo)

	testID := uuid.New().String()[:8]

	return &FavouritesTestFixture{
		t:              t,
		ctx:            context.Background(),
		service:        service,
		favouritesRepo: favouritesRepo,
		userRepo:       userRepo,
		testID:         testID,
	}
}

func (f *FavouritesTestFixture) generateTestUser() (structs.User, string) {
	timestamp := time.Now().UnixNano()
	randomUUID := uuid.New().String()[:8]
	uniqueID := fmt.Sprintf("%s-%d-%s", f.testID, timestamp, randomUUID)

	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	phonePrefix := "89"
	randomNumbers := fmt.Sprintf("%09d", timestamp%1000000000)
	if len(randomNumbers) > 9 {
		randomNumbers = randomNumbers[:9]
	}
	phone := phonePrefix + randomNumbers

	return structs.User{
		Name:          fmt.Sprintf("Test User %s", uniqueID),
		Date_of_birth: dob,
		Mail:          fmt.Sprintf("test%s@example.com", uniqueID),
		Password:      string(hashedPassword),
		Phone:         phone,
		Address:       fmt.Sprintf("123 Test St %s", uniqueID),
		Status:        "active",
		Role:          "обычный пользователь",
	}, plainPassword
}

func (f *FavouritesTestFixture) createUserForTest() (uuid.UUID, structs.User, string) {
	testUser, plainPassword := f.generateTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)
	return userID, testUser, plainPassword
}

func (f *FavouritesTestFixture) cleanupUserData(userID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM favourites_item WHERE id_favourites IN (SELECT id FROM favourites WHERE id_user = $1)", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM favourites WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"user\" WHERE id = $1", userID)
}

func (f *FavouritesTestFixture) createTestProduct() uuid.UUID {
	productID := uuid.New()
	brandID := uuid.New()

	uniqueArt := fmt.Sprintf("TEST-ART-%s", uuid.New().String()[:8])

	_, err := db.ExecContext(f.ctx, "INSERT INTO brand (id, name) VALUES ($1, $2)", brandID, "Test Brand")
	require.NoError(f.t, err)

	_, err = db.ExecContext(f.ctx,
		"INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		productID, "Test Product", "Test Description", 1000.00, brandID, 10, uniqueArt)
	require.NoError(f.t, err)

	return productID
}

func (f *FavouritesTestFixture) createTestFavouritesItem(favouritesID, productID uuid.UUID) structs.FavouritesItem {
	return structs.FavouritesItem{
		IdProduct:    productID,
		IdFavourites: favouritesID,
	}
}

func (f *FavouritesTestFixture) createFavouritesForUser(userID uuid.UUID) (uuid.UUID, error) {
	favourites := structs.Favourites{
		IdUser: userID,
	}
	err := f.favouritesRepo.Create(f.ctx, favourites)
	if err != nil {
		return uuid.Nil, err
	}

	favouritesID, err := f.favouritesRepo.GetFIdByUId(f.ctx, userID)
	return favouritesID, err
}

func TestFavourites_Create_AAA(t *testing.T) {
	fixture := NewFavouritesTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.Favourites, uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successful favourites creation",
			setup: func() (structs.Favourites, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()
				favourites := structs.Favourites{
					IdUser: userID,
				}
				return favourites, userID
			},
			expectedErr: false,
		},
		{
			name: "fail to create favourites for non-existent user",
			setup: func() (structs.Favourites, uuid.UUID) {
				favourites := structs.Favourites{
					IdUser: uuid.New(),
				}
				return favourites, uuid.Nil
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			favourites, userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
		expectedErr bool
	}{
		{
			name: "successfully get favourites by user id",
			setup: func() uuid.UUID {
				userID, _, _ := fixture.createUserForTest()
				_, err := fixture.createFavouritesForUser(userID)
				require.NoError(t, err)
				return userID
			},
			expectedErr: false,
		},
		{
			name: "fail to get favourites for non-existent user",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get favourites with items",
			setup: func() uuid.UUID {
				userID, _, _ := fixture.createUserForTest()
				favouritesID, err := fixture.createFavouritesForUser(userID)
				require.NoError(t, err)

				product1 := fixture.createTestProduct()
				product2 := fixture.createTestProduct()

				item1 := fixture.createTestFavouritesItem(favouritesID, product1)
				item2 := fixture.createTestFavouritesItem(favouritesID, product2)

				err = fixture.service.AddItem(fixture.ctx, item1, userID)
				require.NoError(t, err)
				err = fixture.service.AddItem(fixture.ctx, item2, userID)
				require.NoError(t, err)

				return userID
			},
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "fail to get items for non-existent user",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedCount: 0,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
		expectedErr bool
	}{
		{
			name: "successfully add item to favourites",
			setup: func() (structs.FavouritesItem, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()
				favouritesID, err := fixture.createFavouritesForUser(userID)
				require.NoError(t, err)

				productID := fixture.createTestProduct()
				item := fixture.createTestFavouritesItem(favouritesID, productID)

				return item, userID
			},
			expectedErr: false,
		},
		{
			name: "fail to add item to non-existent user favourites",
			setup: func() (structs.FavouritesItem, uuid.UUID) {
				productID := fixture.createTestProduct()
				item := fixture.createTestFavouritesItem(uuid.New(), productID)
				return item, uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
		expectedErr bool
	}{
		{
			name: "fail to delete non-existent item",
			setup: func() (uuid.UUID, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()
				_, err := fixture.createFavouritesForUser(userID)
				require.NoError(t, err)
				return userID, uuid.New()
			},
			expectedErr: true,
		},
		{
			name: "fail to delete item from non-existent user",
			setup: func() (uuid.UUID, uuid.UUID) {
				return uuid.New(), uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, productID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
