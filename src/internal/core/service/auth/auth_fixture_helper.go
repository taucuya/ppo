package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
	"golang.org/x/crypto/bcrypt"
)

var errTest = errors.New("test error")

type TestFixture struct {
	t        *testing.T
	ctrl     *gomock.Controller
	ctx      context.Context
	testUser structs.User
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := structs.GenId()

	return &TestFixture{
		t:    t,
		ctrl: ctrl,
		ctx:  context.Background(),
		testUser: structs.User{
			Id:       userID,
			Mail:     "test@example.com",
			Password: string(hashedPassword),
		},
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockAuthProvider,
	*mock_structs.MockAuthRepository, *mock_structs.MockAuthUser) {
	mockProv := mock_structs.NewMockAuthProvider(f.ctrl)
	mockRepo := mock_structs.NewMockAuthRepository(f.ctrl)
	mockUser := mock_structs.NewMockAuthUser(f.ctrl)

	service := New(mockProv, mockRepo, mockUser)
	return service, mockProv, mockRepo, mockUser
}

func (f *TestFixture) AssertError(err error, expectedErr error) {
	if expectedErr != nil {
		if err == nil {
			f.t.Errorf("Expected error %v, got nil", expectedErr)
			return
		} else if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
			f.t.Errorf("Expected  error %v, got %v", expectedErr, err)
		}

	} else if err != nil {
		f.t.Errorf("Expected error nil, got %v", err)
		return
	}
}
