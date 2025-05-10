package basket

// import (
// 	"context"
// 	"errors"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/taucuya/ppo/internal/core/mock_structs"
// 	"github.com/taucuya/ppo/internal/core/structs"
// )

// var testError = errors.New("test error")

// func TestCreate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockBasketRepository(ctrl)
// 	service := New(mockRepo)

// 	now := time.Now()
// 	testDate := time.Date(
// 		now.Year(),
// 		now.Month(),
// 		now.Day(),
// 		0, 0, 0, 0,
// 		time.UTC,
// 	)

// 	testID := structs.GenId()
// 	userID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		basket  structs.Basket
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful creation",
// 			basket: structs.Basket{
// 				Id:     testID,
// 				IdUser: userID,
// 				Date:   testDate,
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().Create(gomock.Any(), structs.Basket{
// 					Id:     testID,
// 					IdUser: userID,
// 					Date:   testDate,
// 				}).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository error",
// 			basket: structs.Basket{
// 				Id:     testID,
// 				IdUser: userID,
// 				Date:   testDate,
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(testError)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := service.Create(context.Background(), tt.basket)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestGetById(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockBasketRepository(ctrl)
// 	service := New(mockRepo)

// 	now := time.Now()
// 	testDate := time.Date(
// 		now.Year(),
// 		now.Month(),
// 		now.Day(),
// 		0, 0, 0, 0,
// 		time.UTC,
// 	)

// 	testID := structs.GenId()
// 	userID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		id      uuid.UUID
// 		mock    func()
// 		want    structs.Basket
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful get",
// 			id:   testID,
// 			mock: func() {
// 				mockRepo.EXPECT().GetById(gomock.Any(), testID).Return(structs.Basket{
// 					Id:     testID,
// 					IdUser: userID,
// 					Date:   testDate,
// 				}, nil)
// 			},
// 			want: structs.Basket{
// 				Id:     testID,
// 				IdUser: userID,
// 				Date:   testDate,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "not found",
// 			id:   structs.GenId(),
// 			mock: func() {
// 				mockRepo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(structs.Basket{}, testError)
// 			},
// 			want:    structs.Basket{},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.GetById(context.Background(), tt.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetById() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// // Остальные тесты остаются без изменений, так как они не работают с датами
// func TestGetItems(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockBasketRepository(ctrl)
// 	service := New(mockRepo)

// 	basketID := structs.GenId()
// 	productID := structs.GenId()
// 	itemID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		id      uuid.UUID
// 		mock    func()
// 		want    []structs.BasketItem
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful get items",
// 			id:   basketID,
// 			mock: func() {
// 				mockRepo.EXPECT().GetItems(gomock.Any(), basketID).Return([]structs.BasketItem{
// 					{
// 						Id:        itemID,
// 						IdProduct: productID,
// 						IdBasket:  basketID,
// 						Amount:    2,
// 					},
// 				}, nil)
// 			},
// 			want: []structs.BasketItem{
// 				{
// 					Id:        itemID,
// 					IdProduct: productID,
// 					IdBasket:  basketID,
// 					Amount:    2,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "empty basket",
// 			id:   basketID,
// 			mock: func() {
// 				mockRepo.EXPECT().GetItems(gomock.Any(), basketID).Return([]structs.BasketItem{}, nil)
// 			},
// 			want:    []structs.BasketItem{},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository error",
// 			id:   basketID,
// 			mock: func() {
// 				mockRepo.EXPECT().GetItems(gomock.Any(), basketID).Return(nil, testError)
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.GetItems(context.Background(), tt.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetItems() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAddItem(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockBasketRepository(ctrl)
// 	service := New(mockRepo)

// 	basketID := structs.GenId()
// 	productID := structs.GenId()
// 	itemID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		item    structs.BasketItem
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful add",
// 			item: structs.BasketItem{
// 				Id:        itemID,
// 				IdProduct: productID,
// 				IdBasket:  basketID,
// 				Amount:    1,
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().AddItem(gomock.Any(), structs.BasketItem{
// 					Id:        itemID,
// 					IdProduct: productID,
// 					IdBasket:  basketID,
// 					Amount:    1,
// 				}).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository error",
// 			item: structs.BasketItem{
// 				Id:        itemID,
// 				IdProduct: productID,
// 				IdBasket:  basketID,
// 				Amount:    1,
// 			},
// 			mock: func() {
// 				mockRepo.EXPECT().AddItem(gomock.Any(), gomock.Any()).Return(testError)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := service.AddItem(context.Background(), tt.item)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestDeleteItem(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockBasketRepository(ctrl)
// 	service := New(mockRepo)

// 	itemID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		id      uuid.UUID
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful delete",
// 			id:   itemID,
// 			mock: func() {
// 				mockRepo.EXPECT().DeleteItem(gomock.Any(), itemID).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository error",
// 			id:   itemID,
// 			mock: func() {
// 				mockRepo.EXPECT().DeleteItem(gomock.Any(), itemID).Return(testError)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := service.DeleteItem(context.Background(), tt.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestUpdateItemAmount(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock_structs.NewMockBasketRepository(ctrl)
// 	service := New(mockRepo)

// 	itemID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		id      uuid.UUID
// 		amount  int
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name:   "successful update",
// 			id:     itemID,
// 			amount: 5,
// 			mock: func() {
// 				mockRepo.EXPECT().UpdateItemAmount(gomock.Any(), itemID, 5).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "zero amount",
// 			id:     itemID,
// 			amount: 0,
// 			mock: func() {
// 				mockRepo.EXPECT().UpdateItemAmount(gomock.Any(), itemID, 0).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "negative amount",
// 			id:     itemID,
// 			amount: -1,
// 			mock: func() {
// 				mockRepo.EXPECT().UpdateItemAmount(gomock.Any(), itemID, -1).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "repository error",
// 			id:     itemID,
// 			amount: 5,
// 			mock: func() {
// 				mockRepo.EXPECT().UpdateItemAmount(gomock.Any(), itemID, 5).Return(testError)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := service.UpdateItemAmount(context.Background(), tt.id, tt.amount)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateItemAmount() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 			}
// 		})
// 	}
// }
