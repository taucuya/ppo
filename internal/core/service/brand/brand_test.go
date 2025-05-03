package brand

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var testError = errors.New("test error")

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockBrandRepository(ctrl)
	service := New(mockRepo)

	testBrand := structs.Brand{
		Id:            structs.GenId(),
		Name:          "Test Brand",
		Description:   "Test Description",
		PriceCategory: "Premium",
	}

	tests := []struct {
		name    string
		brand   structs.Brand
		mock    func()
		wantErr bool
	}{
		{
			name:  "successful creation",
			brand: testBrand,
			mock: func() {
				mockRepo.EXPECT().Create(gomock.Any(), testBrand).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "repository error",
			brand: testBrand,
			mock: func() {
				mockRepo.EXPECT().Create(gomock.Any(), testBrand).Return(testError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := service.Create(context.Background(), tt.brand)
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

	mockRepo := mock_structs.NewMockBrandRepository(ctrl)
	service := New(mockRepo)

	testID := structs.GenId()
	testBrand := structs.Brand{
		Id:            testID,
		Name:          "Test Brand",
		Description:   "Test Description",
		PriceCategory: "Premium",
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    structs.Brand
		wantErr bool
	}{
		{
			name: "successful get",
			id:   testID,
			mock: func() {
				mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(testBrand, nil)
			},
			want:    testBrand,
			wantErr: false,
		},
		{
			name: "not found",
			id:   structs.GenId(),
			mock: func() {
				mockRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(structs.Brand{}, testError)
			},
			want:    structs.Brand{},
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
			if got != tt.want {
				t.Errorf("GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockBrandRepository(ctrl)
	service := New(mockRepo)

	testID := structs.GenId()

	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   testID,
			mock: func() {
				mockRepo.EXPECT().Delete(gomock.Any(), testID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			id:   testID,
			mock: func() {
				mockRepo.EXPECT().Delete(gomock.Any(), testID).Return(testError)
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

func TestGetAllBrands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_structs.NewMockBrandRepository(ctrl)
	service := New(mockRepo)

	testBrands := []structs.Brand{
		{
			Id:            structs.GenId(),
			Name:          "Brand 1",
			Description:   "Desc 1",
			PriceCategory: "Standard",
		},
		{
			Id:            structs.GenId(),
			Name:          "Brand 2",
			Description:   "Desc 2",
			PriceCategory: "Premium",
		},
	}

	tests := []struct {
		name    string
		mock    func()
		want    []structs.Brand
		wantErr bool
	}{
		{
			name: "successful get all",
			mock: func() {
				mockRepo.EXPECT().GetAllBrands(gomock.Any()).Return(testBrands, nil)
			},
			want:    testBrands,
			wantErr: false,
		},
		{
			name: "empty list",
			mock: func() {
				mockRepo.EXPECT().GetAllBrands(gomock.Any()).Return([]structs.Brand{}, nil)
			},
			want:    []structs.Brand{},
			wantErr: false,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().GetAllBrands(gomock.Any()).Return(nil, testError)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetAllBrands(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBrands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, testError) {
				t.Errorf("Expected testError, got %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetAllBrands() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetAllBrands() = %v, want %v", got, tt.want)
					break
				}
			}
		})
	}
}
