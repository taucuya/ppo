package basket

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockBasketRepository)
		expectedRet structs.Basket
		expectedErr error
	}{
		{
			name: "successful get",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.basket.Id).Return(structs.Basket{
					Id:     fixture.basket.Id,
					IdUser: fixture.basket.IdUser,
					Date:   fixture.basket.Date,
				}, nil)
			},
			expectedRet: structs.Basket{
				Id:     fixture.basket.Id,
				IdUser: fixture.basket.IdUser,
				Date:   fixture.basket.Date,
			},
			expectedErr: nil,
		},
		{
			name: "error get (not found basket by user)",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(uuid.UUID{}, errTest)
			},
			expectedRet: structs.Basket{},
			expectedErr: errTest,
		},
		{
			name: "error get (not found basket by basket id)",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.basket.Id).Return(structs.Basket{}, errTest)
			},
			expectedRet: structs.Basket{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetById(fixture.ctx, fixture.basket.IdUser)

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

func TestGetItems_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(mockRepo *mock_structs.MockBasketRepository)
		expectedRet []structs.BasketItem
		expectedErr error
	}{
		{
			name: "successful get items",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetItems(gomock.Any(), fixture.basket.Id).Return([]structs.BasketItem{
					{
						Id:        fixture.basketItems[0].Id,
						IdProduct: fixture.basketItems[0].IdProduct,
						IdBasket:  fixture.basketItems[0].IdBasket,
						Amount:    fixture.basketItems[0].Amount,
					},
					{
						Id:        fixture.basketItems[1].Id,
						IdProduct: fixture.basketItems[1].IdProduct,
						IdBasket:  fixture.basketItems[1].IdBasket,
						Amount:    fixture.basketItems[1].Amount,
					},
				}, nil)
			},
			expectedRet: []structs.BasketItem{
				{
					Id:        fixture.basketItems[0].Id,
					IdProduct: fixture.basketItems[0].IdProduct,
					IdBasket:  fixture.basketItems[0].IdBasket,
					Amount:    fixture.basketItems[0].Amount,
				},
				{
					Id:        fixture.basketItems[1].Id,
					IdProduct: fixture.basketItems[1].IdProduct,
					IdBasket:  fixture.basketItems[1].IdBasket,
					Amount:    fixture.basketItems[1].Amount,
				},
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetItems(gomock.Any(), fixture.basket.Id).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetItems(fixture.ctx, fixture.basket.Id)

			if err != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}

		})
	}
	fixture.Cleanup()
}

func TestAddItem_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	itemId := structs.GenId()
	productId := structs.GenId()

	tests := []struct {
		name        string
		item        structs.BasketItem
		setupMocks  func(mockRepo *mock_structs.MockBasketRepository)
		expectedErr error
	}{
		{
			name: "successful add",
			item: structs.BasketItem{
				Id:        itemId,
				IdProduct: productId,
				IdBasket:  fixture.basket.Id,
				Amount:    1,
			},
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().AddItem(fixture.ctx, structs.BasketItem{
					Id:        itemId,
					IdProduct: productId,
					IdBasket:  fixture.basket.Id,
					Amount:    1,
				}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			item: structs.BasketItem{
				Id:        itemId,
				IdProduct: productId,
				IdBasket:  fixture.basket.Id,
				Amount:    1,
			},
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().AddItem(fixture.ctx, gomock.Any()).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.AddItem(fixture.ctx, tt.item, fixture.basket.IdUser)

			fixture.AssertError(err, tt.expectedErr)

		})
	}
	fixture.Cleanup()
}

func TestDeleteItem(t *testing.T) {
	fixture := NewTestFixture(t)

	productId := structs.GenId()

	tests := []struct {
		name        string
		setupMocks  func(mockRepo *mock_structs.MockBasketRepository)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().DeleteItem(fixture.ctx, fixture.basket.Id, productId).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().DeleteItem(fixture.ctx, fixture.basket.Id, productId).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.DeleteItem(fixture.ctx, fixture.basket.IdUser, productId)

			fixture.AssertError(err, tt.expectedErr)

		})
	}
}

func TestUpdateItemAmount(t *testing.T) {
	fixture := NewTestFixture(t)

	productId := structs.GenId()
	amount := 5

	tests := []struct {
		name        string
		amount      int
		setupMocks  func(mockRepo *mock_structs.MockBasketRepository)
		expectedErr error
	}{
		{
			name:   "successful update",
			amount: amount,
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().UpdateItemAmount(fixture.ctx, fixture.basket.Id, productId, amount).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "repository error",
			amount: amount,
			setupMocks: func(mockRepo *mock_structs.MockBasketRepository) {
				mockRepo.EXPECT().GetBIdByUId(fixture.ctx, fixture.basket.IdUser).Return(fixture.basket.Id, nil)
				mockRepo.EXPECT().UpdateItemAmount(fixture.ctx, fixture.basket.Id, productId, 5).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.UpdateItemAmount(fixture.ctx, fixture.basket.IdUser, productId, tt.amount)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
}
