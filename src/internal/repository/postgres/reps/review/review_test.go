package review_rep

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate(t *testing.T) {
	fixture := NewTestFixture(t)

	validReview := fixture.reviewMother.ValidReview()

	tests := []struct {
		name        string
		setupMocks  func(structs.Review)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectExec(`insert into review \(id_product, id_user, rating, r_text, date\)`).
					WithArgs(
						review.IdProduct,
						review.IdUser,
						review.Rating,
						review.Text,
						review.Date,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectExec(`insert into review \(id_product, id_user, rating, r_text, date\)`).
					WithArgs(
						review.IdProduct,
						review.IdUser,
						review.Rating,
						review.Text,
						review.Date,
					).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(validReview)

			err := fixture.repo.Create(fixture.ctx, validReview)
			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetById(t *testing.T) {
	fixture := NewTestFixture(t)

	validReview := fixture.reviewMother.ValidReview()

	tests := []struct {
		name        string
		setupMocks  func(structs.Review)
		expectedRet structs.Review
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMocks: func(review structs.Review) {
				rows := sqlmock.NewRows([]string{"id_product", "id_user", "rating", "r_text", "date"}).
					AddRow(review.IdProduct, review.IdUser, review.Rating, review.Text, review.Date)
				fixture.mock.ExpectQuery(`select \* from review where id = \$1`).
					WithArgs(review.Id).
					WillReturnRows(rows)
			},
			expectedRet: validReview,
			expectedErr: nil,
		},
		{
			name: "review not found",
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectQuery(`select \* from review where id = \$1`).
					WithArgs(review.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedRet: structs.Review{},
			expectedErr: errors.New("failed to get review: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error",
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectQuery(`select \* from review where id = \$1`).
					WithArgs(review.Id).
					WillReturnError(errTest)
			},
			expectedRet: structs.Review{},
			expectedErr: errors.New("failed to get review: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(validReview)

			ret, err := fixture.repo.GetById(fixture.ctx, validReview.Id)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDelete(t *testing.T) {
	fixture := NewTestFixture(t)

	validReview := fixture.reviewMother.ValidReview()

	tests := []struct {
		name        string
		setupMocks  func(uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(reviewID uuid.UUID) {
				fixture.mock.ExpectExec(`delete from review where id = \$1`).
					WithArgs(reviewID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "review not found",
			setupMocks: func(reviewID uuid.UUID) {
				fixture.mock.ExpectExec(`delete from review where id = \$1`).
					WithArgs(reviewID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: errors.New("review with id " + validReview.Id.String() + " not found"),
		},
		{
			name: "database error",
			setupMocks: func(reviewID uuid.UUID) {
				fixture.mock.ExpectExec(`delete from review where id = \$1`).
					WithArgs(reviewID).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(validReview.Id)

			err := fixture.repo.Delete(fixture.ctx, validReview.Id)
			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestReviewsForProduct(t *testing.T) {
	fixture := NewTestFixture(t)

	productID := fixture.reviewMother.ValidReview().IdProduct
	reviewsForProduct := fixture.reviewMother.ReviewsForProduct(productID)
	emptyReviews := fixture.reviewMother.EmptyReviews()

	tests := []struct {
		name        string
		productID   uuid.UUID
		setupMocks  func(uuid.UUID, []structs.Review)
		expectedRet []structs.Review
		expectedErr error
	}{
		{
			name:      "successful get reviews for product",
			productID: productID,
			setupMocks: func(productID uuid.UUID, reviews []structs.Review) {
				rows := sqlmock.NewRows([]string{"id", "id_product", "id_user", "rating", "r_text", "date"})
				for _, review := range reviews {
					rows.AddRow(review.Id, review.IdProduct, review.IdUser, review.Rating, review.Text, review.Date)
				}
				fixture.mock.ExpectQuery(`select \* from review where id_product = \$1 order by date desc`).
					WithArgs(productID).
					WillReturnRows(rows)
			},
			expectedRet: reviewsForProduct,
			expectedErr: nil,
		},
		{
			name:      "no reviews for product",
			productID: uuid.MustParse("f6a7b8c9-c0d1-2345-fabc-678901234567"),
			setupMocks: func(productID uuid.UUID, reviews []structs.Review) {
				rows := sqlmock.NewRows([]string{"id", "id_product", "id_user", "rating", "r_text", "date"})
				fixture.mock.ExpectQuery(`select \* from review where id_product = \$1 order by date desc`).
					WithArgs(productID).
					WillReturnRows(rows)
			},
			expectedRet: emptyReviews,
			expectedErr: nil,
		},
		{
			name:      "database error",
			productID: productID,
			setupMocks: func(productID uuid.UUID, reviews []structs.Review) {
				fixture.mock.ExpectQuery(`select \* from review where id_product = \$1 order by date desc`).
					WithArgs(productID).
					WillReturnError(errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(tt.productID, tt.expectedRet)

			ret, err := fixture.repo.ReviewsForProduct(fixture.ctx, tt.productID)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
}

func TestCreateWithDifferentReviews(t *testing.T) {
	fixture := NewTestFixture(t)
	defer fixture.Cleanup()

	tests := []struct {
		name        string
		review      structs.Review
		setupMocks  func(structs.Review)
		expectedErr error
	}{
		{
			name:   "create low rating review",
			review: fixture.reviewMother.LowRatingReview(),
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectExec(`insert into review`).
					WithArgs(
						review.IdProduct,
						review.IdUser,
						review.Rating,
						review.Text,
						review.Date,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:   "create empty text review",
			review: fixture.reviewMother.EmptyTextReview(),
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectExec(`insert into review`).
					WithArgs(
						review.IdProduct,
						review.IdUser,
						review.Rating,
						review.Text,
						review.Date,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:   "create review for different product",
			review: fixture.reviewMother.DifferentProductReview(),
			setupMocks: func(review structs.Review) {
				fixture.mock.ExpectExec(`insert into review`).
					WithArgs(
						review.IdProduct,
						review.IdUser,
						review.Rating,
						review.Text,
						review.Date,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(tt.review)

			err := fixture.repo.Create(fixture.ctx, tt.review)
			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
}
