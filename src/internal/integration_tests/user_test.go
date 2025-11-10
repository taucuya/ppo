package integrationtests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/user"
	"github.com/taucuya/ppo/internal/core/structs"
	basket_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/basket"
	favourites_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/favourites"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error

	dsn := "postgres://test_user:test_password@postgres:5432/test_db?sslmode=disable"

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("database not available: " + err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

type UserTestFixture struct {
	t          *testing.T
	ctx        context.Context
	service    *user.Service
	userRepo   *user_rep.Repository
	basketRepo *basket_rep.Repository
	favRepo    *favourites_rep.Repository
	testID     string
}

func NewUserTestFixture(t *testing.T) *UserTestFixture {
	userRepo := user_rep.New(db)
	basketRepo := basket_rep.New(db)
	favRepo := favourites_rep.New(db)

	service := user.New(userRepo, basketRepo, favRepo)

	testID := uuid.New().String()[:8]

	return &UserTestFixture{
		t:          t,
		ctx:        context.Background(),
		service:    service,
		userRepo:   userRepo,
		basketRepo: basketRepo,
		favRepo:    favRepo,
		testID:     testID,
	}
}

func (f *UserTestFixture) generateTestUser() structs.User {
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%s-%d", f.testID, timestamp)

	dob, _ := time.Parse("2006-01-02", "1990-01-01")

	phoneSuffix := fmt.Sprintf("%09d", timestamp%1000000000)
	phone := "89" + phoneSuffix
	if len(phone) > 11 {
		phone = phone[:11]
	}

	return structs.User{
		Name:          fmt.Sprintf("Test User %s", uniqueID),
		Date_of_birth: dob,
		Mail:          fmt.Sprintf("test%s@example.com", uniqueID),
		Password:      "password123",
		Phone:         phone,
		Address:       "123 Test St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *UserTestFixture) generateAnotherTestUser() structs.User {
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%s-%d", f.testID, timestamp)

	dob, _ := time.Parse("2006-01-02", "1995-05-15")

	phoneSuffix := fmt.Sprintf("%09d", (timestamp+1)%1000000000)
	phone := "89" + phoneSuffix
	if len(phone) > 11 {
		phone = phone[:11]
	}

	return structs.User{
		Name:          fmt.Sprintf("Another Test User %s", uniqueID),
		Date_of_birth: dob,
		Mail:          fmt.Sprintf("another%s@example.com", uniqueID),
		Password:      "password456",
		Phone:         phone,
		Address:       "456 Another St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *UserTestFixture) assertUserEqual(expected, actual structs.User) {
	require.Equal(f.t, expected.Name, actual.Name)
	require.Equal(f.t, expected.Mail, actual.Mail)
	require.Equal(f.t, expected.Phone, actual.Phone)
	require.Equal(f.t, expected.Address, actual.Address)
	require.Equal(f.t, expected.Status, actual.Status)
	require.Equal(f.t, expected.Role, actual.Role)
	require.True(f.t, expected.Date_of_birth.Equal(actual.Date_of_birth))
}

func (f *UserTestFixture) cleanupUserData(userID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM favourites WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM basket_item WHERE id_basket IN (SELECT id FROM basket WHERE id_user = $1)", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM basket WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"user\" WHERE id = $1", userID)
}

func TestUser_Create_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.User, []uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successful user creation with basket and favourites",
			setup: func() (structs.User, []uuid.UUID) {
				user := fixture.generateTestUser()
				return user, []uuid.UUID{}
			},
			expectedErr: false,
		},
		{
			name: "fail to create user with duplicate email",
			setup: func() (structs.User, []uuid.UUID) {
				user := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, user)
				require.NoError(t, err)
				return user, []uuid.UUID{userID}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, cleanupIDs := tt.setup()

			defer func() {
				for _, id := range cleanupIDs {
					if id != uuid.Nil {
						fixture.cleanupUserData(id)
					}
				}
			}()

			err := fixture.service.Create(fixture.ctx, user)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				createdUser, err := fixture.userRepo.GetByMail(fixture.ctx, user.Mail)
				require.NoError(t, err)
				fixture.assertUserEqual(user, createdUser)

				var basketCount int
				err = db.Get(&basketCount, "SELECT COUNT(*) FROM basket WHERE id_user = $1", createdUser.Id)
				require.NoError(t, err)
				require.Equal(t, 1, basketCount)

				var favCount int
				err = db.Get(&favCount, "SELECT COUNT(*) FROM favourites WHERE id_user = $1", createdUser.Id)
				require.NoError(t, err)
				require.Equal(t, 1, favCount)

				defer fixture.cleanupUserData(createdUser.Id)
			}
		})
	}
}

func TestUser_GetById_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully get user by id",
			setup: func() uuid.UUID {
				user := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, user)
				require.NoError(t, err)
				return userID
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent user by id",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

			result, err := fixture.service.GetById(fixture.ctx, userID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, userID, result.Id)
			}
		})
	}
}

func TestUser_GetByMail_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (string, []uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successfully get user by mail",
			setup: func() (string, []uuid.UUID) {
				user := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, user)
				require.NoError(t, err)
				return user.Mail, []uuid.UUID{userID}
			},
			expectedErr: false,
		},
		{
			name: "fail to get user by non-existent mail",
			setup: func() (string, []uuid.UUID) {
				return "nonexistent@example.com", []uuid.UUID{}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mail, cleanupIDs := tt.setup()

			defer func() {
				for _, id := range cleanupIDs {
					if id != uuid.Nil {
						fixture.cleanupUserData(id)
					}
				}
			}()

			result, err := fixture.service.GetByMail(fixture.ctx, mail)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, mail, result.Mail)
			}
		})
	}
}

func TestUser_GetByPhone_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (string, []uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successfully get user by phone",
			setup: func() (string, []uuid.UUID) {
				user := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, user)
				require.NoError(t, err)
				return user.Phone, []uuid.UUID{userID}
			},
			expectedErr: false,
		},
		{
			name: "fail to get user by non-existent phone",
			setup: func() (string, []uuid.UUID) {
				return "89000000000", []uuid.UUID{}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, cleanupIDs := tt.setup()

			defer func() {
				for _, id := range cleanupIDs {
					if id != uuid.Nil {
						fixture.cleanupUserData(id)
					}
				}
			}()

			result, err := fixture.service.GetByPhone(fixture.ctx, phone)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, phone, result.Phone)
			}
		})
	}
}
