package review_rep

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type ReviewMother struct{}

func NewReviewMother() *ReviewMother {
	return &ReviewMother{}
}

func (m *ReviewMother) ValidReview() structs.Review {
	return structs.Review{
		IdProduct: uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012"),
		IdUser:    uuid.MustParse("c3d4e5f6-a7b8-9012-cdef-345678901234"),
		Rating:    5,
		Text:      "Excellent product!",
		Date:      time.Now(),
	}
}

func (m *ReviewMother) AnotherValidReview() structs.Review {
	return structs.Review{
		Id:        uuid.MustParse("d4e5f6a7-b8c9-0123-defa-456789012345"),
		IdProduct: uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012"),
		IdUser:    uuid.MustParse("e5f6a7b8-b9c0-1234-efab-567890123456"),
		Rating:    4,
		Text:      "Very good, but could be better",
		Date:      time.Now().Add(-24 * time.Hour),
	}
}

func (m *ReviewMother) DifferentProductReview() structs.Review {
	review := m.ValidReview()
	review.IdProduct = uuid.MustParse("f6a7a8b9-c0d1-2345-fabc-678901234567")
	return review
}

func (m *ReviewMother) LowRatingReview() structs.Review {
	review := m.ValidReview()
	review.Rating = 2
	review.Text = "Not satisfied"
	return review
}

func (m *ReviewMother) EmptyTextReview() structs.Review {
	review := m.ValidReview()
	review.Text = ""
	return review
}

func (m *ReviewMother) ReviewsForProduct(productID uuid.UUID) []structs.Review {
	review1 := m.ValidReview()
	review1.IdProduct = productID

	review2 := m.AnotherValidReview()
	review2.IdProduct = productID

	return []structs.Review{review1, review2}
}

func (m *ReviewMother) EmptyReviews() []structs.Review {
	return nil
}

type TestFixture struct {
	t            *testing.T
	db           *sql.DB
	sqlxDB       *sqlx.DB
	mock         sqlmock.Sqlmock
	repo         *Repository
	ctx          context.Context
	reviewMother *ReviewMother
}

func NewTestFixture(t *testing.T) *TestFixture {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return &TestFixture{
		t:            t,
		db:           db,
		sqlxDB:       sqlxDB,
		mock:         mock,
		repo:         New(sqlxDB),
		ctx:          context.Background(),
		reviewMother: NewReviewMother(),
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
