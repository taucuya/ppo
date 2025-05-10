package auth

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// 	"github.com/taucuya/ppo/internal/core/mock_structs"
// 	"github.com/taucuya/ppo/internal/core/structs"
// )

// var testError = errors.New("test error")

// func TestNew(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockProvider := mock_structs.NewMockAuthProvider(ctrl)
// 	mockRepo := mock_structs.NewMockAuthRepository(ctrl)
// 	mockUser := mock_structs.NewMockAuthUser(ctrl)

// 	service := New(mockProvider, mockRepo, mockUser)

// 	if service == nil {
// 		t.Error("Expected service instance, got nil")
// 	}
// }

// func TestSignIn(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUser := mock_structs.NewMockAuthUser(ctrl)
// 	service := New(nil, nil, mockUser)

// 	user := structs.User{
// 		Id:   structs.GenId(),
// 		Mail: "test@example.com",
// 	}

// 	tests := []struct {
// 		name    string
// 		user    structs.User
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful sign-in",
// 			user: user,
// 			mock: func() {
// 				mockUser.EXPECT().Create(gomock.Any(), user).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "repository error",
// 			user: user,
// 			mock: func() {
// 				mockUser.EXPECT().Create(gomock.Any(), user).Return(testError)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := service.SignIn(context.Background(), tt.user)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestLogIn(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockProvider := mock_structs.NewMockAuthProvider(ctrl)
// 	mockRepo := mock_structs.NewMockAuthRepository(ctrl)
// 	mockUser := mock_structs.NewMockAuthUser(ctrl)
// 	service := New(mockProvider, mockRepo, mockUser)

// 	user := structs.User{
// 		Id:   structs.GenId(),
// 		Mail: "test@example.com",
// 	}
// 	atoken, rtoken := "access-token", "refresh-token"

// 	tests := []struct {
// 		name       string
// 		mail       string
// 		password   string
// 		mock       func()
// 		wantAToken string
// 		wantRToken string
// 		wantErr    bool
// 	}{
// 		{
// 			name:     "successful login",
// 			mail:     user.Mail,
// 			password: "password123",
// 			mock: func() {
// 				mockUser.EXPECT().GetByMail(gomock.Any(), user.Mail).Return(user, nil)
// 				mockProvider.EXPECT().GenToken(gomock.Any(), user.Id).Return(atoken, rtoken, nil)
// 				mockRepo.EXPECT().CreateToken(gomock.Any(), user.Id, rtoken).Return(nil)
// 			},
// 			wantAToken: atoken,
// 			wantRToken: rtoken,
// 			wantErr:    false,
// 		},
// 		{
// 			name:     "user not found",
// 			mail:     user.Mail,
// 			password: "wrongpassword",
// 			mock: func() {
// 				mockUser.EXPECT().GetByMail(gomock.Any(), user.Mail).Return(structs.User{}, testError)
// 			},
// 			wantAToken: "",
// 			wantRToken: "",
// 			wantErr:    true,
// 		},
// 		{
// 			name:     "token generation error",
// 			mail:     user.Mail,
// 			password: "password123",
// 			mock: func() {
// 				mockUser.EXPECT().GetByMail(gomock.Any(), user.Mail).Return(user, nil)
// 				mockProvider.EXPECT().GenToken(gomock.Any(), user.Id).Return("", "", testError)
// 			},
// 			wantAToken: "",
// 			wantRToken: "",
// 			wantErr:    true,
// 		},
// 		{
// 			name:     "repository error",
// 			mail:     user.Mail,
// 			password: "password123",
// 			mock: func() {
// 				mockUser.EXPECT().GetByMail(gomock.Any(), user.Mail).Return(user, nil)
// 				mockProvider.EXPECT().GenToken(gomock.Any(), user.Id).Return(atoken, rtoken, nil)
// 				mockRepo.EXPECT().CreateToken(gomock.Any(), user.Id, rtoken).Return(testError)
// 			},
// 			wantAToken: "",
// 			wantRToken: "",
// 			wantErr:    true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			gotAToken, gotRToken, err := service.LogIn(context.Background(), tt.mail, tt.password)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("LogIn() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if gotAToken != tt.wantAToken || gotRToken != tt.wantRToken {
// 				t.Errorf("LogIn() = %v, %v, want %v, %v", gotAToken, gotRToken, tt.wantAToken, tt.wantRToken)
// 			}
// 		})
// 	}
// }

// func TestLogOut(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockProvider := mock_structs.NewMockAuthProvider(ctrl)
// 	mockRepo := mock_structs.NewMockAuthRepository(ctrl)
// 	service := New(mockProvider, mockRepo, nil)

