package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var testError = errors.New("test error")

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockUserRepository(ctrl)
	service := New(mockRepo)
	testUser := structs.User{
		Id:       structs.GenId(),
		Name:     "Test User",
		Mail:     "test@example.com",
		Password: "securepassword",
		Phone:    "1234567890",
		Address:  "Test Address",
		Status:   "Active",
		Role:     "User",
	}

	t.Run("successful creation", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), testUser).Return(nil).Times(1)

		err := service.Create(context.Background(), testUser)
		if err != nil {
			t.Errorf("Create() unexpected error = %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), testUser).Return(testError).Times(1)

		err := service.Create(context.Background(), testUser)
		if !errors.Is(err, testError) {
			t.Errorf("Create() error = %v, want %v", err, testError)
		}
	})
}

func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockUserRepository(ctrl)
	service := New(mockRepo)
	testID := structs.GenId()
	testUser := structs.User{
		Id: testID, Name: "Test User", Mail: "test@example.com", Phone: "1234567890",
	}

	t.Run("successful get", func(t *testing.T) {
		mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(testUser, nil).Times(1)

		got, err := service.GetById(context.Background(), testID)
		if err != nil {
			t.Errorf("GetById() unexpected error = %v", err)
		}
		if got != testUser {
			t.Errorf("GetById() = %v, want %v", got, testUser)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(structs.User{}, testError).Times(1)

		_, err := service.GetById(context.Background(), testID)
		if !errors.Is(err, testError) {
			t.Errorf("GetById() error = %v, want %v", err, testError)
		}
	})
}

func TestGetByMail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockUserRepository(ctrl)
	service := New(mockRepo)

	mail := "test@example.com"
	testUser := structs.User{Mail: mail, Name: "Test User", Phone: "1234567890"}

	t.Run("successful get", func(t *testing.T) {
		mockRepo.EXPECT().GetByMail(gomock.Any(), mail).Return(testUser, nil).Times(1)

		got, err := service.GetByMail(context.Background(), mail)
		if err != nil {
			t.Errorf("GetByMail() unexpected error = %v", err)
		}
		if got != testUser {
			t.Errorf("GetByMail() = %v, want %v", got, testUser)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByMail(gomock.Any(), mail).Return(structs.User{}, testError).Times(1)

		_, err := service.GetByMail(context.Background(), mail)
		if !errors.Is(err, testError) {
			t.Errorf("GetByMail() error = %v, want %v", err, testError)
		}
	})
}

func TestGetByPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockUserRepository(ctrl)
	service := New(mockRepo)

	phone := "1234567890"
	testUser := structs.User{Phone: phone, Name: "Test User", Mail: "test@example.com"}

	t.Run("successful get", func(t *testing.T) {
		mockRepo.EXPECT().GetByPhone(gomock.Any(), phone).Return(testUser, nil).Times(1)

		got, err := service.GetByPhone(context.Background(), phone)
		if err != nil {
			t.Errorf("GetByPhone() unexpected error = %v", err)
		}
		if got != testUser {
			t.Errorf("GetByPhone() = %v, want %v", got, testUser)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByPhone(gomock.Any(), phone).Return(structs.User{}, testError).Times(1)

		_, err := service.GetByPhone(context.Background(), phone)
		if !errors.Is(err, testError) {
			t.Errorf("GetByPhone() error = %v, want %v", err, testError)
		}
	})
}
