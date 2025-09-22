package favourites

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockFavouritesRepository)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().Create(fixture.ctx, fixture.favourites).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().Create(fixture.ctx, fixture.favourites).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.Create(fixture.ctx, fixture.favourites)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockFavouritesRepository)
		expectedRet structs.Favourites
		expectedErr error
	}{
		{
			name: "successful get",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favourites.Id, nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.favourites.Id).Return(fixture.favourites, nil)
			},
			expectedRet: fixture.favourites,
			expectedErr: nil,
		},
		{
			name: "error get (not found favourites id by user)",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(uuid.UUID{}, errTest)
			},
			expectedRet: structs.Favourites{},
			expectedErr: errTest,
		},
		{
			name: "error get (not found favourites by id)",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favourites.Id, nil)
				mockRepo.EXPECT().GetById(fixture.ctx, fixture.favourites.Id).Return(structs.Favourites{}, errTest)
			},
			expectedRet: structs.Favourites{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetById(fixture.ctx, fixture.favourites.IdUser)

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
		setupMocks  func(*mock_structs.MockFavouritesRepository)
		expectedRet []structs.FavouritesItem
		expectedErr error
	}{
		{
			name: "successful get items",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetItems(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favouritesItems, nil)
			},
			expectedRet: fixture.favouritesItems,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetItems(fixture.ctx, fixture.favourites.IdUser).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			ret, err := service.GetItems(fixture.ctx, fixture.favourites.IdUser)

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

func TestAddItem_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	item := structs.FavouritesItem{
		Id:           structs.GenId(),
		IdProduct:    structs.GenId(),
		IdFavourites: fixture.favourites.Id,
	}

	tests := []struct {
		name        string
		item        structs.FavouritesItem
		setupMocks  func(*mock_structs.MockFavouritesRepository)
		expectedErr error
	}{
		{
			name: "successful add",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favourites.Id, nil)
				mockRepo.EXPECT().AddItem(fixture.ctx, item).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error get favourites id",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(uuid.UUID{}, errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favourites.Id, nil)
				expectedItem := item
				expectedItem.IdFavourites = fixture.favourites.Id
				mockRepo.EXPECT().AddItem(fixture.ctx, expectedItem).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.AddItem(fixture.ctx, item, fixture.favourites.IdUser)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestDeleteItem_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	productID := structs.GenId()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockFavouritesRepository)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favourites.Id, nil)
				mockRepo.EXPECT().DeleteItem(fixture.ctx, fixture.favourites.Id, productID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error get favourites id",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(uuid.UUID{}, errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockFavouritesRepository) {
				mockRepo.EXPECT().GetFIdByUId(fixture.ctx, fixture.favourites.IdUser).Return(fixture.favourites.Id, nil)
				mockRepo.EXPECT().DeleteItem(fixture.ctx, fixture.favourites.Id, productID).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo)

			err := service.DeleteItem(fixture.ctx, fixture.favourites.IdUser, productID)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}
