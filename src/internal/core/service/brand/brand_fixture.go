package brand

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type TestFixture struct {
	t     *testing.T
	ctrl  *gomock.Controller
	ctx   context.Context
	brand structs.Brand
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	brandId := structs.GenId()
	name := "Test brand"
	description := "Test description"
	priceCategory := "Test price category"

	return &TestFixture{
		t:    t,
		ctrl: ctrl,
		ctx:  ctx,
		brand: structs.Brand{
			Id:            brandId,
			Name:          name,
			Description:   description,
			PriceCategory: priceCategory,
		},
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockBrandRepository) {
	mockRepo := mock_structs.NewMockBrandRepository(f.ctrl)

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
