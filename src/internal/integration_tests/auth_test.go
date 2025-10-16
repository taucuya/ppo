package integrationtests

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/taucuya/ppo/internal/core/service/auth"
	"github.com/taucuya/ppo/internal/core/service/basket"
	"github.com/taucuya/ppo/internal/core/service/favourites"
	"github.com/taucuya/ppo/internal/core/service/user"
	"github.com/taucuya/ppo/internal/core/structs"
	auth_prov "github.com/taucuya/ppo/internal/providers/jwt/auth"
	auth_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/auth"
	basket_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/basket"
	favourites_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/favourites"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

type AuthTestFixture struct {
	t        *testing.T
	ctx      context.Context
	service  *auth.Service
	authRepo *auth_rep.Repository
	userRepo *user_rep.Repository
	provider *auth_prov.Provider
}

func NewAuthTestFixture(t *testing.T) *AuthTestFixture {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	authRepo := auth_rep.New(db)
	basketRepo := basket_rep.New(db)
	basket := basket.New(basketRepo)
	favouritesRepo := favourites_rep.New(db)
	favourites := favourites.New(favouritesRepo)
	userRepo := user_rep.New(db)
	userServ := user.New(userRepo, basket, favourites)
	provider := auth_prov.New(key, 15*time.Minute, 24*time.Hour)

	service := auth.New(provider, authRepo, userServ)

	return &AuthTestFixture{
		t:        t,
		ctx:      context.Background(),
		service:  service,
		authRepo: authRepo,
		userRepo: userRepo,
		provider: provider,
	}
}

func (f *AuthTestFixture) createTestUser() structs.User {
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	return structs.User{
		Name:          "Test User",
		Date_of_birth: dob,
		Mail:          "test@example.com",
		Password:      string(hashedPassword),
		Phone:         "89016475843",
		Address:       "123 Test St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *AuthTestFixture) setupRegularUser() uuid.UUID {
	testUser := f.createTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)
	return userID
}

func TestAuth_SignUp_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)
	testUser := fixture.createTestUser()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		user        structs.User
		expectedErr bool
	}{
		{
			name: "successful user registration",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			user:        testUser,
			expectedErr: false,
		},
		{
			name: "fail to register user with duplicate email",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			user:        testUser,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			err := fixture.service.SignUp(fixture.ctx, tt.user)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				createdUser, err := fixture.userRepo.GetByMail(fixture.ctx, tt.user.Mail)
				require.NoError(t, err)
				require.Equal(t, tt.user.Name, createdUser.Name)
				require.Equal(t, tt.user.Mail, createdUser.Mail)

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

func TestAuth_LogIn_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)
	testUser := fixture.createTestUser()

	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		mail        string
		password    string
		expectedErr bool
	}{
		{
			name: "successful login with correct credentials",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			mail:        testUser.Mail,
			password:    testUser.Password,
			expectedErr: false,
		},
		{
			name: "fail to login with incorrect password",
			setup: func() {
				truncateTables(t)
				_, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			mail:        testUser.Mail,
			password:    "password123",
			expectedErr: true,
		},
		{
			name: "fail to login with non-existent email",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			mail:        "nonexistent@example.com",
			password:    "password123",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			atoken, rtoken, err := fixture.service.LogIn(fixture.ctx, tt.mail, tt.password)

			if tt.expectedErr {
				require.Error(t, err)
				require.Empty(t, atoken)
				require.Empty(t, rtoken)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, atoken)
				require.NotEmpty(t, rtoken)

				tokenID, err := fixture.authRepo.VerifyToken(fixture.ctx, rtoken)
				require.NoError(t, err)
				require.NotEqual(t, uuid.Nil, tokenID)
			}
		})
	}
}

func TestAuth_LogOut_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name        string
		setup       func() string
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful logout with valid token",
			setup: func() string {
				truncateTables(t)
				userID := fixture.setupRegularUser()

				_, rtoken, err := fixture.provider.GenToken(fixture.ctx, userID)
				require.NoError(t, err)

				err = fixture.authRepo.CreateToken(fixture.ctx, userID, rtoken)
				require.NoError(t, err)

				return rtoken
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to logout with invalid token",
			setup: func() string {
				truncateTables(t)
				return "invalid-token"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtoken := tt.setup()
			defer tt.cleanup()

			err := fixture.service.LogOut(fixture.ctx, rtoken)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := fixture.authRepo.VerifyToken(fixture.ctx, rtoken)
				require.Error(t, err)
			}
		})
	}
}

