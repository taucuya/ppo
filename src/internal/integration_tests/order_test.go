package integrationtests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/order"
	"github.com/taucuya/ppo/internal/core/structs"
	order_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/order"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

type OrderTestFixture struct {
	t         *testing.T
	ctx       context.Context
	service   *order.Service
	orderRepo *order_rep.Repository
	userRepo  *user_rep.Repository
}

func NewOrderTestFixture(t *testing.T) *OrderTestFixture {
	orderRepo := order_rep.New(db)
	userRepo := user_rep.New(db)

	service := order.New(orderRepo)

	return &OrderTestFixture{
		t:         t,
		ctx:       context.Background(),
		service:   service,
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (f *OrderTestFixture) createTestUser() structs.User {
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

func (f *OrderTestFixture) createTestOrder(userID uuid.UUID) structs.Order {
	return structs.Order{
		IdUser:  userID,
		Address: "123 Test St",
		Status:  "непринятый",
	}
}

func (f *OrderTestFixture) setupUserWithOrder() (uuid.UUID, uuid.UUID) {
	testUser := f.createTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)

	basketID := uuid.New()
	_, err = db.Exec("INSERT INTO basket (id, id_user, date) VALUES ($1, $2, $3)",
		basketID, userID, time.Now())
	require.NoError(f.t, err)

	brandID := uuid.New()
	_, err = db.Exec("INSERT INTO brand (id, name, description, price_category) VALUES ($1, $2, $3, $4)",
		brandID, "Test Brand", "Test Description", "premium")
	require.NoError(f.t, err)

	product1 := uuid.New()
	product2 := uuid.New()
	_, err = db.Exec(`INSERT INTO product (id, name, description, price, id_brand, amount, art) 
		VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)`,
		product1, "Product 1", "Description 1", 1000.0, brandID, 20, "ART-001",
		product2, "Product 2", "Description 2", 2000.0, brandID, 15, "ART-002")
	require.NoError(f.t, err)

	_, err = db.Exec(`INSERT INTO basket_item (id_product, id_basket, amount) 
		VALUES ($1, $2, $3), ($4, $5, $6)`,
		product1, basketID, 2,
		product2, basketID, 1)
	require.NoError(f.t, err)

	testOrder := f.createTestOrder(userID)
	err = f.orderRepo.Create(f.ctx, testOrder)
	require.NoError(f.t, err)

	var orders []struct {
		ID uuid.UUID `db:"id"`
	}
	err = db.SelectContext(f.ctx, &orders, `SELECT id FROM "order" WHERE id_user = $1`, userID)
	require.NoError(f.t, err)
	require.Len(f.t, orders, 1)
	return userID, orders[0].ID
}

func TestOrder_Create_AAA(t *testing.T) {
	fixture := NewOrderTestFixture(t)

	tests := []struct {
		name        string
		setup       func() structs.Order
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful order creation",
			setup: func() structs.Order {
				truncateTables(t)
				testUser := fixture.createTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)

				basketID := uuid.New()
				_, err = db.Exec("INSERT INTO basket (id, id_user, date) VALUES ($1, $2, $3)",
					basketID, userID, time.Now())
				require.NoError(t, err)

				brandID := uuid.New()
				_, err = db.Exec("INSERT INTO brand (id, name, description, price_category) VALUES ($1, $2, $3, $4)",
					brandID, "Test Brand", "Test Description", "premium")
				require.NoError(t, err)

				product1 := uuid.New()
				_, err = db.Exec("INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					product1, "Test Product", "Test Description", 1000.0, brandID, 10, "TEST-ART-123")
				require.NoError(t, err)

				_, err = db.Exec("INSERT INTO basket_item (id_product, id_basket, amount) VALUES ($1, $2, $3)",
					product1, basketID, 2)
				require.NoError(t, err)

				return fixture.createTestOrder(userID)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to create order for non-existent user",
			setup: func() structs.Order {
				truncateTables(t)
				return fixture.createTestOrder(uuid.New())
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, order)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrder_GetItems_AAA(t *testing.T) {
	fixture := NewOrderTestFixture(t)

	tests := []struct {
		name          string
		setup         func() uuid.UUID
		cleanup       func()
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get order items",
			setup: func() uuid.UUID {
				truncateTables(t)
				_, orderID := fixture.setupUserWithOrder()
				return orderID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "get empty items for non-existent order",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID := tt.setup()
			defer tt.cleanup()

			items, err := fixture.service.GetItems(fixture.ctx, orderID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, items, tt.expectedCount)
			}
		})
	}
}

func TestOrder_GetStatus_AAA(t *testing.T) {
	fixture := NewOrderTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expected    string
		expectedErr bool
	}{
		{
			name: "successfully get order status",
			setup: func() uuid.UUID {
				truncateTables(t)
				_, orderID := fixture.setupUserWithOrder()
				return orderID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expected:    "непринятый",
			expectedErr: false,
		},
		{
			name: "fail to get status for non-existent order",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID := tt.setup()
			defer tt.cleanup()

			status, err := fixture.service.GetStatus(fixture.ctx, orderID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, status)
			}
		})
	}
}

func TestOrder_ChangeOrderStatus_AAA(t *testing.T) {
	fixture := NewOrderTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (uuid.UUID, string)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully change order status",
			setup: func() (uuid.UUID, string) {
				truncateTables(t)
				_, orderID := fixture.setupUserWithOrder()
				return orderID, "принятый"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "no change when status is the same",
			setup: func() (uuid.UUID, string) {
				truncateTables(t)
				_, orderID := fixture.setupUserWithOrder()
				return orderID, "непринятый"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to change status for non-existent order",
			setup: func() (uuid.UUID, string) {
				truncateTables(t)
				return uuid.New(), "принятый"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID, status := tt.setup()
			defer tt.cleanup()

			err := fixture.service.ChangeOrderStatus(fixture.ctx, orderID, status)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrder_GetFreeOrders_AAA(t *testing.T) {
	fixture := NewOrderTestFixture(t)

	tests := []struct {
		name          string
		setup         func()
		cleanup       func()
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "get free orders",
			setup: func() {
				truncateTables(t)

				user1 := fixture.createTestUser()
				user1.Mail = "test1@example.com"
				user1.Phone = "89016475841"
				user1ID, err := fixture.userRepo.Create(fixture.ctx, user1)
				require.NoError(t, err)

				basketID1 := uuid.New()
				_, err = db.Exec("INSERT INTO basket (id, id_user, date) VALUES ($1, $2, $3)",
					basketID1, user1ID, time.Now())
				require.NoError(t, err)

				brandID := uuid.New()
				_, err = db.Exec("INSERT INTO brand (id, name, description, price_category) VALUES ($1, $2, $3, $4)",
					brandID, "Test Brand", "Test Description", "premium")
				require.NoError(t, err)

				product1 := uuid.New()
				_, err = db.Exec("INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					product1, "Product 1", "Description 1", 1000.0, brandID, 20, "ART-001")
				require.NoError(t, err)

				_, err = db.Exec("INSERT INTO basket_item (id_product, id_basket, amount) VALUES ($1, $2, $3)",
					product1, basketID1, 2)
				require.NoError(t, err)

				orderID1 := uuid.New()
				_, err = db.Exec(`INSERT INTO "order" (id, id_user, address, status, price) VALUES ($1, $2, $3, $4, $5)`,
					orderID1, user1ID, "123 Test St", "непринятый", 2000.0)
				require.NoError(t, err)

				user2 := fixture.createTestUser()
				user2.Mail = "test2@example.com"
				user2.Phone = "89016475842"
				user2ID, err := fixture.userRepo.Create(fixture.ctx, user2)
				require.NoError(t, err)

				basketID2 := uuid.New()
				_, err = db.Exec("INSERT INTO basket (id, id_user, date) VALUES ($1, $2, $3)",
					basketID2, user2ID, time.Now())
				require.NoError(t, err)

				product2 := uuid.New()
				_, err = db.Exec("INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					product2, "Product 2", "Description 2", 1500.0, brandID, 15, "ART-002")
				require.NoError(t, err)

				_, err = db.Exec("INSERT INTO basket_item (id_product, id_basket, amount) VALUES ($1, $2, $3)",
					product2, basketID2, 1)
				require.NoError(t, err)

				orderID2 := uuid.New()
				_, err = db.Exec(`INSERT INTO "order" (id, id_user, address, status, price) VALUES ($1, $2, $3, $4, $5)`,
					orderID2, user2ID, "456 Another St", "непринятый", 1500.0)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "get empty free orders list",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			orders, err := fixture.service.GetFreeOrders(fixture.ctx)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, orders, tt.expectedCount)
			}
		})
	}
}

func TestOrder_Delete_AAA(t *testing.T) {
	fixture := NewOrderTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully delete order",
			setup: func() uuid.UUID {
				truncateTables(t)
				_, orderID := fixture.setupUserWithOrder()
				return orderID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent order",
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
			orderID := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Delete(fixture.ctx, orderID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := fixture.service.GetById(fixture.ctx, orderID)
				require.Error(t, err)
			}
		})
	}
}
