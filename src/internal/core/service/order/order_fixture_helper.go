package order

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type TestFixture struct {
	t          *testing.T
	ctrl       *gomock.Controller
	ctx        context.Context
	order      structs.Order
	orderItems []structs.OrderItem
	orders     []structs.Order
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)
	orderId := structs.GenId()
	userId := structs.GenId()
	productId1 := structs.GenId()
	productId2 := structs.GenId()

	return &TestFixture{
		t:    t,
		ctrl: ctrl,
		ctx:  context.Background(),
		order: structs.Order{
			Id:     orderId,
			IdUser: userId,
			Status: "pending",
		},
		orderItems: []structs.OrderItem{
			{
				Id:        structs.GenId(),
				IdProduct: productId1,
				IdOrder:   orderId,
				Amount:    2,
			},
			{
				Id:        structs.GenId(),
				IdProduct: productId2,
				IdOrder:   orderId,
				Amount:    1,
			},
		},
		orders: []structs.Order{
			{
				Id:     structs.GenId(),
				IdUser: structs.GenId(),
				Status: "free",
			},
			{
				Id:     structs.GenId(),
				IdUser: structs.GenId(),
				Status: "free",
			},
		},
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

// func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockOrderRepository) {
// mockRepo := mock_structs.NewMockOrderRepository(f.ctrl)
// service := New(mockRepo)
// return service, mockRepo
// }

func (f *TestFixture) AssertError(err error, expectedErr error) {
	if expectedErr != nil {
		if err == nil {
			f.t.Error("Expected error, got nil")
			return
		}
		if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
			f.t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	} else if err != nil {
		f.t.Errorf("Unexpected error: %v", err)
	}
}