func TestAuth_RefreshToken_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (string, string)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful token refresh with existing user",
			setup: func() (string, string) {
				truncateTables(t)

				testUser := fixture.createTestUser()
				err := fixture.service.SignUp(fixture.ctx, testUser)
				require.NoError(t, err)

				atoken, rtoken, err := fixture.service.LogIn(fixture.ctx, testUser.Mail, testUser.Password)
				require.NoError(t, err)

				tokenID, err := fixture.authRepo.VerifyToken(fixture.ctx, rtoken)
				require.NoError(t, err)
				require.NotEqual(t, uuid.Nil, tokenID)

				return atoken, rtoken
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to refresh with non-existent user token",
			setup: func() (string, string) {
				truncateTables(t)

				testUser1 := fixture.createTestUser()
				err := fixture.service.SignUp(fixture.ctx, testUser1)
				require.NoError(t, err)
				atoken1, _, err := fixture.service.LogIn(fixture.ctx, testUser1.Mail, testUser1.Password)
				require.NoError(t, err)

				testUser2 := structs.User{
					Name:          "Another User",
					Date_of_birth: testUser1.Date_of_birth,
					Mail:          "another@example.com",
					Password:      "password456",
					Phone:         "89016475844",
					Address:       "456 Another St",
					Status:        "active",
					Role:          "обычный пользователь",
				}
				err = fixture.service.SignUp(fixture.ctx, testUser2)
				require.NoError(t, err)
				_, rtoken2, err := fixture.service.LogIn(fixture.ctx, testUser2.Mail, testUser2.Password)
				require.NoError(t, err)

				tokenID, err := fixture.authRepo.VerifyToken(fixture.ctx, rtoken2)
				require.NoError(t, err)
				err = fixture.authRepo.DeleteToken(fixture.ctx, tokenID)
				require.NoError(t, err)

				_, err = fixture.authRepo.VerifyToken(fixture.ctx, rtoken2)
				require.Error(t, err)

				return atoken1, rtoken2
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
		{
			name: "fail to refresh with completely invalid token",
			setup: func() (string, string) {
				truncateTables(t)

				testUser := fixture.createTestUser()
				err := fixture.service.SignUp(fixture.ctx, testUser)
				require.NoError(t, err)

				atoken, _, err := fixture.service.LogIn(fixture.ctx, testUser.Mail, testUser.Password)
				require.NoError(t, err)

				return atoken, "completely.invalid.token.12345"
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
		{
			name: "fail to refresh with mismatched tokens",
			setup: func() (string, string) {
				truncateTables(t)

				testUser1 := fixture.createTestUser()
				err := fixture.service.SignUp(fixture.ctx, testUser1)
				require.NoError(t, err)

				testUser2 := structs.User{
					Name:          "Another User",
					Date_of_birth: testUser1.Date_of_birth,
					Mail:          "another@example.com",
					Password:      "password456",
					Phone:         "89016475844",
					Address:       "456 Another St",
					Status:        "active",
					Role:          "обычный пользователь",
				}
				err = fixture.service.SignUp(fixture.ctx, testUser2)
				require.NoError(t, err)

				atoken1, _, err := fixture.service.LogIn(fixture.ctx, testUser1.Mail, testUser1.Password)
				require.NoError(t, err)

				_, rtoken2, err := fixture.service.LogIn(fixture.ctx, testUser2.Mail, testUser2.Password)
				require.NoError(t, err)

				return atoken1, rtoken2
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atoken, rtoken := tt.setup()
			defer tt.cleanup()

			newToken, err := fixture.service.RefreshToken(fixture.ctx, atoken, rtoken)

			if tt.expectedErr {
				require.Error(t, err)
				require.Empty(t, newToken)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, newToken)

				valid, err := fixture.provider.VerifyToken(fixture.ctx, newToken)
				require.NoError(t, err)
				require.True(t, valid)
			}
		})
	}
}

func TestAuth_VerifyTokens_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name            string
		setup           func() (string, string)
		cleanup         func()
		expectedAccess  bool
		expectedRefresh bool
		expectNewToken  bool
	}{
		{
			name: "both tokens valid",
			setup: func() (string, string) {
				truncateTables(t)
				userID := fixture.setupRegularUser()

				atoken, rtoken, err := fixture.provider.GenToken(fixture.ctx, userID)
				require.NoError(t, err)

				err = fixture.authRepo.CreateToken(fixture.ctx, userID, rtoken)
				require.NoError(t, err)

				return atoken, rtoken
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedAccess:  true,
			expectedRefresh: true,
			expectNewToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atoken, rtoken := tt.setup()
			defer tt.cleanup()

			newToken, accessValid, refreshValid, err := fixture.service.VerifyTokens(fixture.ctx, atoken, rtoken)

			require.NoError(t, err)
			require.Equal(t, tt.expectedAccess, accessValid)
			require.Equal(t, tt.expectedRefresh, refreshValid)

			if tt.expectNewToken {
				require.NotEmpty(t, newToken)
			} else {
				require.Empty(t, newToken)
			}
		})
	}
}

func TestAuth_GetId_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (string, uuid.UUID)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully extract user ID from valid token",
			setup: func() (string, uuid.UUID) {
				truncateTables(t)
				userID := fixture.setupRegularUser()

				atoken, _, err := fixture.provider.GenToken(fixture.ctx, userID)
				require.NoError(t, err)

				return atoken, userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to extract user ID from invalid token",
			setup: func() (string, uuid.UUID) {
				truncateTables(t)
				return "invalid-token", uuid.Nil
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expectedID := tt.setup()
			defer tt.cleanup()

			userID, err := fixture.service.GetId(token)

			if tt.expectedErr {
				require.Error(t, err)
				require.Equal(t, uuid.Nil, userID)
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedID, userID)
			}
		})
	}
}
