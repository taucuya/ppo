package auth

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUp_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockAuthUser)
		expectedErr error
	}{
		{
			name: "successful sign-up",
			setupMocks: func(mockUser *mock_structs.MockAuthUser) {
				mockUser.EXPECT().Create(fixture.ctx, fixture.testUser).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "sign-up error(repository)",
			setupMocks: func(mockUser *mock_structs.MockAuthUser) {
				mockUser.EXPECT().Create(fixture.ctx, fixture.testUser).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _, mockUser := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockUser)

			err := service.SignUp(fixture.ctx, fixture.testUser)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestLogIn_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	atoken, rtoken := "access-token", "refresh-token"
	password := "password123"

	tests := []struct {
		name           string
		mail           string
		password       string
		setupMocks     func(*mock_structs.MockAuthUser, *mock_structs.MockAuthProvider, *mock_structs.MockAuthRepository)
		expectedAToken string
		expectedRToken string
		expectedErr    error
	}{
		{
			name:     "successful login",
			mail:     fixture.testUser.Mail,
			password: password,
			setupMocks: func(mockUser *mock_structs.MockAuthUser, mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockUser.EXPECT().GetByMail(fixture.ctx, fixture.testUser.Mail).Return(fixture.testUser, nil)
				mockProv.EXPECT().GenToken(fixture.ctx, fixture.testUser.Id).Return(atoken, rtoken, nil)
				mockRepo.EXPECT().CreateToken(fixture.ctx, fixture.testUser.Id, rtoken).Return(nil)
			},
			expectedAToken: atoken,
			expectedRToken: rtoken,
			expectedErr:    nil,
		},
		{
			name:     "login error (user not found)",
			mail:     fixture.testUser.Mail,
			password: fixture.testUser.Password,
			setupMocks: func(mockUser *mock_structs.MockAuthUser, mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockUser.EXPECT().GetByMail(fixture.ctx, fixture.testUser.Mail).Return(structs.User{}, errTest)
			},
			expectedAToken: "",
			expectedRToken: "",
			expectedErr:    errTest,
		},
		{
			name:     "login error (wrong password)",
			mail:     fixture.testUser.Mail,
			password: "wrongpassword",
			setupMocks: func(mockUser *mock_structs.MockAuthUser, mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockUser.EXPECT().GetByMail(fixture.ctx, fixture.testUser.Mail).Return(fixture.testUser, nil)
			},
			expectedAToken: "",
			expectedRToken: "",
			expectedErr:    bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:     "login error (token generation error)",
			mail:     fixture.testUser.Mail,
			password: password,
			setupMocks: func(mockUser *mock_structs.MockAuthUser, mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockUser.EXPECT().GetByMail(fixture.ctx, fixture.testUser.Mail).Return(fixture.testUser, nil)
				mockProv.EXPECT().GenToken(fixture.ctx, fixture.testUser.Id).Return("", "", errTest)
			},
			expectedAToken: "",
			expectedRToken: "",
			expectedErr:    errTest,
		},
		{
			name:     "login error (repository error)",
			mail:     fixture.testUser.Mail,
			password: password,
			setupMocks: func(mockUser *mock_structs.MockAuthUser, mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockUser.EXPECT().GetByMail(fixture.ctx, fixture.testUser.Mail).Return(fixture.testUser, nil)
				mockProv.EXPECT().GenToken(fixture.ctx, fixture.testUser.Id).Return(atoken, rtoken, nil)
				mockRepo.EXPECT().CreateToken(fixture.ctx, fixture.testUser.Id, rtoken).Return(errTest)
			},
			expectedAToken: "",
			expectedRToken: "",
			expectedErr:    errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockProv, mockRepo, mockUser := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockUser, mockProv, mockRepo)
			gotAToken, gotRToken, err := service.LogIn(fixture.ctx, tt.mail, tt.password)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAToken, gotAToken)
				assert.Equal(t, tt.expectedRToken, gotRToken)
			}
		})
	}
	fixture.Cleanup()
}

func TestLogOut_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	rtoken := "refresh-token"

	tests := []struct {
		name        string
		rtoken      string
		setupMocks  func(*mock_structs.MockAuthProvider, *mock_structs.MockAuthRepository)
		expectedErr error
	}{
		{
			name:   "successful logout",
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(true, nil)
				mockRepo.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(fixture.testUser.Id, nil)
				mockRepo.EXPECT().DeleteToken(fixture.ctx, fixture.testUser.Id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "logout error (token verification provider error)",
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(false, errTest)
			},
			expectedErr: errTest,
		},
		{
			name:   "logout error (token verification repo error)",
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(true, nil)
				mockRepo.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(uuid.UUID{}, errTest)
			},
			expectedErr: errTest,
		},
		{
			name:   "logout error (repository delete token error)",
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(true, nil)
				mockRepo.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(fixture.testUser.Id, nil)
				mockRepo.EXPECT().DeleteToken(fixture.ctx, fixture.testUser.Id).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockProv, mockRepo, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockProv, mockRepo)

			err := service.LogOut(fixture.ctx, tt.rtoken)

			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestRefreshToken_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	atoken := "access-token"
	rtoken := "refresh-token"
	newToken := "new-access-token"

	tests := []struct {
		name        string
		atoken      string
		rtoken      string
		setupMocks  func(*mock_structs.MockAuthProvider)
		expected    string
		expectedErr error
	}{
		{
			name:   "successful refresh",
			atoken: atoken,
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider) {
				mockProv.EXPECT().RefreshToken(fixture.ctx, atoken, rtoken).Return(newToken, nil)
			},
			expected:    newToken,
			expectedErr: nil,
		},
		{
			name:   "refresh error",
			atoken: atoken,
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider) {
				mockProv.EXPECT().RefreshToken(fixture.ctx, atoken, rtoken).Return("", errTest)
			},
			expected:    "",
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockProv, _, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockProv)

			got, err := service.RefreshToken(fixture.ctx, tt.atoken, tt.rtoken)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
	fixture.Cleanup()
}