// 	rtoken := "refresh-token"
// 	userID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		rtoken  string
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name:   "successful logout",
// 			rtoken: rtoken,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(false, nil)
// 				mockRepo.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(userID, nil)
// 				mockRepo.EXPECT().DeleteToken(gomock.Any(), userID).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:   "token verification provider error",
// 			rtoken: rtoken,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(false, testError)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "token verification repo error",
// 			rtoken: rtoken,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(false, nil)
// 				mockRepo.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(uuid.UUID{}, testError)
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name:   "repository delete error",
// 			rtoken: rtoken,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(false, nil)
// 				mockRepo.EXPECT().VerifyToken(gomock.Any(), rtoken).Return(userID, nil)
// 				mockRepo.EXPECT().DeleteToken(gomock.Any(), userID).Return(testError)
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := service.LogOut(context.Background(), tt.rtoken)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("LogOut() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.wantErr && !errors.Is(err, testError) {
// 				t.Errorf("Expected testError, got %v", err)
// 			}
// 		})
// 	}
// }

// func TestVerifyAToken(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockProvider := mock_structs.NewMockAuthProvider(ctrl)
// 	svc := New(mockProvider, nil, nil)

// 	token := "access-token"
// 	tests := []struct {
// 		name    string
// 		token   string
// 		mock    func()
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name:  "valid token",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(true, nil)
// 			},
// 			want:    true,
// 			wantErr: false,
// 		},
// 		{
// 			name:  "invalid token",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(false, nil)
// 			},
// 			want:    false,
// 			wantErr: false,
// 		},
// 		{
// 			name:  "token error",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(false, testError)
// 			},
// 			want:    false,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := svc.VerifyAToken(context.Background(), tt.token)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("VerifyAToken() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if got != tt.want {
// 				t.Errorf("VerifyAToken() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestVerifyRToken(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockProvider := mock_structs.NewMockAuthProvider(ctrl)
// 	mockRepo := mock_structs.NewMockAuthRepository(ctrl)
// 	svc := New(mockProvider, mockRepo, nil)

// 	token := "refresh-token"
// 	userID := structs.GenId()

// 	tests := []struct {
// 		name    string
// 		token   string
// 		mock    func()
// 		wantID  uuid.UUID
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name:  "valid token",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(false, nil)
// 				mockRepo.EXPECT().VerifyToken(gomock.Any(), token).Return(userID, nil)
// 			},
// 			wantID:  userID,
// 			want:    false,
// 			wantErr: false,
// 		},
// 		{
// 			name:  "expired token",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(true, nil)
// 				mockRepo.EXPECT().VerifyToken(gomock.Any(), token).Return(userID, nil)
// 			},
// 			wantID:  userID,
// 			want:    true,
// 			wantErr: false,
// 		},
// 		{
// 			name:  "provider error",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(false, testError)
// 			},
// 			wantID:  uuid.UUID{},
// 			want:    false,
// 			wantErr: true,
// 		},
// 		{
// 			name:  "repo error",
// 			token: token,
// 			mock: func() {
// 				mockProvider.EXPECT().VerifyToken(gomock.Any(), token).Return(false, nil)
// 				mockRepo.EXPECT().VerifyToken(gomock.Any(), token).Return(uuid.UUID{}, testError)
// 			},
// 			wantID:  uuid.UUID{},
// 			want:    false,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			gotID, got, err := svc.VerifyRToken(context.Background(), tt.token)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("VerifyRToken() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if gotID != tt.wantID || got != tt.want {
// 				t.Errorf("VerifyRToken() = %v, %v, want %v, %v", gotID, got, tt.wantID, tt.want)
// 			}
// 		})
// 	}
// }

// func TestRefreshToken(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockProvider := mock_structs.NewMockAuthProvider(ctrl)
// 	service := New(mockProvider, nil, nil)

// 	atoken := "access-token"
// 	rtoken := "refresh-token"
// 	newToken := "new-access-token"

// 	tests := []struct {
// 		name    string
// 		atoken  string
// 		rtoken  string
// 		mock    func()
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "successful refresh",
// 			atoken: atoken,
// 			rtoken: rtoken,
// 			mock: func() {
// 				mockProvider.EXPECT().RefreshToken(gomock.Any(), atoken, rtoken).Return(newToken, nil)
// 			},
// 			want:    newToken,
// 			wantErr: false,
// 		},
// 		{
// 			name:   "refresh error",
// 			atoken: atoken,
// 			rtoken: rtoken,
// 			mock: func() {
// 				mockProvider.EXPECT().RefreshToken(gomock.Any(), atoken, rtoken).Return("", testError)
// 			},
// 			want:    "",
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := service.RefreshToken(context.Background(), tt.atoken, tt.rtoken)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if got != tt.want {
// 				t.Errorf("RefreshToken() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
