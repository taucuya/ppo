package integrationtests

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

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
	testID   string
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

	testID := uuid.New().String()[:8]

	return &AuthTestFixture{
		t:        t,
		ctx:      context.Background(),
		service:  service,
		authRepo: authRepo,
		userRepo: userRepo,
		provider: provider,
		testID:   testID,
	}
}

func (f *AuthTestFixture) generateTestUser() (structs.User, string) {
	timestamp := time.Now().UnixNano()
	randomUUID := uuid.New().String()[:8]

	uniqueID := fmt.Sprintf("%s-%d-%s", f.testID, timestamp, randomUUID)

	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	plainPassword := "password123"

	phonePrefix := "89"
	randomNumbers := fmt.Sprintf("%09d", time.Now().UnixNano()%1000000000)
	if len(randomNumbers) > 9 {
		randomNumbers = randomNumbers[:9]
	}
	phone := phonePrefix + randomNumbers

	return structs.User{
		Name:          fmt.Sprintf("Test User %s", uniqueID),
		Date_of_birth: dob,
		Mail:          fmt.Sprintf("test%s@example.com", uniqueID),
		Password:      plainPassword,
		Phone:         phone,
		Address:       fmt.Sprintf("123 Test St %s", uniqueID),
		Role:          "обычный пользователь",
	}, plainPassword
}

func (f *AuthTestFixture) createUserForTest() (uuid.UUID, structs.User, string) {
	testUser, plainPassword := f.generateTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)
	return userID, testUser, plainPassword
}

func (f *AuthTestFixture) cleanupUserData(userID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM basket_item WHERE id_basket IN (SELECT id FROM basket WHERE id_user = $1)", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM basket WHERE id_user = $1", userID)

	_, _ = db.ExecContext(f.ctx, "DELETE FROM favourites_item WHERE id_favourites IN (SELECT id FROM favourites WHERE id_user = $1)", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM favourites WHERE id_user = $1", userID)

	_, _ = db.ExecContext(f.ctx, "DELETE FROM review WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM order_item WHERE id_order IN (SELECT id FROM \"order\" WHERE id_user = $1)", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"order\" WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM worker WHERE id_user = $1", userID)

	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"user\" WHERE id = $1", userID)
}

func (f *AuthTestFixture) cleanupUserByEmail(email string) {
	var userID uuid.UUID
	err := db.GetContext(f.ctx, &userID, "SELECT id FROM \"user\" WHERE mail = $1", email)
	if err == nil && userID != uuid.Nil {
		f.cleanupUserData(userID)
	}
}

func TestAuth_SignUp_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name        string
		user        structs.User
		expectedErr bool
	}{
		{
			name:        "successful user registration",
			user:        structs.User{},
			expectedErr: false,
		},
		{
			name:        "fail to register user with duplicate email",
			user:        structs.User{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "fail to register user with duplicate email" {
				existingUser, _ := fixture.generateTestUser()
				_, err := fixture.userRepo.Create(fixture.ctx, existingUser)
				require.NoError(t, err)
				defer fixture.cleanupUserByEmail(existingUser.Mail)

				tt.user = existingUser
			} else {
				testUser, _ := fixture.generateTestUser()
				tt.user = testUser
			}

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
				err = db.GetContext(fixture.ctx, &basketCount, "SELECT COUNT(*) FROM basket WHERE id_user = $1", createdUser.Id)

				var favCount int
				err = db.GetContext(fixture.ctx, &favCount, "SELECT COUNT(*) FROM favourites WHERE id_user = $1", createdUser.Id)

				fixture.cleanupUserData(createdUser.Id)
			}
		})
	}
}