func TestVerifyTokens_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	atoken := "access-token"
	rtoken := "refresh-token"
	newToken := "new-access-token"

	tests := []struct {
		name            string
		atoken          string
		rtoken          string
		setupMocks      func(*mock_structs.MockAuthProvider, *mock_structs.MockAuthRepository)
		expectedToken   string
		expectedAccess  bool
		expectedRefresh bool
		expectedErr     error
	}{
		{
			name:   "both tokens valid",
			atoken: atoken,
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, atoken).Return(true, nil)
				mockProv.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(true, nil)
				mockRepo.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(fixture.testUser.Id, nil)
			},
			expectedAccess:  true,
			expectedRefresh: true,
			expectedErr:     nil,
		},
		{
			name:   "access expired, refresh valid - successful refresh",
			atoken: "expired-token",
			rtoken: rtoken,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, "expired-token").Return(false, jwt.ErrTokenExpired)
				mockProv.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(true, nil)
				mockRepo.EXPECT().VerifyToken(fixture.ctx, rtoken).Return(fixture.testUser.Id, nil)
				mockProv.EXPECT().RefreshToken(fixture.ctx, "expired-token", rtoken).Return(newToken, nil)
			},
			expectedToken:   newToken,
			expectedAccess:  true,
			expectedRefresh: true,
			expectedErr:     nil,
		},
		{
			name:   "refresh token invalid",
			atoken: atoken,
			rtoken: "invalid-refresh",
			setupMocks: func(mockProv *mock_structs.MockAuthProvider, mockRepo *mock_structs.MockAuthRepository) {
				mockProv.EXPECT().VerifyToken(fixture.ctx, atoken).Return(true, nil)
				mockProv.EXPECT().VerifyToken(fixture.ctx, "invalid-refresh").Return(false, errTest)
			},
			expectedAccess:  true,
			expectedRefresh: false,
			expectedErr:     errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockProv, mockRepo, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockProv, mockRepo)

			newToken, accessValid, refreshValid, err := service.VerifyTokens(fixture.ctx, tt.atoken, tt.rtoken)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, newToken)
				assert.Equal(t, tt.expectedAccess, accessValid)
				assert.Equal(t, tt.expectedRefresh, refreshValid)
			}
		})
	}
	fixture.Cleanup()
}

func TestCheckRoles_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	tests := []struct {
		name           string
		setupMock      func(*mock_structs.MockAuthRepository) bool
		checkFunc      func(*Service, context.Context, uuid.UUID) bool
		expectedResult bool
	}{
		{
			name: "User is admin",
			setupMock: func(mockRepo *mock_structs.MockAuthRepository) bool {
				mockRepo.EXPECT().CheckAdmin(fixture.ctx, fixture.testUser.Id).Return(true)
				return true
			},
			checkFunc: func(s *Service, ctx context.Context, id uuid.UUID) bool {
				return s.CheckAdmin(ctx, id)
			},
			expectedResult: true,
		},
		{
			name: "User is not admin",
			setupMock: func(mockRepo *mock_structs.MockAuthRepository) bool {
				mockRepo.EXPECT().CheckAdmin(fixture.ctx, fixture.testUser.Id).Return(false)
				return false
			},
			checkFunc: func(s *Service, ctx context.Context, id uuid.UUID) bool {
				return s.CheckAdmin(ctx, id)
			},
			expectedResult: false,
		},
		{
			name: "User is worker",
			setupMock: func(mockRepo *mock_structs.MockAuthRepository) bool {
				mockRepo.EXPECT().CheckWorker(fixture.ctx, fixture.testUser.Id).Return(true)
				return true
			},
			checkFunc: func(s *Service, ctx context.Context, id uuid.UUID) bool {
				return s.CheckWorker(ctx, id)
			},
			expectedResult: true,
		},
		{
			name: "User is not worker",
			setupMock: func(mockRepo *mock_structs.MockAuthRepository) bool {
				mockRepo.EXPECT().CheckWorker(fixture.ctx, fixture.testUser.Id).Return(false)
				return false
			},
			checkFunc: func(s *Service, ctx context.Context, id uuid.UUID) bool {
				return s.CheckWorker(ctx, id)
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, mockRepo, _ := fixture.CreateServiceWithMocks()
			expected := tt.setupMock(mockRepo)

			result := tt.checkFunc(service, fixture.ctx, fixture.testUser.Id)

			assert.Equal(t, expected, result)
		})
	}
	fixture.Cleanup()
}

func TestGetId_AAA(t *testing.T) {
	fixture := NewTestFixture(t)

	token := "test-token"

	tests := []struct {
		name        string
		token       string
		setupMocks  func(*mock_structs.MockAuthProvider)
		expectedID  uuid.UUID
		expectedErr error
	}{
		{
			name:  "successful ID extraction",
			token: token,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider) {
				mockProv.EXPECT().ExtractUserID(token).Return(fixture.testUser.Id, nil)
			},
			expectedID:  fixture.testUser.Id,
			expectedErr: nil,
		},
		{
			name:  "extraction error",
			token: token,
			setupMocks: func(mockProv *mock_structs.MockAuthProvider) {
				mockProv.EXPECT().ExtractUserID(token).Return(uuid.UUID{}, errTest)
			},
			expectedID:  uuid.UUID{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockProv, _, _ := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockProv)

			gotID, err := service.GetId(tt.token)

			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, gotID)
			}
		})
	}
	fixture.Cleanup()
}
