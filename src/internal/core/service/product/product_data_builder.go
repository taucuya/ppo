package product

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type ProductBuilder struct {
	product structs.Product
}

func NewProductBuilder() *ProductBuilder {
	return &ProductBuilder{
		product: structs.Product{
			Id:          structs.GenId(),
			Name:        "Test Product",
			Description: "Test Description",
			Articule:    "TEST123",
			Category:    "electronics",
			IdBrand:     structs.GenId(),
			Price:       999.99,
			Amount:      10,
		},
	}
}

func (b *ProductBuilder) WithID(id uuid.UUID) *ProductBuilder {
	b.product.Id = id
	return b
}

func (b *ProductBuilder) WithName(name string) *ProductBuilder {
	b.product.Name = name
	return b
}

func (b *ProductBuilder) WithDescription(description string) *ProductBuilder {
	b.product.Description = description
	return b
}

func (b *ProductBuilder) WithArticule(articule string) *ProductBuilder {
	b.product.Articule = articule
	return b
}

func (b *ProductBuilder) WithCategory(category string) *ProductBuilder {
	b.product.Category = category
	return b
}

func (b *ProductBuilder) WithIdBrand(brandId uuid.UUID) *ProductBuilder {
	b.product.IdBrand = brandId
	return b
}

func (b *ProductBuilder) WithPrice(price float64) *ProductBuilder {
	b.product.Price = price
	return b
}

func (b *ProductBuilder) WithAmount(amount int) *ProductBuilder {
	b.product.Amount = amount
	return b
}

func (b *ProductBuilder) Build() structs.Product {
	return b.product
}

type TestFixture struct {
	t              *testing.T
	ctrl           *gomock.Controller
	ctx            context.Context
	productBuilder *ProductBuilder
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)

	return &TestFixture{
		t:              t,
		ctrl:           ctrl,
		ctx:            context.Background(),
		productBuilder: NewProductBuilder(),
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockProductRepository) {
	mockRepo := mock_structs.NewMockProductRepository(f.ctrl)
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
