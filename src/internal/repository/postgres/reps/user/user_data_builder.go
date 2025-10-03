package user_rep

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	structs "github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type UserBuilder struct {
	user structs.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: structs.User{
			Id:            uuid.New(),
			Name:          "Test User",
			Date_of_birth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Mail:          "test@example.com",
			Phone:         "+1234567890",
			Address:       "123 Test Street",
			Status:        "active",
			Role:          "обычный пользователь",
		},
	}
}

func (b *UserBuilder) WithID(id uuid.UUID) *UserBuilder {
	b.user.Id = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

func (b *UserBuilder) WithDateOfBirth(dob time.Time) *UserBuilder {
	b.user.Date_of_birth = dob
	return b
}

func (b *UserBuilder) WithMail(mail string) *UserBuilder {
	b.user.Mail = mail
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.user.Password = password
	return b
}

func (b *UserBuilder) WithPhone(phone string) *UserBuilder {
	b.user.Phone = phone
	return b
}

func (b *UserBuilder) WithAddress(address string) *UserBuilder {
	b.user.Address = address
	return b
}

func (b *UserBuilder) WithStatus(status string) *UserBuilder {
	b.user.Status = status
	return b
}

func (b *UserBuilder) WithRole(role string) *UserBuilder {
	b.user.Role = role
	return b
}

func (b *UserBuilder) Build() structs.User {
	return b.user
}

type TestFixture struct {
	t           *testing.T
	db          *sql.DB
	sqlxDB      *sqlx.DB
	mock        sqlmock.Sqlmock
	repo        *Repository
	ctx         context.Context
	userBuilder *UserBuilder
}

func NewTestFixture(t *testing.T) *TestFixture {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return &TestFixture{
		t:           t,
		db:          db,
		sqlxDB:      sqlxDB,
		mock:        mock,
		repo:        New(sqlxDB),
		ctx:         context.Background(),
		userBuilder: NewUserBuilder(),
	}
}

func (f *TestFixture) Cleanup() {
	f.db.Close()
}

func (f *TestFixture) AssertError(actual, expected error) {
	if expected == nil {
		assert.NoError(f.t, actual)
	} else {
		assert.EqualError(f.t, actual, expected.Error())
	}
}
