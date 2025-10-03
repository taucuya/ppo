package order

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
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().Create(fixture.ctx, fixture.order).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().Create(fixture.ctx, fixture.order).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.Create(fixture.ctx, fixture.order)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedRet structs.Order
		expectedErr error
	}{
		{
			name: "successful get",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.order.Id).Return(fixture.order, nil)
			},
			expectedRet: fixture.order,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.order.Id).Return(structs.Order{}, errTest)
			},
			expectedRet: structs.Order{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetById(fixture.ctx, fixture.order.Id)

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

func TestGetItems_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedRet []structs.OrderItem
		expectedErr error
	}{
		{
			name: "successful get items",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetItems(fixture.ctx, fixture.order.Id).Return(fixture.orderItems, nil)
			},
			expectedRet: fixture.orderItems,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetItems(fixture.ctx, fixture.order.Id).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetItems(fixture.ctx, fixture.order.Id)

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

func TestGetFreeOrders_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedRet []structs.Order
		expectedErr error
	}{
		{
			name: "successful get free orders",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetFreeOrders(fixture.ctx).Return(fixture.orders, nil)
			},
			expectedRet: fixture.orders,
			expectedErr: nil,
		},
		{
			name: "empty list",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetFreeOrders(fixture.ctx).Return([]structs.Order{}, nil)
			},
			expectedRet: []structs.Order{},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetFreeOrders(fixture.ctx).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetFreeOrders(fixture.ctx)

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

func TestGetOrdersByUser_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedRet []structs.Order
		expectedErr error
	}{
		{
			name: "successful get free orders",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetOrdersByUser(fixture.ctx, fixture.order.IdUser).Return(fixture.orders, nil)
			},
			expectedRet: fixture.orders,
			expectedErr: nil,
		},
		{
			name: "empty list",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetOrdersByUser(fixture.ctx, fixture.order.IdUser).Return([]structs.Order{}, nil)
			},
			expectedRet: []structs.Order{},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetOrdersByUser(fixture.ctx, fixture.order.IdUser).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetOrdersByUser(fixture.ctx, fixture.order.IdUser)

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

func TestGetStatus_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedRet string
		expectedErr error
	}{
		{
			name: "successful get status",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("pending", nil)
			},
			expectedRet: "pending",
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("", errTest)
			},
			expectedRet: "",
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetStatus(fixture.ctx, fixture.order.Id)

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

func TestChangeOrderStatus_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		status      string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedErr error
	}{
		{
			name:   "successful status change",
			status: "completed",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("pending", nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.order.Id).Return(fixture.order, nil)
				mockRepo.EXPECT().UpdateStatus(fixture.ctx, fixture.order.Id, "completed").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "same status - no update needed",
			status: "pending",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("pending", nil)
			},
			expectedErr: nil,
		},
		{
			name:   "error getting status",
			status: "completed",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("", errTest)
			},
			expectedErr: errTest,
		},
		{
			name:   "error getting order",
			status: "completed",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("pending", nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.order.Id).Return(structs.Order{}, errTest)
			},
			expectedErr: errTest,
		},
		{
			name:   "error updating status",
			status: "completed",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().GetStatus(fixture.ctx, fixture.order.Id).Return("pending", nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.order.Id).Return(fixture.order, nil)
				mockRepo.EXPECT().UpdateStatus(fixture.ctx, fixture.order.Id, "completed").Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.ChangeOrderStatus(fixture.ctx, fixture.order.Id, tt.status)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestDelete_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockOrderRepository)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().Delete(fixture.ctx, fixture.order.Id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockOrderRepository) {
				mockRepo.EXPECT().Delete(fixture.ctx, fixture.order.Id).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.Delete(fixture.ctx, fixture.order.Id)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}
