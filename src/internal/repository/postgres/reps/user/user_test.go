package user_rep

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate(t *testing.T) {
	fixture := NewTestFixture(t)
	testUser := fixture.userBuilder.Build()

	tests := []struct {
		name        string
		setupMock   func(structs.User)
		expectedID  uuid.UUID
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMock: func(user structs.User) {

				fixture.mock.ExpectQuery(`insert into "user"`).
					WithArgs(
						user.Name,
						user.Date_of_birth,
						user.Mail,
						sqlmock.AnyArg(),
						user.Phone,
						user.Address,
						user.Status,
						user.Role,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.Id))
			},
			expectedID:  testUser.Id,
			expectedErr: nil,
		},
		{
			name: "database error",
			setupMock: func(user structs.User) {

				fixture.mock.ExpectQuery(`insert into "user"`).
					WithArgs(
						user.Name,
						user.Date_of_birth,
						user.Mail,
						sqlmock.AnyArg(),
						user.Phone,
						user.Address,
						user.Status,
						user.Role,
					).
					WillReturnError(errTest)
			},
			expectedID:  uuid.UUID{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(testUser)

			id, err := fixture.repo.Create(fixture.ctx, testUser)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expectedID, id)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetById(t *testing.T) {
	fixture := NewTestFixture(t)

	testUser := fixture.userBuilder.Build()

	tests := []struct {
		name        string
		setupMock   func()
		expected    structs.User
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "mail", "phone", "address", "status", "role"}).
					AddRow(testUser.Id, testUser.Name, testUser.Date_of_birth, testUser.Mail, testUser.Phone, testUser.Address, testUser.Status, testUser.Role)
				fixture.mock.ExpectQuery(`select \* from "user" where id = \$1`).
					WithArgs(testUser.Id).
					WillReturnRows(rows)
			},
			expected:    testUser,
			expectedErr: nil,
		},
		{
			name: "user not found",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where id = \$1`).
					WithArgs(testUser.Id).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.User{},
			expectedErr: errors.New("failed to get user: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where id = \$1`).
					WithArgs(testUser.Id).
					WillReturnError(errTest)
			},
			expected:    structs.User{},
			expectedErr: errors.New("failed to get user: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetById(fixture.ctx, testUser.Id)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetByMail(t *testing.T) {
	fixture := NewTestFixture(t)

	testUser := fixture.userBuilder.WithMail("unique@example.com").Build()

	tests := []struct {
		name        string
		setupMock   func()
		expected    structs.User
		expectedErr error
	}{
		{
			name: "successful get by mail",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "mail", "password", "phone", "address", "status", "role"}).
					AddRow(testUser.Id, testUser.Name, testUser.Date_of_birth, testUser.Mail, "hashed_password", testUser.Phone, testUser.Address, testUser.Status, testUser.Role)
				fixture.mock.ExpectQuery(`select \* from "user" where mail = \$1`).
					WithArgs(testUser.Mail).
					WillReturnRows(rows)
			},
			expected: structs.User{
				Id:            testUser.Id,
				Name:          testUser.Name,
				Date_of_birth: testUser.Date_of_birth,
				Mail:          testUser.Mail,
				Password:      "hashed_password",
				Phone:         testUser.Phone,
				Address:       testUser.Address,
				Status:        testUser.Status,
				Role:          testUser.Role,
			},
			expectedErr: nil,
		},
		{
			name: "user not found by mail",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where mail = \$1`).
					WithArgs(testUser.Mail).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.User{},
			expectedErr: errors.New("failed to get user by mail: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where mail = \$1`).
					WithArgs(testUser.Mail).
					WillReturnError(errTest)
			},
			expected:    structs.User{},
			expectedErr: errors.New("failed to get user by mail: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetByMail(fixture.ctx, testUser.Mail)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetAllUsers(t *testing.T) {
	fixture := NewTestFixture(t)

	testUsers := []structs.User{
		fixture.userBuilder.WithName("User 1").Build(),
		fixture.userBuilder.WithName("User 2").WithMail("user2@example.com").Build(),
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []structs.User
		expectedErr error
	}{
		{
			name: "successful get all users",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "mail", "phone", "address", "status", "role"})
				for _, user := range testUsers {
					rows.AddRow(user.Id, user.Name, user.Date_of_birth, user.Mail, user.Phone, user.Address, user.Status, user.Role)
				}
				fixture.mock.ExpectQuery(`select \* from "user" where role = \$1`).
					WithArgs("обычный пользователь").
					WillReturnRows(rows)
			},
			expected:    testUsers,
			expectedErr: nil,
		},
		{
			name: "empty list",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "mail", "phone", "address", "status", "role"})
				fixture.mock.ExpectQuery(`select \* from "user" where role = \$1`).
					WithArgs("обычный пользователь").
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where role = \$1`).
					WithArgs("обычный пользователь").
					WillReturnError(errTest)
			},
			expected:    nil,
			expectedErr: errors.New("failed to get user by mail: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetAllUsers(fixture.ctx)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestGetByPhone(t *testing.T) {
	fixture := NewTestFixture(t)

	testUser := fixture.userBuilder.WithPhone("+1234567890").Build()

	tests := []struct {
		name        string
		setupMock   func()
		expected    structs.User
		expectedErr error
	}{
		{
			name: "successful get by phone",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "mail", "password", "phone", "address", "status", "role"}).
					AddRow(testUser.Id, testUser.Name, testUser.Date_of_birth, testUser.Mail, "hashed_password", testUser.Phone, testUser.Address, testUser.Status, testUser.Role)
				fixture.mock.ExpectQuery(`select \* from "user" where phone = \$1`).
					WithArgs(testUser.Phone).
					WillReturnRows(rows)
			},
			expected: structs.User{
				Id:            testUser.Id,
				Name:          testUser.Name,
				Date_of_birth: testUser.Date_of_birth,
				Mail:          testUser.Mail,
				Password:      "hashed_password",
				Phone:         testUser.Phone,
				Address:       testUser.Address,
				Status:        testUser.Status,
				Role:          testUser.Role,
			},
			expectedErr: nil,
		},
		{
			name: "user not found by phone",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where phone = \$1`).
					WithArgs(testUser.Phone).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    structs.User{},
			expectedErr: errors.New("failed to get user by phone: " + sql.ErrNoRows.Error()),
		},
		{
			name: "database error",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select \* from "user" where phone = \$1`).
					WithArgs(testUser.Phone).
					WillReturnError(errTest)
			},
			expected:    structs.User{},
			expectedErr: errors.New("failed to get user by phone: " + errTest.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := fixture.repo.GetByPhone(fixture.ctx, testUser.Phone)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expected, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
