package review

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

type ReviewMother struct{}

func NewReviewMother() *ReviewMother {
	return &ReviewMother{}
}

func (m *ReviewMother) ValidReview() structs.Review {

	return structs.Review{
		Id:        uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
		IdProduct: uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012"),
		IdUser:    uuid.MustParse("c3d4e5f6-a7b8-9012-cdef-345678901234"),
		Rating:    5,
		Text:      "Excellent product!",
	}
}

func (m *ReviewMother) AnotherValidReview() structs.Review {
	return structs.Review{
		Id:        uuid.MustParse("d4e5f6a7-b8c9-0123-defa-456789012345"),
		IdProduct: uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012"),
		IdUser:    uuid.MustParse("e5f6a7b8-b9c0-1234-efab-567890123456"),
		Rating:    4,
		Text:      "Very good, but could be better",
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
	return []structs.Review{}
}

type TestFixture struct {
	t            *testing.T
	ctrl         *gomock.Controller
	ctx          context.Context
	reviewMother *ReviewMother
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)

	return &TestFixture{
		t:            t,
		ctrl:         ctrl,
		ctx:          context.Background(),
		reviewMother: NewReviewMother(),
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockReviewRepository) {
	mockRepo := mock_structs.NewMockReviewRepository(f.ctrl)
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
