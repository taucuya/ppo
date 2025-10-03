package basket

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type TestFixture struct {
	t           *testing.T
	ctrl        *gomock.Controller
	ctx         context.Context
	basket      structs.Basket
	basketItems []structs.BasketItem
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	basketID := structs.GenId()
	idUser := structs.GenId()
	date := time.Now()
	idItem1 := structs.GenId()
	idItem2 := structs.GenId()
	idProduct1 := structs.GenId()
	idProduct2 := structs.GenId()

	return &TestFixture{
		t:    t,
		ctrl: ctrl,
		ctx:  ctx,
		basket: structs.Basket{
			Id:     basketID,
			IdUser: idUser,
			Date:   date,
		},
		basketItems: []structs.BasketItem{
			{
				Id:        idItem1,
				IdProduct: idProduct1,
				IdBasket:  basketID,
				Amount:    2,
			},
			{
				Id:        idItem2,
				IdProduct: idProduct2,
				IdBasket:  basketID,
				Amount:    1,
			},
		},
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockBasketRepository) {
	mockRepo := mock_structs.NewMockBasketRepository(f.ctrl)

	service := New(mockRepo)
	return service, mockRepo
}

func (f *TestFixture) AssertError(err error, expectedErr error) {
	if expectedErr != nil {
		if err == nil {
			f.t.Errorf("Expected error %v, got nil", expectedErr)
			return
		} else if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
			f.t.Errorf("Expected  error %v, got %v", expectedErr, err)
		}

	} else if err != nil {
		f.t.Errorf("Expected error nil, got %v", err)
		return
	}
}
