package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

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
	testID     string
}

func NewBasketTestFixture(t *testing.T) *BasketTestFixture {
	basketRepo := basket_rep.New(db)
	userRepo := user_rep.New(db)
	service := basket.New(basketRepo)

	testID := uuid.New().String()[:8]

	return &BasketTestFixture{
		t:          t,
		ctx:        context.Background(),
		service:    service,
		basketRepo: basketRepo,
		userRepo:   userRepo,
		testID:     testID,
	}
}

func (f *BasketTestFixture) generateTestUser() (structs.User, string) {
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

func (f *BasketTestFixture) createUserForTest() (uuid.UUID, structs.User, string) {
	testUser, plainPassword := f.generateTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)
	return userID, testUser, plainPassword
}

func (f *BasketTestFixture) cleanupUserData(userID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM basket_item WHERE id_basket IN (SELECT id FROM basket WHERE id_user = $1)", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM basket WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"user\" WHERE id = $1", userID)
}

func (f *BasketTestFixture) createTestProduct() uuid.UUID {
	productID := uuid.New()
	brandID := uuid.New()

	uniqueArt := fmt.Sprintf("%s-PROD-%d", f.testID, time.Now().UnixNano())

	_, err := db.ExecContext(f.ctx, "INSERT INTO brand (id, name) VALUES ($1, $2)", brandID, fmt.Sprintf("Test Brand %s", f.testID))
	require.NoError(f.t, err)

	_, err = db.ExecContext(f.ctx,
		"INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		productID, fmt.Sprintf("Test Product %s", uniqueArt), "Test Description", 1000.00, brandID, 10, uniqueArt)
	require.NoError(f.t, err)

	return productID
}

func (f *BasketTestFixture) createTestBasketItem(basketID, productID uuid.UUID) structs.BasketItem {
	return structs.BasketItem{
		Id:        uuid.New(),
		IdProduct: productID,
		IdBasket:  basketID,
		Amount:    2,
	}
}

func TestBasket_Create_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.Basket, uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successful basket creation",
			setup: func() (structs.Basket, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()
				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				return basket, userID
			},
			expectedErr: false,
		},
		{
			name: "fail to create basket for non-existent user",
			setup: func() (structs.Basket, uuid.UUID) {
				basket := structs.Basket{
					IdUser: uuid.New(),
					Date:   time.Now(),
				}
				return basket, uuid.Nil
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basket, userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
		expectedErr bool
	}{
		{
			name: "successfully get basket by user id",
			setup: func() uuid.UUID {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				return userID
			},
			expectedErr: false,
		},
		{
			name: "fail to get basket for non-existent user",
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

func TestBasket_GetItems_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name          string
		setup         func() uuid.UUID
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get empty basket items",
			setup: func() uuid.UUID {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				return userID
			},
			expectedCount: 0,
			expectedErr:   false,
		},
		{
			name: "successfully get basket with items",
			setup: func() uuid.UUID {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				basketObj, err := fixture.service.GetById(fixture.ctx, userID)
				require.NoError(t, err)

				product1 := fixture.createTestProduct()
				product2 := fixture.createTestProduct()

				item1 := fixture.createTestBasketItem(basketObj.Id, product1)
				item2 := fixture.createTestBasketItem(basketObj.Id, product2)

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

func TestBasket_AddItem_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.BasketItem, uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successfully add item to basket",
			setup: func() (structs.BasketItem, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				basketObj, err := fixture.service.GetById(fixture.ctx, userID)
				require.NoError(t, err)

				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(basketObj.Id, productID)

				return item, userID
			},
			expectedErr: false,
		},
		{
			name: "successfully update amount when adding existing item",
			setup: func() (structs.BasketItem, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				basketObj, err := fixture.service.GetById(fixture.ctx, userID)
				require.NoError(t, err)

				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(basketObj.Id, productID)

				err = fixture.service.AddItem(fixture.ctx, item, userID)
				require.NoError(t, err)

				sameItem := fixture.createTestBasketItem(basketObj.Id, productID)
				sameItem.Amount = 3

				return sameItem, userID
			},
			expectedErr: false,
		},
		{
			name: "fail to add item to non-existent user basket",
			setup: func() (structs.BasketItem, uuid.UUID) {
				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(uuid.New(), productID)
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
		expectedErr bool
	}{
		{
			name: "successfully delete item from basket",
			setup: func() (uuid.UUID, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				basketObj, err := fixture.service.GetById(fixture.ctx, userID)
				require.NoError(t, err)

				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(basketObj.Id, productID)

				err = fixture.service.AddItem(fixture.ctx, item, userID)
				require.NoError(t, err)

				return userID, productID
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent item",
			setup: func() (uuid.UUID, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
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

func TestBasket_UpdateItemAmount_AAA(t *testing.T) {
	fixture := NewBasketTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (uuid.UUID, uuid.UUID, int)
		expectedErr bool
	}{
		{
			name: "successfully update item amount",
			setup: func() (uuid.UUID, uuid.UUID, int) {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				basketObj, err := fixture.service.GetById(fixture.ctx, userID)
				require.NoError(t, err)

				productID := fixture.createTestProduct()
				item := fixture.createTestBasketItem(basketObj.Id, productID)

				err = fixture.service.AddItem(fixture.ctx, item, userID)
				require.NoError(t, err)

				return userID, productID, 5
			},
			expectedErr: false,
		},
		{
			name: "fail to update amount for non-existent item",
			setup: func() (uuid.UUID, uuid.UUID, int) {
				userID, _, _ := fixture.createUserForTest()

				basket := structs.Basket{
					IdUser: userID,
					Date:   time.Now(),
				}
				err := fixture.service.Create(fixture.ctx, basket)
				require.NoError(t, err)

				return userID, uuid.New(), 5
			},
			expectedErr: false,
		},
		{
			name: "fail to update amount for non-existent user",
			setup: func() (uuid.UUID, uuid.UUID, int) {
				return uuid.New(), uuid.New(), 5
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, productID, amount := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
