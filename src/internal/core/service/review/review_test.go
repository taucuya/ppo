package review

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	validReview := fixture.reviewMother.ValidReview()

	tests := []struct {
		name        string
		review      structs.Review
		setupMocks  func(*mock_structs.MockReviewRepository, structs.Review)
		expectedErr error
	}{
		{
			name:   "successful creation",
			review: validReview,
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, review structs.Review) {
				mockRepo.EXPECT().Create(fixture.ctx, review).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "repository error",
			review: validReview,
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, review structs.Review) {
				mockRepo.EXPECT().Create(fixture.ctx, review).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, tt.review)

			err := service.Create(fixture.ctx, tt.review)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	validReview := fixture.reviewMother.ValidReview()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockReviewRepository, structs.Review)
		expectedRet structs.Review
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, review structs.Review) {
				mockRepo.EXPECT().GetById(fixture.ctx, review.Id).Return(review, nil)
			},
			expectedRet: validReview,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, review structs.Review) {
				mockRepo.EXPECT().GetById(fixture.ctx, review.Id).Return(structs.Review{}, errTest)
			},
			expectedRet: structs.Review{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, validReview)

			ret, err := service.GetById(fixture.ctx, validReview.Id)

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

func TestDelete_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	validReview := fixture.reviewMother.ValidReview()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockReviewRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, reviewID uuid.UUID) {
				mockRepo.EXPECT().Delete(fixture.ctx, reviewID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, reviewID uuid.UUID) {
				mockRepo.EXPECT().Delete(fixture.ctx, reviewID).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, validReview.Id)

			err := service.Delete(fixture.ctx, validReview.Id)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestReviewsForProduct_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	productID := fixture.reviewMother.ValidReview().IdProduct
	reviewsForProduct := fixture.reviewMother.ReviewsForProduct(productID)
	emptyReviews := fixture.reviewMother.EmptyReviews()

	tests := []struct {
		name        string
		productID   uuid.UUID
		setupMocks  func(*mock_structs.MockReviewRepository, uuid.UUID, []structs.Review)
		expectedRet []structs.Review
		expectedErr error
	}{
		{
			name:      "successful get reviews for product",
			productID: productID,
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, productID uuid.UUID, reviews []structs.Review) {
				mockRepo.EXPECT().ReviewsForProduct(fixture.ctx, productID).Return(reviews, nil)
			},
			expectedRet: reviewsForProduct,
			expectedErr: nil,
		},
		{
			name:      "no reviews for product",
			productID: uuid.MustParse("f6a7b8c9-c0d1-2345-fabc-678901234567"),
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, productID uuid.UUID, reviews []structs.Review) {
				mockRepo.EXPECT().ReviewsForProduct(fixture.ctx, productID).Return(emptyReviews, nil)
			},
			expectedRet: emptyReviews,
			expectedErr: nil,
		},
		{
			name:      "repository error",
			productID: productID,
			setupMocks: func(mockRepo *mock_structs.MockReviewRepository, productID uuid.UUID, reviews []structs.Review) {
				mockRepo.EXPECT().ReviewsForProduct(fixture.ctx, productID).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, tt.productID, tt.expectedRet)

			ret, err := service.ReviewsForProduct(fixture.ctx, tt.productID)

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
