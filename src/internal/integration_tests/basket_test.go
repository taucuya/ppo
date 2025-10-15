package integrationtests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/basket"
	"github.com/taucuya/ppo/internal/core/structs"
	basket_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/basket"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

type BasketTestFixture struct {
	t          *testing.T
	ctx        context.Context
	service    *basket.Service
	basketRepo *basket_rep.Repository
	userRepo   *user_rep.Repository
}

func NewBasketTestFixture(t *testing.T) *BasketTestFixture {
	basketRepo := basket_rep.New(db)
	userRepo := user_rep.New(db)

	service := basket.New(basketRepo)

	return &BasketTestFixture{
		t:          t,
		ctx:        context.Background(),
		service:    service,
		basketRepo: basketRepo,
		userRepo:   userRepo,
	}
}

func (f *BasketTestFixture) createTestUser() structs.User {
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

func (f *BasketTestFixture) createTestBasket(userID uuid.UUID) structs.Basket {
	return structs.Basket{
		Id:     uuid.New(),
		IdUser: userID,
		Date:   time.Now(),
	}
}

func (f *BasketTestFixture) createTestProduct() uuid.UUID {
	productID := uuid.New()
	brandID := uuid.New()
	_, err := db.Exec("INSERT INTO brand (id, name) VALUES ($1, $2)", brandID, "Test Brand")
	require.NoError(f.t, err)
	_, err = db.Exec("INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		productID, "Test Product", "Test Description", 1000, brandID, 10, "TEST-ART-123")
	require.NoError(f.t, err)
	return productID
}

func (f *BasketTestFixture) setupUserWithBasket() (uuid.UUID, uuid.UUID) {
	testUser := f.createTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)

	testBasket := structs.Basket{
		IdUser: userID,
		Date:   time.Now(),
	}
	err = f.basketRepo.Create(f.ctx, testBasket)
	require.NoError(f.t, err)

	basketID, err := f.basketRepo.GetBIdByUId(f.ctx, userID)
	require.NoError(f.t, err)

	return userID, basketID
}

func (f *BasketTestFixture) createTestBasketItem(basketID, productID uuid.UUID) structs.BasketItem {
	return structs.BasketItem{
		IdProduct: productID,
		IdBasket:  basketID,
		Amount:    2,
	}
}
func TestBasket_Create_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() structs.Basket
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful basket creation",
			setup: func() structs.Basket {
				truncateTables(t)
				testUser := fixture.createTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
				return fixture.createTestBasket(userID)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to create basket for non-existent user",
			setup: func() structs.Basket {
				truncateTables(t)
				return fixture.createTestBasket(uuid.New())
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basket := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, basket)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBasket_GetById_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully get basket by user id",
			setup: func() uuid.UUID {
				truncateTables(t)
				userID, _ := fixture.setupUserWithBasket()
				return userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to get basket for non-existent user",
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

func TestBasket_GetItems_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name          string
		setup         func() uuid.UUID
		cleanup       func()
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get empty basket items",
			setup: func() uuid.UUID {
				truncateTables(t)
				userID, _ := fixture.setupUserWithBasket()
				return userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
		{
			name: "successfully get basket with items",
			setup: func() uuid.UUID {
				truncateTables(t)
				userID, basketID := fixture.setupUserWithBasket()

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

				item1 := fixture.createTestBasketItem(basketID, product1)
				item2 := fixture.createTestBasketItem(basketID, product2)

				err = fixture.basketRepo.AddItem(fixture.ctx, item1)
				require.NoError(t, err)
				err = fixture.basketRepo.AddItem(fixture.ctx, item2)
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

func TestBasket_AddItem_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.BasketItem, uuid.UUID)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully add item to basket",
			setup: func() (structs.BasketItem, uuid.UUID) {
				truncateTables(t)
				userID, basketID := fixture.setupUserWithBasket()
				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(basketID, productID)
				return item, userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "successfully update amount when adding existing item",
			setup: func() (structs.BasketItem, uuid.UUID) {
				truncateTables(t)
				userID, basketID := fixture.setupUserWithBasket()
				productID := fixture.createTestProduct()

				item := fixture.createTestBasketItem(basketID, productID)
				err := fixture.basketRepo.AddItem(fixture.ctx, item)
				require.NoError(t, err)

				sameItem := fixture.createTestBasketItem(basketID, productID)
				sameItem.Amount = 3
				return sameItem, userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to add item to non-existent user basket",
			setup: func() (structs.BasketItem, uuid.UUID) {
				truncateTables(t)
				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(uuid.New(), productID)
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

				if tt.name == "successfully update amount when adding existing item" {
					require.Len(t, items, 1)
					require.Equal(t, 5, items[0].Amount)
				} else {
					require.Len(t, items, 1)
					require.Equal(t, item.IdProduct, items[0].IdProduct)
				}
			}
		})
	}
}

func TestBasket_DeleteItem_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (uuid.UUID, uuid.UUID)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully delete item from basket",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				userID, basketID := fixture.setupUserWithBasket()
				productID := fixture.createTestProduct()

				item := fixture.createTestBasketItem(basketID, productID)
				err := fixture.basketRepo.AddItem(fixture.ctx, item)
				require.NoError(t, err)

				return userID, productID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent item",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				userID, _ := fixture.setupUserWithBasket()
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

func TestBasket_UpdateItemAmount_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (uuid.UUID, uuid.UUID, int)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully update item amount",
			setup: func() (uuid.UUID, uuid.UUID, int) {
				truncateTables(t)
				userID, basketID := fixture.setupUserWithBasket()
				productID := fixture.createTestProduct()

				item := fixture.createTestBasketItem(basketID, productID)
				err := fixture.basketRepo.AddItem(fixture.ctx, item)
				require.NoError(t, err)

				return userID, productID, 5
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to update amount for non-existent item",
			setup: func() (uuid.UUID, uuid.UUID, int) {
				truncateTables(t)
				userID, _ := fixture.setupUserWithBasket()
				return userID, uuid.New(), 5
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to update amount for non-existent user",
			setup: func() (uuid.UUID, uuid.UUID, int) {
				truncateTables(t)
				return uuid.New(), uuid.New(), 5
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, productID, amount := tt.setup()
			defer tt.cleanup()

			err := fixture.service.UpdateItemAmount(fixture.ctx, userID, productID, amount)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.name == "successfully update item amount" {
					items, err := fixture.service.GetItems(fixture.ctx, userID)
					require.NoError(t, err)
					require.Len(t, items, 1)
					require.Equal(t, amount, items[0].Amount)
				}
			}
		})
	}
}
