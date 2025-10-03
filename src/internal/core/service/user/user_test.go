package user

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	defer fixture.Cleanup()

	testUser := fixture.userBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockUserRepository, *mock_structs.MockUsrBasket, *mock_structs.MockUsrFavourites, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, mockBasket *mock_structs.MockUsrBasket, mockFav *mock_structs.MockUsrFavourites, userID uuid.UUID) {
				mockRepo.EXPECT().Create(fixture.ctx, testUser).Return(userID, nil)

				mockBasket.EXPECT().Create(fixture.ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, basket structs.Basket) error {
					assert.Equal(t, userID, basket.IdUser)
					return nil
				})

				mockFav.EXPECT().Create(fixture.ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fav structs.Favourites) error {
					assert.Equal(t, userID, fav.IdUser)
					return nil
				})
			},
			expectedErr: nil,
		},
		{
			name: "repository error on user creation",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, mockBasket *mock_structs.MockUsrBasket, mockFav *mock_structs.MockUsrFavourites, userID uuid.UUID) {
				mockRepo.EXPECT().Create(fixture.ctx, testUser).Return(uuid.Nil, errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "repository error on basket creation",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, mockBasket *mock_structs.MockUsrBasket, mockFav *mock_structs.MockUsrFavourites, userID uuid.UUID) {
				mockRepo.EXPECT().Create(fixture.ctx, testUser).Return(userID, nil)
				mockBasket.EXPECT().Create(fixture.ctx, gomock.Any()).Return(errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "repository error on favourites creation",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, mockBasket *mock_structs.MockUsrBasket, mockFav *mock_structs.MockUsrFavourites, userID uuid.UUID) {
				mockRepo.EXPECT().Create(fixture.ctx, testUser).Return(userID, nil)
				mockBasket.EXPECT().Create(fixture.ctx, gomock.Any()).Return(nil)
				mockFav.EXPECT().Create(fixture.ctx, gomock.Any()).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, mockBasket, mockFav := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, mockBasket, mockFav, testUser.Id)

			err := service.Create(fixture.ctx, testUser)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testUser := fixture.userBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockUserRepository, structs.User)
		expectedRet structs.User
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, user structs.User) {
				mockRepo.EXPECT().GetById(fixture.ctx, user.Id).Return(user, nil)
			},
			expectedRet: testUser,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, user structs.User) {
				mockRepo.EXPECT().GetById(fixture.ctx, user.Id).Return(structs.User{}, errTest)
			},
			expectedRet: structs.User{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testUser)

			ret, err := service.GetById(fixture.ctx, testUser.Id)

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

func TestGetByMail_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testUser := fixture.userBuilder.WithMail("test@example.com").Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockUserRepository, structs.User)
		expectedRet structs.User
		expectedErr error
	}{
		{
			name: "successful get by mail",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, user structs.User) {
				mockRepo.EXPECT().GetByMail(fixture.ctx, user.Mail).Return(user, nil)
			},
			expectedRet: testUser,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, user structs.User) {
				mockRepo.EXPECT().GetByMail(fixture.ctx, user.Mail).Return(structs.User{}, errTest)
			},
			expectedRet: structs.User{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testUser)

			ret, err := service.GetByMail(fixture.ctx, testUser.Mail)

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

func TestGetAllUsers_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testUsers := []structs.User{
		fixture.userBuilder.WithName("User 1").Build(),
		fixture.userBuilder.WithName("User 2").WithMail("user2@example.com").Build(),
	}

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockUserRepository, []structs.User)
		expectedRet []structs.User
		expectedErr error
	}{
		{
			name: "successful get all users",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, users []structs.User) {
				mockRepo.EXPECT().GetAllUsers(fixture.ctx).Return(users, nil)
			},
			expectedRet: testUsers,
			expectedErr: nil,
		},
		{
			name: "empty list",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, users []structs.User) {
				mockRepo.EXPECT().GetAllUsers(fixture.ctx).Return([]structs.User{}, nil)
			},
			expectedRet: []structs.User{},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, users []structs.User) {
				mockRepo.EXPECT().GetAllUsers(fixture.ctx).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testUsers)

			ret, err := service.GetAllUsers(fixture.ctx)

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

func TestGetByPhone_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	testUser := fixture.userBuilder.WithPhone("+1234567890").Build()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockUserRepository, structs.User)
		expectedRet structs.User
		expectedErr error
	}{
		{
			name: "successful get by phone",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, user structs.User) {
				mockRepo.EXPECT().GetByPhone(fixture.ctx, user.Phone).Return(user, nil)
			},
			expectedRet: testUser,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockUserRepository, user structs.User) {
				mockRepo.EXPECT().GetByPhone(fixture.ctx, user.Phone).Return(structs.User{}, errTest)
			},
			expectedRet: structs.User{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testUser)

			ret, err := service.GetByPhone(fixture.ctx, testUser.Phone)

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