func TestAuth_LogIn_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name               string
		setupUser          bool
		useCorrectPassword bool
		mail               string
		password           string
		expectedErr        bool
	}{
		{
			name:               "successful login with correct credentials",
			setupUser:          true,
			useCorrectPassword: true,
			expectedErr:        false,
		},
		{
			name:               "fail to login with incorrect password",
			setupUser:          true,
			useCorrectPassword: false,
			expectedErr:        true,
		},
		{
			name:        "fail to login with non-existent email",
			setupUser:   false,
			mail:        "nonexistent@example.com",
			password:    "password123",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testUser structs.User
			var userID uuid.UUID

			if tt.setupUser {
				userID, testUser, _ = fixture.createUserForTest()
				defer fixture.cleanupUserData(userID)
			}

			var mail, password string
			if tt.setupUser {
				mail = testUser.Mail
				if tt.useCorrectPassword {
					password = "password123"
				} else {
					password = "wrongpassword"
				}
			} else {
				mail = tt.mail
				password = tt.password
			}

			atoken, rtoken, err := fixture.service.LogIn(fixture.ctx, mail, password)

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

				fixture.authRepo.DeleteToken(fixture.ctx, tokenID)
			}
		})
	}
}

func TestAuth_RefreshToken_AAA(t *testing.T) {
	fixture := NewAuthTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (string, string, uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successful token refresh with existing user",
			setup: func() (string, string, uuid.UUID) {
				userID, testUser, _ := fixture.createUserForTest()

				atoken, rtoken, err := fixture.service.LogIn(fixture.ctx, testUser.Mail, "password123")
				require.NoError(t, err)

				return atoken, rtoken, userID
			},
			expectedErr: false,
		},
		{
			name: "fail to refresh with non-existent user token",
			setup: func() (string, string, uuid.UUID) {
				userID1, testUser1, plainPassword1 := fixture.createUserForTest()
				atoken1, _, err := fixture.service.LogIn(fixture.ctx, testUser1.Mail, plainPassword1)
				require.NoError(t, err)

				userID2, testUser2, plainPassword2 := fixture.createUserForTest()
				_, rtoken2, err := fixture.service.LogIn(fixture.ctx, testUser2.Mail, plainPassword2)
				require.NoError(t, err)

				tokenID, err := fixture.authRepo.VerifyToken(fixture.ctx, rtoken2)
				require.NoError(t, err)
				fixture.authRepo.DeleteToken(fixture.ctx, tokenID)

				defer fixture.cleanupUserData(userID1)
				defer fixture.cleanupUserData(userID2)

				return atoken1, rtoken2, userID2
			},
			expectedErr: true,
		},
		{
			name: "fail to refresh with completely invalid token",
			setup: func() (string, string, uuid.UUID) {
				userID, testUser, plainPassword := fixture.createUserForTest()
				defer fixture.cleanupUserData(userID)

				atoken, _, err := fixture.service.LogIn(fixture.ctx, testUser.Mail, plainPassword)
				require.NoError(t, err)

				return atoken, "completely.invalid.token.12345", userID
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atoken, rtoken, userID := tt.setup()
			if userID != uuid.Nil {
				defer fixture.cleanupUserData(userID)
			}

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
		setup           func() (string, string, uuid.UUID)
		expectedAccess  bool
		expectedRefresh bool
		expectNewToken  bool
	}{
		{
			name: "both tokens valid",
			setup: func() (string, string, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()
				atoken, rtoken, err := fixture.provider.GenToken(fixture.ctx, userID)
				require.NoError(t, err)

				err = fixture.authRepo.CreateToken(fixture.ctx, userID, rtoken)
				require.NoError(t, err)

				return atoken, rtoken, userID
			},
			expectedAccess:  true,
			expectedRefresh: true,
			expectNewToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atoken, rtoken, userID := tt.setup()
			defer fixture.cleanupUserData(userID)

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
		expectedErr bool
	}{
		{
			name: "successfully extract user ID from valid token",
			setup: func() (string, uuid.UUID) {
				userID, _, _ := fixture.createUserForTest()
				atoken, _, err := fixture.provider.GenToken(fixture.ctx, userID)
				require.NoError(t, err)
				return atoken, userID
			},
			expectedErr: false,
		},
		{
			name: "fail to extract user ID from invalid token",
			setup: func() (string, uuid.UUID) {
				return "invalid-token", uuid.Nil
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expectedID := tt.setup()
			if expectedID != uuid.Nil {
				defer fixture.cleanupUserData(expectedID)
			}

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
