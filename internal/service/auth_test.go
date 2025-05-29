package service

import (
	"context"
	"medods-test-task/internal/models"
	"medods-test-task/internal/service/mocks"
	"medods-test-task/pkg/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_RefreshToken(t *testing.T) {
	manager := &utils.Manager{}
	userID := uuid.New()
	ip := "127.0.0.1"
	email := "test@email.com"

	_, refresh, _ := manager.NewTokenPair(userID, ip)
	hashed, _ := manager.HashToken(refresh)

	type (
		repoMockBehavior  func(r *mocks.AuthRepo, userID uuid.UUID, ip, hashedToken, newHashedToken string)
		tokenMockBehavior func(t *mocks.TokenManager, userID uuid.UUID, token, hashedToken, newAccessToken, newRefreshToken, newHashedToken string)
		args              struct {
			ctx          context.Context
			refreshToken string
			IPAddress    string
		}
	)

	tests := []struct {
		name            string
		args            args
		userID          uuid.UUID
		newAccessToken  string
		newRefreshToken string
		hashedToken     string
		newHashedToken  string
		repoMock        repoMockBehavior
		tokenMock       tokenMockBehavior
		expectedErr     error
	}{
		{
			name:            "OK",
			userID:          userID,
			newAccessToken:  "access",
			newRefreshToken: "refresh",
			hashedToken:     hashed,
			newHashedToken:  "hashed",
			expectedErr:     nil,
			args: args{
				ctx:          context.Background(),
				refreshToken: refresh,
				IPAddress:    ip,
			},
			repoMock: func(r *mocks.AuthRepo, userID uuid.UUID, ip, hashedToken, newHashedToken string) {
				r.On("GetSessionByUserID", mock.Anything, userID).Return(&models.RefreshSession{
					UserID:    userID,
					IP:        ip,
					Token:     hashedToken,
					ExpiresAt: time.Now().Add(720 * time.Hour),
				}, nil)
				r.On("DeleteSessionByUserID", mock.Anything, userID).Return(nil)
				r.On("CreateSession", mock.Anything, mock.Anything).Return(nil)
			},
			tokenMock: func(m *mocks.TokenManager, userID uuid.UUID, token, hashedToken, newAccessToken, newRefreshToken, newHashedToken string) {
				m.On("ParseRefreshToken", token).Return(userID, nil)
				m.On("ValidateToken", token, hashedToken).Return(nil)
				m.On("NewTokenPair", userID, mock.Anything).Return(newAccessToken, newRefreshToken, nil)
				m.On("HashToken", newRefreshToken).Return(newHashedToken, nil)
				m.On("GetRefreshTTL").Return(time.Duration(720 * time.Hour))
			},
		},
		{
			name:            "Invalid Token",
			userID:          userID,
			newAccessToken:  "access",
			newRefreshToken: "refresh",
			hashedToken:     "hashedinvalid",
			newHashedToken:  "hashed",
			expectedErr:     models.ErrInvalidToken,
			args: args{
				ctx:          context.Background(),
				refreshToken: "inValId-Tokn",
				IPAddress:    ip,
			},
			repoMock: func(r *mocks.AuthRepo, userID uuid.UUID, ip, hashedToken, newHashedToken string) {

			},
			tokenMock: func(m *mocks.TokenManager, userID uuid.UUID, token, hashedToken, newAccessToken, newRefreshToken, newHashedToken string) {
				m.On("ParseRefreshToken", token).Return(uuid.UUID{}, models.ErrInvalidToken)
			},
		},
		{
			name:            "Token expired",
			userID:          userID,
			newAccessToken:  "access",
			newRefreshToken: "refresh",
			hashedToken:     "hashedinvalid",
			newHashedToken:  "hashed",
			expectedErr:     models.ErrTokenExpired,
			args: args{
				ctx:          context.Background(),
				refreshToken: "inValId-Tokn",
				IPAddress:    ip,
			},
			repoMock: func(r *mocks.AuthRepo, userID uuid.UUID, ip, hashedToken, newHashedToken string) {
				r.On("GetSessionByUserID", mock.Anything, userID).Return(&models.RefreshSession{
					UserID:    userID,
					IP:        ip,
					Token:     hashedToken,
					ExpiresAt: time.Now().Add(-24 * time.Hour),
				}, nil)
				r.On("DeleteSessionByUserID", mock.Anything, userID).Return(nil)
			},
			tokenMock: func(m *mocks.TokenManager, userID uuid.UUID, token, hashedToken, newAccessToken, newRefreshToken, newHashedToken string) {
				m.On("ParseRefreshToken", token).Return(userID, nil)
				m.On("ValidateToken", token, hashedToken).Return(nil)
			},
		},
		{
			name:            "Wrong IP",
			userID:          userID,
			newAccessToken:  "access",
			newRefreshToken: "refresh",
			hashedToken:     "hashedinvalid",
			newHashedToken:  "hashed",
			expectedErr:     models.ErrInvalidSession,
			args: args{
				ctx:          context.Background(),
				refreshToken: "inValId-Tokn",
				IPAddress:    ip,
			},
			repoMock: func(r *mocks.AuthRepo, userID uuid.UUID, ip, hashedToken, newHashedToken string) {
				r.On("GetSessionByUserID", mock.Anything, userID).Return(&models.RefreshSession{
					UserID:    userID,
					IP:        "127.1.0.1",
					Token:     hashedToken,
					ExpiresAt: time.Now().Add(720 * time.Hour),
				}, nil)
				r.On("DeleteSessionByUserID", mock.Anything, userID).Return(nil)
				r.On("GetUserByID", mock.Anything, userID).Return(&models.User{
					ID:    userID,
					Email: email,
				}, nil)
			},
			tokenMock: func(m *mocks.TokenManager, userID uuid.UUID, token, hashedToken, newAccessToken, newRefreshToken, newHashedToken string) {
				m.On("ParseRefreshToken", token).Return(userID, nil)
				m.On("ValidateToken", token, hashedToken).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewAuthRepo(t)
			m := mocks.NewTokenManager(t)
			e := mocks.NewEmailService(t)

			s := &AuthService{
				authRepo:     r,
				tokenManager: m,
				emailService: e,
			}

			tt.repoMock(r, tt.userID, tt.args.IPAddress, tt.hashedToken, tt.newHashedToken)
			tt.tokenMock(m, tt.userID, tt.args.refreshToken, tt.hashedToken, tt.newAccessToken, tt.newRefreshToken, tt.newHashedToken)

			_, _, err := s.RefreshToken(tt.args.ctx, tt.args.refreshToken, tt.args.IPAddress)
			if err != tt.expectedErr {
				t.Errorf("error = %v, expectedError %v", err, tt.expectedErr)
				return
			}
		})
	}
}
