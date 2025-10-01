package product_rep

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type ProductBuilder struct {
	product structs.Product
}

func NewProductBuilder() *ProductBuilder {
	return &ProductBuilder{
		product: structs.Product{
			Id:          uuid.New(),
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			Category:    "electronics",
			Amount:      10,
			IdBrand:     uuid.New(),
			PicLink:     "http://example.com/pic.jpg",
			Articule:    "ART123",
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

func (b *ProductBuilder) WithPrice(price float64) *ProductBuilder {
	b.product.Price = price
	return b
}

func (b *ProductBuilder) WithCategory(category string) *ProductBuilder {
	b.product.Category = category
	return b
}

func (b *ProductBuilder) WithAmount(amount int) *ProductBuilder {
	b.product.Amount = amount
	return b
}

func (b *ProductBuilder) WithIdBrand(idBrand uuid.UUID) *ProductBuilder {
	b.product.IdBrand = idBrand
	return b
}

func (b *ProductBuilder) WithPicLink(picLink string) *ProductBuilder {
	b.product.PicLink = picLink
	return b
}

func (b *ProductBuilder) WithArticule(articule string) *ProductBuilder {
	b.product.Articule = articule
	return b
}

func (b *ProductBuilder) Build() structs.Product {
	return b.product
}

type TestFixture struct {
	t              *testing.T
	db             *sql.DB
	sqlxDB         *sqlx.DB
	mock           sqlmock.Sqlmock
	repo           *Repository
	ctx            context.Context
	productBuilder *ProductBuilder
}

func NewTestFixture(t *testing.T) *TestFixture {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return &TestFixture{
		t:              t,
		db:             db,
		sqlxDB:         sqlxDB,
		mock:           mock,
		repo:           New(sqlxDB),
		ctx:            context.Background(),
		productBuilder: NewProductBuilder(),
	}
}

func (f *TestFixture) Cleanup() {
	f.db.Close()
}

func (f *TestFixture) AssertError(actual, expected error) {
	if expected == nil {
		assert.NoError(f.t, actual)
	} else {
		assert.EqualError(f.t, actual, expected.Error())
	}
}
