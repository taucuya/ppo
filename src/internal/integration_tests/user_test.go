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

	err = runSQLScripts(db, []string{
		"/app/internal/database/sql/01-create.sql",
		"/app/internal/database/sql/02-constraints.sql",
		"/app/internal/database/sql/trigger_accept.sql",
		"/app/internal/database/sql/trigger_order.sql",
	})

	// err = runSQLScripts(db, []string{
	// 	"/home/runner/work/ppo/ppo/src/internal/database/sql/01-create.sql",
	// 	"/home/runner/work/ppo/ppo/src/internal/database/sql/02-constraints.sql",
	// 	"/home/runner/work/ppo/ppo/src/internal/database/sql/03-inserts.sql",
	// 	"/home/runner/work/ppo/ppo/src/internal/database/sql/trigger_accept.sql",
	// 	"/home/runner/work/ppo/ppo/src/internal/database/sql/trigger_order.sql",
	// })

	if err != nil {
		panic("failed to run SQL scripts: " + err.Error())
	}

	code := m.Run()

	_ = db.Close()
	os.Exit(code)
}

func runSQLScripts(db *sqlx.DB, scripts []string) error {
	for _, path := range scripts {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", path, err)
		}
	}
	return nil
}

func truncateTables(t *testing.T) {
	tables := []string{
		"review", "product", "brand", "order_item", "order_worker",
		"\"order\"", "basket_item", "basket", "worker", "\"user\"", "favourites",
	}

	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

type UserTestFixture struct {
	t        *testing.T
	ctx      context.Context
	service  *user.Service
	userRepo *user_rep.Repository
}

func NewUserTestFixture(t *testing.T) *UserTestFixture {
	userRepo := user_rep.New(db)
	basketRepo := basket_rep.New(db)
	favRepo := favourites_rep.New(db)

	service := user.New(userRepo, basketRepo, favRepo)

	return &UserTestFixture{
		t:        t,
		ctx:      context.Background(),
		service:  service,
		userRepo: userRepo,
	}
}

func (f *UserTestFixture) createTestUser() structs.User {
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	return structs.User{
		Name:          "Test User",
		Date_of_birth: dob,
		Mail:          "test@example.com",
		Password:      "password123",
		Phone:         "89016475843",
		Address:       "123 Test St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *UserTestFixture) createAnotherTestUser() structs.User {
	dob, _ := time.Parse("2006-01-02", "1995-05-15")
	return structs.User{
		Name:          "Another Test User",
		Date_of_birth: dob,
		Mail:          "another@example.com",
		Password:      "password456",
		Phone:         "89016475844",
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

func TestUser_Create_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)
	testUser := fixture.createTestUser()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		user        structs.User
		expectedErr error
	}{
		{
			name: "successful user creation with basket and favourites",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			user:        testUser,
			expectedErr: nil,
		},
		{
			name: "fail to create user with duplicate email",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			user:        testUser,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, tt.user)

			if tt.expectedErr != nil {
				require.Error(t, err)
			} else if tt.name == "successful user creation with basket and favourites" {
				require.NoError(t, err)

				createdUser, err := fixture.userRepo.GetByMail(fixture.ctx, tt.user.Mail)
				require.NoError(t, err)
				fixture.assertUserEqual(tt.user, createdUser)

				var basketCount int
				err = db.Get(&basketCount, "SELECT COUNT(*) FROM basket WHERE id_user = $1", createdUser.Id)
				require.NoError(t, err)
				require.Equal(t, 1, basketCount)

				var favCount int
				err = db.Get(&favCount, "SELECT COUNT(*) FROM favourites WHERE id_user = $1", createdUser.Id)
				require.NoError(t, err)
				require.Equal(t, 1, favCount)
			}
		})
	}
}

func TestUser_GetById_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)
	testUser := fixture.createTestUser()

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr error
	}{
		{
			name: "successfully get user by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				id, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
				return id
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: nil,
		},
		{
			name: "fail to get non-existent user by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			defer tt.cleanup()

			result, err := fixture.service.GetById(fixture.ctx, userID)

			if tt.name == "successfully get user by id" {
				require.NoError(t, err)
				fixture.assertUserEqual(testUser, result)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestUser_GetByMail_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)
	testUser := fixture.createTestUser()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		mail        string
		expectedErr error
	}{
		{
			name: "successfully get user by mail",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			mail:        testUser.Mail,
			expectedErr: nil,
		},
		{
			name: "fail to get user by non-existent mail",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			mail:        "nonexistent@example.com",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			result, err := fixture.service.GetByMail(fixture.ctx, tt.mail)

			if tt.name == "successfully get user by mail" {
				require.NoError(t, err)
				fixture.assertUserEqual(testUser, result)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestUser_GetByPhone_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)
	testUser := fixture.createTestUser()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		phone       string
		expectedErr error
	}{
		{
			name: "successfully get user by phone",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			phone:       testUser.Phone,
			expectedErr: nil,
		},
		{
			name: "fail to get user by non-existent phone",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			phone:       "89000000000",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			result, err := fixture.service.GetByPhone(fixture.ctx, tt.phone)

			if tt.name == "successfully get user by phone" {
				require.NoError(t, err)
				fixture.assertUserEqual(testUser, result)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestUser_GetAllUsers_AAA(t *testing.T) {
	fixture := NewUserTestFixture(t)
	testUser1 := fixture.createTestUser()
	testUser2 := fixture.createAnotherTestUser()

	tests := []struct {
		name          string
		setup         func()
		cleanup       func()
		expectedCount int
		expectedErr   error
	}{
		{
			name: "successfully get all users",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser1)
				require.NoError(t, err)
				_, err = fixture.userRepo.Create(fixture.ctx, testUser2)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 2,
			expectedErr:   nil,
		},
		{
			name: "get empty users list",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			users, err := fixture.service.GetAllUsers(fixture.ctx)

			require.NoError(t, err)
			require.Len(t, users, tt.expectedCount)
		})
	}
}
