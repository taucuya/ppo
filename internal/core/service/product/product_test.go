package product

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/taucuya/ppo/internal/core/mock_structs"
// 	"github.com/taucuya/ppo/internal/core/structs"
// )

// var testError = errors.New("test error")

// func TestCreate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockProductRepository(ctrl)
// 	service := New(mockRepo)

// 	testProduct := structs.Product{
// 		Id:          structs.GenId(),
// 		Name:        "Test Product",
// 		Description: "Test Description",
// 		Price:       99.99,
// 		Category:    "Test Category",
// 		Amount:      10,
// 		IdBrand:     structs.GenId(),
// 		PicLink:     "http://example.com/image.jpg",
// 	}

// 	t.Run("successful creation", func(t *testing.T) {
// 		mockRepo.EXPECT().
// 			Create(gomock.Any(), testProduct).
// 			Return(nil).
// 			Times(1)

// 		err := service.Create(context.Background(), testProduct)
// 		if err != nil {
// 			t.Errorf("Create() unexpected error = %v", err)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().
// 			Create(gomock.Any(), testProduct).
// 			Return(testError).
// 			Times(1)

// 		err := service.Create(context.Background(), testProduct)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("Create() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestGetById(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockProductRepository(ctrl)
// 	service := New(mockRepo)

// 	testID := structs.GenId()
// 	testProduct := structs.Product{
// 		Id:          testID,
// 		Name:        "Test Product",
// 		Description: "Test Description",
// 		Price:       99.99,
// 		Category:    "Test Category",
// 		Amount:      10,
// 		IdBrand:     structs.GenId(),
// 		PicLink:     "http://example.com/image.jpg",
// 	}

// 	t.Run("successful get", func(t *testing.T) {
// 		mockRepo.EXPECT().
// 			GetById(gomock.Any(), testID).
// 			Return(testProduct, nil).
// 			Times(1)

// 		got, err := service.GetById(context.Background(), testID)
// 		if err != nil {
// 			t.Errorf("GetById() unexpected error = %v", err)
// 		}
// 		if got != testProduct {
// 			t.Errorf("GetById() = %v, want %v", got, testProduct)
// 		}
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		mockRepo.EXPECT().
// 			GetById(gomock.Any(), testID).
// 			Return(structs.Product{}, testError).
// 			Times(1)

// 		_, err := service.GetById(context.Background(), testID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("GetById() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestDelete(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockProductRepository(ctrl)
// 	service := New(mockRepo)

// 	productID := structs.GenId()

// 	t.Run("successful delete", func(t *testing.T) {
// 		mockRepo.EXPECT().
// 			Delete(gomock.Any(), productID).
// 			Return(nil).
// 			Times(1)

// 		err := service.Delete(context.Background(), productID)
// 		if err != nil {
// 			t.Errorf("Delete() unexpected error = %v", err)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().
// 			Delete(gomock.Any(), productID).
// 			Return(testError).
// 			Times(1)

// 		err := service.Delete(context.Background(), productID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("Delete() error = %v, want %v", err, testError)
// 		}
// 	})
// }
