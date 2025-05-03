package order

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var testError = errors.New("test error")

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockOrderRepository(ctrl)
	service := New(mockRepo)

	now := time.Now()
	testDate := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	testOrder := structs.Order{
		Id:       structs.GenId(),
		Date:     testDate,
		IdUser:   structs.GenId(),
		Address:  "Test Address",
		Status:   "pending",
		Price:    100.50,
		IdWorker: structs.GenId(),
	}

	tests := []struct {
		name    string
		order   structs.Order
		mock    func()
		wantErr bool
	}{
		{
			name:  "successful creation",
			order: testOrder,
			mock: func() {
				mockRepo.EXPECT().Create(gomock.Any(), testOrder).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "repository error",
			order: testOrder,
			mock: func() {
				mockRepo.EXPECT().Create(gomock.Any(), testOrder).Return(testError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := service.Create(context.Background(), tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, testError) {
				t.Errorf("Expected testError, got %v", err)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockOrderRepository(ctrl)
	service := New(mockRepo)

	now := time.Now()
	testDate := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	testID := structs.GenId()
	testOrder := structs.Order{
		Id:       testID,
		Date:     testDate,
		IdUser:   structs.GenId(),
		Address:  "Test Address",
		Status:   "pending",
		Price:    100.50,
		IdWorker: structs.GenId(),
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    structs.Order
		wantErr bool
	}{
		{
			name: "successful get",
			id:   testID,
			mock: func() {
				mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(testOrder, nil)
			},
			want:    testOrder,
			wantErr: false,
		},
		{
			name: "not found",
			id:   structs.GenId(),
			mock: func() {
				mockRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(structs.Order{}, testError)
			},
			want:    structs.Order{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetById(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, testError) {
				t.Errorf("Expected testError, got %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockOrderRepository(ctrl)
	service := New(mockRepo)

	orderID := structs.GenId()
	productID := structs.GenId()
	itemID := structs.GenId()

	testItems := []structs.OrderItem{
		{
			Id:        itemID,
			IdProduct: productID,
			IdOrder:   orderID,
			Amount:    2,
		},
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    []structs.OrderItem
		wantErr bool
	}{
		{
			name: "successful get items",
			id:   orderID,
			mock: func() {
				mockRepo.EXPECT().GetItems(gomock.Any(), orderID).Return(testItems, nil)
			},
			want:    testItems,
			wantErr: false,
		},
		{
			name: "empty order",
			id:   orderID,
			mock: func() {
				mockRepo.EXPECT().GetItems(gomock.Any(), orderID).Return([]structs.OrderItem{}, nil)
			},
			want:    []structs.OrderItem{},
			wantErr: false,
		},
		{
			name: "repository error",
			id:   orderID,
			mock: func() {
				mockRepo.EXPECT().GetItems(gomock.Any(), orderID).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetItems(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, testError) {
				t.Errorf("Expected testError, got %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockOrderRepository(ctrl)
	service := New(mockRepo)

	orderID := structs.GenId()
	testStatus := "shipped"

	t.Run("successful get status", func(t *testing.T) {
		mockRepo.EXPECT().
			GetStatus(gomock.Any(), orderID).
			Return(testStatus, nil).
			Times(1)

		got, err := service.GetStatus(context.Background(), orderID)

		if err != nil {
			t.Errorf("GetStatus() unexpected error = %v", err)
		}
		if got != testStatus {
			t.Errorf("GetStatus() = %v, want %v", got, testStatus)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetStatus(gomock.Any(), orderID).
			Return("", testError).
			Times(1)

		_, err := service.GetStatus(context.Background(), orderID)

		if !errors.Is(err, testError) {
			t.Errorf("GetStatus() error = %v, want %v", err, testError)
		}
	})
}

func TestChangeOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockOrderRepository(ctrl)
	service := New(mockRepo)

	orderID := structs.GenId()
	currentStatus := "processing"
	newStatus := "shipped"

	tests := []struct {
		name      string
		id        uuid.UUID
		newStatus string
		mock      func()
		wantErr   bool
	}{
		{
			name:      "successful status change",
			id:        orderID,
			newStatus: newStatus,
			mock: func() {
				mockRepo.EXPECT().GetStatus(gomock.Any(), orderID).Return(currentStatus, nil)
				mockRepo.EXPECT().GetById(gomock.Any(), orderID).Return(structs.Order{}, nil)
				mockRepo.EXPECT().UpdateStatus(gomock.Any(), orderID, newStatus).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "same status",
			id:        orderID,
			newStatus: currentStatus,
			mock: func() {
				mockRepo.EXPECT().GetStatus(gomock.Any(), orderID).Return(currentStatus, nil)
			},
			wantErr: false,
		},
		{
			name:      "get status error",
			id:        orderID,
			newStatus: newStatus,
			mock: func() {
				mockRepo.EXPECT().GetStatus(gomock.Any(), orderID).Return("", testError)
			},
			wantErr: true,
		},
		{
			name:      "get order error",
			id:        orderID,
			newStatus: newStatus,
			mock: func() {
				mockRepo.EXPECT().GetStatus(gomock.Any(), orderID).Return(currentStatus, nil)
				mockRepo.EXPECT().GetById(gomock.Any(), orderID).Return(structs.Order{}, testError)
			},
			wantErr: true,
		},
		{
			name:      "update status error",
			id:        orderID,
			newStatus: newStatus,
			mock: func() {
				mockRepo.EXPECT().GetStatus(gomock.Any(), orderID).Return(currentStatus, nil)
				mockRepo.EXPECT().GetById(gomock.Any(), orderID).Return(structs.Order{}, nil)
				mockRepo.EXPECT().UpdateStatus(gomock.Any(), orderID, newStatus).Return(testError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := service.ChangeOrderStatus(context.Background(), tt.id, tt.newStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, testError) {
				t.Errorf("Expected testError, got %v", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockOrderRepository(ctrl)
	service := New(mockRepo)

	orderID := structs.GenId()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   orderID,
			mock: func() {
				mockRepo.EXPECT().Delete(gomock.Any(), orderID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			id:   orderID,
			mock: func() {
				mockRepo.EXPECT().Delete(gomock.Any(), orderID).Return(testError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := service.Delete(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, testError) {
				t.Errorf("Expected testError, got %v", err)
			}
		})
	}
}
