package user

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/service/basket"
	"github.com/taucuya/ppo/internal/core/service/favourites"
	"github.com/taucuya/ppo/internal/core/structs"
	basket_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/basket"
	favourites_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/favourites"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

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
	fixture.Cleanup()
}

func TestClassicCreate_AAA(t *testing.T) {
	fixture := NewTestClassicFixture(t)
	testUser := fixture.userBuilder.Build()

	tests := []struct {
		name        string
		setupMocks  func(sqlmock.Sqlmock, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(mock sqlmock.Sqlmock, userID uuid.UUID) {
				mock.ExpectQuery(`insert into "user"`).
					WithArgs(
						testUser.Name,
						testUser.Date_of_birth,
						testUser.Mail,
						sqlmock.AnyArg(),
						testUser.Phone,
						testUser.Address,
						testUser.Status,
						testUser.Role,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

				mock.ExpectExec(`insert into basket`).
					WithArgs(
						userID,
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(`insert into favourites`).
					WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "repository error on user creation",
			setupMocks: func(mock sqlmock.Sqlmock, userID uuid.UUID) {
				mock.ExpectQuery(`insert into "user"`).
					WithArgs(
						testUser.Name,
						testUser.Date_of_birth,
						testUser.Mail,
						sqlmock.AnyArg(),
						testUser.Phone,
						testUser.Address,
						testUser.Status,
						testUser.Role,
					).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "repository error on basket creation",
			setupMocks: func(mock sqlmock.Sqlmock, userID uuid.UUID) {
				mock.ExpectQuery(`insert into "user"`).
					WithArgs(
						testUser.Name,
						testUser.Date_of_birth,
						testUser.Mail,
						sqlmock.AnyArg(),
						testUser.Phone,
						testUser.Address,
						testUser.Status,
						testUser.Role,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

				mock.ExpectExec(`insert into basket`).
					WithArgs(
						userID,
						sqlmock.AnyArg(),
					).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
		{
			name: "repository error on favourites creation",
			setupMocks: func(mock sqlmock.Sqlmock, userID uuid.UUID) {
				mock.ExpectQuery(`insert into "user"`).
					WithArgs(
						testUser.Name,
						testUser.Date_of_birth,
						testUser.Mail,
						sqlmock.AnyArg(),
						testUser.Phone,
						testUser.Address,
						testUser.Status,
						testUser.Role,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

				mock.ExpectExec(`insert into basket`).
					WithArgs(
						userID,
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(`insert into favourites`).
					WithArgs(userID).
					WillReturnError(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			defer sqlxDB.Close()
			userRepo := user_rep.New(sqlxDB)
			basketRepo := basket_rep.New(sqlxDB)
			favRepo := favourites_rep.New(sqlxDB)
			favouritesSvc := favourites.New(favRepo)
			basketSvc := basket.New(basketRepo)
			service := New(userRepo, basketSvc, favouritesSvc)
			tt.setupMocks(mock, testUser.Id)

			err = service.Create(fixture.ctx, testUser)
			fixture.AssertClassicError(err, tt.expectedErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
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
