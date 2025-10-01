package auth_rep

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		token       string
		expectedErr error
	}{
		{
			name: "successful token creation",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into token \\(rtoken\\) values \\(\\$1\\)").
					WithArgs(fixture.token).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			token:       fixture.token,
			expectedErr: nil,
		},
		{
			name: "token creation error",
			setupMock: func() {
				fixture.mock.ExpectExec("insert into token \\(rtoken\\) values \\(\\$1\\)").
					WithArgs(fixture.token).
					WillReturnError(errTest)
			},
			token:       fixture.token,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.CreateToken(fixture.ctx, fixture.userID, tt.token)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestCheckAdmin(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         uuid.UUID
		expectedResult bool
	}{
		{
			name: "user is admin",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.userID)
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1 and job_title = \$2`).
					WithArgs(fixture.userID, "admin").
					WillReturnRows(rows)
			},
			userID:         fixture.userID,
			expectedResult: true,
		},
		{
			name: "user is not admin",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1 and job_title = \$2`).
					WithArgs(fixture.userID, "admin").
					WillReturnError(sql.ErrNoRows)
			},
			userID:         fixture.userID,
			expectedResult: false,
		},
		{
			name: "database error when checking admin",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1 and job_title = \$2`).
					WithArgs(fixture.userID, "admin").
					WillReturnError(errTest)
			},
			userID:         fixture.userID,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result := fixture.repo.CheckAdmin(fixture.ctx, tt.userID)

			assert.Equal(t, tt.expectedResult, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestCheckWorker(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         uuid.UUID
		expectedResult bool
	}{
		{
			name: "user is worker",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.userID)
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1 and job_title = \$2`).
					WithArgs(fixture.userID, "worker").
					WillReturnRows(rows)
			},
			userID:         fixture.userID,
			expectedResult: true,
		},
		{
			name: "user is not worker",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1 and job_title = \$2`).
					WithArgs(fixture.userID, "worker").
					WillReturnError(sql.ErrNoRows)
			},
			userID:         fixture.userID,
			expectedResult: false,
		},
		{
			name: "database error when checking worker",
			setupMock: func() {
				fixture.mock.ExpectQuery(`select id from worker where id_user = \$1 and job_title = \$2`).
					WithArgs(fixture.userID, "worker").
					WillReturnError(errTest)
			},
			userID:         fixture.userID,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result := fixture.repo.CheckWorker(fixture.ctx, tt.userID)

			assert.Equal(t, tt.expectedResult, result)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestVerifyToken(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		token       string
		expectedID  uuid.UUID
		expectedErr error
	}{
		{
			name: "successful token verification",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(fixture.userID)
				fixture.mock.ExpectQuery("select id from token where rtoken = \\$1").
					WithArgs(fixture.token).
					WillReturnRows(rows)
			},
			token:       fixture.token,
			expectedID:  fixture.userID,
			expectedErr: nil,
		},
		{
			name: "token not found",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from token where rtoken = \\$1").
					WithArgs(fixture.token).
					WillReturnError(sql.ErrNoRows)
			},
			token:       fixture.token,
			expectedID:  uuid.UUID{},
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "database error when verifying token",
			setupMock: func() {
				fixture.mock.ExpectQuery("select id from token where rtoken = \\$1").
					WithArgs(fixture.token).
					WillReturnError(errTest)
			},
			token:       fixture.token,
			expectedID:  uuid.UUID{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			id, err := fixture.repo.VerifyToken(fixture.ctx, tt.token)

			fixture.AssertError(err, tt.expectedErr)
			assert.Equal(t, tt.expectedID, id)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}

func TestDeleteToken(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMock   func()
		tokenID     uuid.UUID
		expectedErr error
	}{
		{
			name: "successful token deletion",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from token where id = \\$1").
					WithArgs(fixture.userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			tokenID:     fixture.userID,
			expectedErr: nil,
		},
		{
			name: "database error when deleting token",
			setupMock: func() {
				fixture.mock.ExpectExec("delete from token where id = \\$1").
					WithArgs(fixture.userID).
					WillReturnError(errTest)
			},
			tokenID:     fixture.userID,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := fixture.repo.DeleteToken(fixture.ctx, tt.tokenID)

			fixture.AssertError(err, tt.expectedErr)
			require.NoError(t, fixture.mock.ExpectationsWereMet())
		})
	}
	fixture.Cleanup()
}
