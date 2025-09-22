package favourites

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
	t               *testing.T
	ctrl            *gomock.Controller
	ctx             context.Context
	favourites      structs.Favourites
	favouritesItems []structs.FavouritesItem
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)
	userId := structs.GenId()
	favouritesId := structs.GenId()
	productId1 := structs.GenId()
	productId2 := structs.GenId()

	return &TestFixture{
		t:    t,
		ctrl: ctrl,
		ctx:  context.Background(),
		favourites: structs.Favourites{
			Id:     favouritesId,
			IdUser: userId,
		},
		favouritesItems: []structs.FavouritesItem{
			{
				Id:           structs.GenId(),
				IdProduct:    productId1,
				IdFavourites: favouritesId,
			},
			{
				Id:           structs.GenId(),
				IdProduct:    productId2,
				IdFavourites: favouritesId,
			},
		},
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockFavouritesRepository) {
	mockRepo := mock_structs.NewMockFavouritesRepository(f.ctrl)
	service := New(mockRepo)
	return service, mockRepo
}

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
