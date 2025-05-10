package review

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"
// 	"github.com/taucuya/ppo/internal/core/mock_structs"
// 	"github.com/taucuya/ppo/internal/core/structs"
// )

// var testError = errors.New("test error")

// func TestCreate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockReviewRepository(ctrl)
// 	service := New(mockRepo)

// 	testReview := structs.Review{
// 		Id:        structs.GenId(),
// 		IdProduct: structs.GenId(),
// 		IdUser:    structs.GenId(),
// 		Rating:    5,
// 		Text:      "Great product!",
// 		Date:      time.Time{},
// 	// }

// 	t.Run("successful creation", func(t *testing.T) {
// 		mockRepo.EXPECT().Create(gomock.Any(), testReview).Return(nil).Times(1)

// 		err := service.Create(context.Background(), testReview)
// 		if err != nil {
// 			t.Errorf("Create() unexpected error = %v", err)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().Create(gomock.Any(), testReview).Return(testError).Times(1)

// 		err := service.Create(context.Background(), testReview)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("Create() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestGetById(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockReviewRepository(ctrl)
// 	service := New(mockRepo)
// 	testID := structs.GenId()
// 	testReview := structs.Review{
// 		Id: testID, IdProduct: structs.GenId(), IdUser: structs.GenId(), Rating: 5, Text: "Nice", Date: time.Time{},
// 	}

// 	t.Run("successful get", func(t *testing.T) {
// 		mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(testReview, nil).Times(1)

// 		got, err := service.GetById(context.Background(), testID)
// 		if err != nil {
// 			t.Errorf("GetById() unexpected error = %v", err)
// 		}
// 		if got != testReview {
// 			t.Errorf("GetById() = %v, want %v", got, testReview)
// 		}
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(structs.Review{}, testError).Times(1)

// 		_, err := service.GetById(context.Background(), testID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("GetById() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestDelete(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockReviewRepository(ctrl)
// 	service := New(mockRepo)
// 	reviewID := structs.GenId()

// 	t.Run("successful delete", func(t *testing.T) {
// 		mockRepo.EXPECT().Delete(gomock.Any(), reviewID).Return(nil).Times(1)

// 		err := service.Delete(context.Background(), reviewID)
// 		if err != nil {
// 			t.Errorf("Delete() unexpected error = %v", err)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().Delete(gomock.Any(), reviewID).Return(testError).Times(1)

// 		err := service.Delete(context.Background(), reviewID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("Delete() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestReviewsForProduct(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockReviewRepository(ctrl)
// 	service := New(mockRepo)

// 	productID := structs.GenId()
// 	testReviews := []structs.Review{
// 		{Id: structs.GenId(), IdProduct: productID, IdUser: structs.GenId(), Rating: 4, Text: "Good", Date: time.Time{}},
// 		{Id: structs.GenId(), IdProduct: productID, IdUser: structs.GenId(), Rating: 5, Text: "Excellent", Date: time.Time{}},
// 	}

// 	t.Run("successful fetch", func(t *testing.T) {
// 		mockRepo.EXPECT().ReviewsForProduct(gomock.Any(), productID).Return(testReviews, nil).Times(1)

// 		got, err := service.ReviewsForProduct(context.Background(), productID)
// 		if err != nil {
// 			t.Errorf("ReviewsForProduct() unexpected error = %v", err)
// 		}
// 		if len(got) != len(testReviews) {
// 			t.Errorf("ReviewsForProduct() = %v, want %v", got, testReviews)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().ReviewsForProduct(gomock.Any(), productID).Return(nil, testError).Times(1)

// 		_, err := service.ReviewsForProduct(context.Background(), productID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("ReviewsForProduct() error = %v, want %v", err, testError)
// 		}
// 	})
// }
