package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type UserBuilder struct {
	user structs.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: structs.User{
			Id:       uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
			Name:     "Test User",
			Mail:     "test@example.com",
			Phone:    "+1234567890",
			Password: "hashed_password_123",
			Role:     "customer",
		},
	}
}

func (b *UserBuilder) WithID(id uuid.UUID) *UserBuilder {
	b.user.Id = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

func (b *UserBuilder) WithMail(mail string) *UserBuilder {
	b.user.Mail = mail
	return b
}

func (b *UserBuilder) WithPhone(phone string) *UserBuilder {
	b.user.Phone = phone
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.user.Password = password
	return b
}

func (b *UserBuilder) WithRole(role string) *UserBuilder {
	b.user.Role = role
	return b
}

func (b *UserBuilder) Build() structs.User {
	return b.user
}

type BasketBuilder struct {
	basket structs.Basket
}

func NewBasketBuilder() *BasketBuilder {
	return &BasketBuilder{
		basket: structs.Basket{
			Id:   uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012"),
			Date: time.Now(),
		},
	}
}

func (b *BasketBuilder) WithUserID(userID uuid.UUID) *BasketBuilder {
	b.basket.IdUser = userID
	return b
}

func (b *BasketBuilder) WithDate(date time.Time) *BasketBuilder {
	b.basket.Date = date
	return b
}

func (b *BasketBuilder) Build() structs.Basket {
	return b.basket
}

type FavouritesBuilder struct {
	favourites structs.Favourites
}

func NewFavouritesBuilder() *FavouritesBuilder {
	return &FavouritesBuilder{
		favourites: structs.Favourites{
			Id: uuid.MustParse("c3d4e5f6-a7b8-9012-cdef-345678901234"),
		},
	}
}

func (b *FavouritesBuilder) WithUserID(userID uuid.UUID) *FavouritesBuilder {
	b.favourites.IdUser = userID
	return b
}

func (b *FavouritesBuilder) Build() structs.Favourites {
	return b.favourites
}

type TestFixture struct {
	t             *testing.T
	ctrl          *gomock.Controller
	ctx           context.Context
	userBuilder   *UserBuilder
	basketBuilder *BasketBuilder
	favBuilder    *FavouritesBuilder
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)

	return &TestFixture{
		t:             t,
		ctrl:          ctrl,
		ctx:           context.Background(),
		userBuilder:   NewUserBuilder(),
		basketBuilder: NewBasketBuilder(),
		favBuilder:    NewFavouritesBuilder(),
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockUserRepository, *mock_structs.MockUsrBasket, *mock_structs.MockUsrFavourites) {
	mockRepo := mock_structs.NewMockUserRepository(f.ctrl)
	mockBasket := mock_structs.NewMockUsrBasket(f.ctrl)
	mockFav := mock_structs.NewMockUsrFavourites(f.ctrl)

	service := New(mockRepo, mockBasket, mockFav)
	return service, mockRepo, mockBasket, mockFav
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
