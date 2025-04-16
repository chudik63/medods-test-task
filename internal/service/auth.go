package service

import (
	"context"
	"errors"
	"medods-test-task/internal/models"
	"medods-test-task/pkg/utils"
	"time"
)

type AuthRepo interface {
	CreateSession(ctx context.Context, session *models.RefreshSession) error
	DeleteSessionByUserID(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type AuthService struct {
	authRepo     AuthRepo
	tokenManager utils.TokenManager
}

func NewAuthService(auth AuthRepo, token utils.TokenManager) *AuthService {
	return &AuthService{
		authRepo:     auth,
		tokenManager: token,
	}
}

func (s *AuthService) NewSession(ctx context.Context, userID, IPAddress string) (string, string, error) {
	if userID == "" {
		return "", "", models.ErrEmptyUserID
	}

	access, refresh, err := s.tokenManager.NewJWT(userID, IPAddress)
	if err != nil {
		return "", "", err
	}

	_, err = s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, models.ErrUserNotFound) {
			return "", "", err
		}

		err := s.authRepo.CreateUser(ctx, &models.User{
			ID: userID,
		})
		if err != nil {
			return "", "", err
		}
	}

	err = s.authRepo.DeleteSessionByUserID(ctx, userID)
	if err != nil && !errors.Is(err, models.ErrSessionNotFound) {
		return "", "", err
	}

	err = s.authRepo.CreateSession(ctx, &models.RefreshSession{
		UserID:    userID,
		Token:     refresh,
		IP:        IPAddress,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.tokenManager.GetRefreshTTL()),
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, err
}
