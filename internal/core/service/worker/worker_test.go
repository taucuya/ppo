package worker

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

// 	mockRepo := mock_structs.NewMockWorkerRepository(ctrl)
// 	service := New(mockRepo)
// 	testWorker := structs.Worker{
// 		Id:       structs.GenId(),
// 		IdUser:   structs.GenId(),
// 		JobTitle: "Engineer",
// 	}

// 	t.Run("successful creation", func(t *testing.T) {
// 		mockRepo.EXPECT().Create(gomock.Any(), testWorker).Return(nil).Times(1)

// 		err := service.Create(context.Background(), testWorker)
// 		if err != nil {
// 			t.Errorf("Create() unexpected error = %v", err)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().Create(gomock.Any(), testWorker).Return(testError).Times(1)

// 		err := service.Create(context.Background(), testWorker)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("Create() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestGetById(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockWorkerRepository(ctrl)
// 	service := New(mockRepo)
// 	testID := structs.GenId()
// 	testWorker := structs.Worker{Id: testID, IdUser: structs.GenId(), JobTitle: "Engineer"}

// 	t.Run("successful get", func(t *testing.T) {
// 		mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(testWorker, nil).Times(1)

// 		got, err := service.GetById(context.Background(), testID)
// 		if err != nil {
// 			t.Errorf("GetById() unexpected error = %v", err)
// 		}
// 		if got != testWorker {
// 			t.Errorf("GetById() = %v, want %v", got, testWorker)
// 		}
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(structs.Worker{}, testError).Times(1)

// 		_, err := service.GetById(context.Background(), testID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("GetById() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestDelete(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockWorkerRepository(ctrl)
// 	service := New(mockRepo)
// 	testID := structs.GenId()

// 	t.Run("successful delete", func(t *testing.T) {
// 		mockRepo.EXPECT().Delete(gomock.Any(), testID).Return(nil).Times(1)

// 		err := service.Delete(context.Background(), testID)
// 		if err != nil {
// 			t.Errorf("Delete() unexpected error = %v", err)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().Delete(gomock.Any(), testID).Return(testError).Times(1)

// 		err := service.Delete(context.Background(), testID)
// 		if !errors.Is(err, testError) {
// 			t.Errorf("Delete() error = %v, want %v", err, testError)
// 		}
// 	})
// }

// func TestGetAllWorkers(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockWorkerRepository(ctrl)
// 	service := New(mockRepo)
// 	testWorkers := []structs.Worker{
// 		{Id: structs.GenId(), IdUser: structs.GenId(), JobTitle: "Engineer"},
// 		{Id: structs.GenId(), IdUser: structs.GenId(), JobTitle: "Manager"},
// 	}

// 	t.Run("successful get all", func(t *testing.T) {
// 		mockRepo.EXPECT().GetAllWorkers(gomock.Any()).Return(testWorkers, nil).Times(1)

// 		got, err := service.GetAllWorkers(context.Background())
// 		if err != nil {
// 			t.Errorf("GetAllWorkers() unexpected error = %v", err)
// 		}
// 		if len(got) != len(testWorkers) {
// 			t.Errorf("GetAllWorkers() = %v, want %v", got, testWorkers)
// 		}
// 	})

// 	t.Run("repository error", func(t *testing.T) {
// 		mockRepo.EXPECT().GetAllWorkers(gomock.Any()).Return(nil, testError).Times(1)

// 		_, err := service.GetAllWorkers(context.Background())
// 		if !errors.Is(err, testError) {
// 			t.Errorf("GetAllWorkers() error = %v, want %v", err, testError)
// 		}
// 	})
// }
