package service

import (
	"context"
	"errors"
	"log"
	"medods-test-task/internal/models"
	"time"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name EmailService
type EmailService interface {
	SendIPWarningEmail(ctx context.Context, email string)
}

//go:generate go run github.com/vektra/mockery/v2@latest --name TokenManager
type TokenManager interface {
	NewTokenPair(userID uuid.UUID, IPAddress string) (string, string, error)
	ParseRefreshToken(refreshToken string) (uuid.UUID, error)
	HashToken(password string) (string, error)
	ValidateToken(token, hashedToken string) error
	GetRefreshTTL() time.Duration
}

//go:generate go run github.com/vektra/mockery/v2@latest --name AuthRepo
type AuthRepo interface {
	CreateSession(ctx context.Context, session *models.RefreshSession) error
	DeleteSessionByUserID(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*models.RefreshSession, error)
}

type AuthService struct {
	authRepo     AuthRepo
	tokenManager TokenManager
	emailService EmailService
}

func NewAuthService(auth AuthRepo, token TokenManager, email EmailService) *AuthService {
	return &AuthService{
		authRepo:     auth,
		tokenManager: token,
		emailService: email,
	}
}

func (s *AuthService) NewSession(ctx context.Context, userID, IPAddress string) (string, string, error) {
	if userID == "" {
		return "", "", models.ErrEmptyUserID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", "", models.ErrInvalidUserID
	}

	access, refresh, err := s.tokenManager.NewTokenPair(userUUID, IPAddress)
	if err != nil {
		return "", "", err
	}

	_, err = s.authRepo.GetUserByID(ctx, userUUID)
	if err != nil {
		if !errors.Is(err, models.ErrUserNotFound) {
			return "", "", err
		}

		err := s.authRepo.CreateUser(ctx, &models.User{
			ID: userUUID,
		})
		if err != nil {
			return "", "", err
		}
	}

	err = s.authRepo.DeleteSessionByUserID(ctx, userUUID)
	if err != nil && !errors.Is(err, models.ErrSessionNotFound) {
		return "", "", err
	}

	hashedRefresh, err := s.tokenManager.HashToken(refresh)
	if err != nil {
		return "", "", err
	}

	err = s.authRepo.CreateSession(ctx, &models.RefreshSession{
		UserID:    userUUID,
		Token:     hashedRefresh,
		IP:        IPAddress,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.tokenManager.GetRefreshTTL()),
	})
	if err != nil {
		return "", "", err
	}

	log.Print(hashedRefresh)
	log.Print(refresh)

	return access, refresh, err
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, IPAddress string) (string, string, error) {
	userID, err := s.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	session, err := s.authRepo.GetSessionByUserID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	log.Print(session.Token)
	log.Print(refreshToken)

	err = s.tokenManager.ValidateToken(refreshToken, session.Token)
	if err != nil {
		return "", "", err
	}

	err = s.authRepo.DeleteSessionByUserID(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return "", "", models.ErrTokenExpired
	}

	if session.IP != IPAddress {
		user, err := s.authRepo.GetUserByID(ctx, session.UserID)
		if err != nil {
			return "", "", err
		}

		go s.emailService.SendIPWarningEmail(ctx, user.Email)

		return "", "", models.ErrInvalidSession
	}

	accessToken, newRefreshToken, err := s.tokenManager.NewTokenPair(session.UserID, IPAddress)
	if err != nil {
		return "", "", err
	}

	hashedRefresh, err := s.tokenManager.HashToken(newRefreshToken)
	if err != nil {
		return "", "", err
	}

	err = s.authRepo.CreateSession(ctx, &models.RefreshSession{
		UserID:    session.UserID,
		Token:     hashedRefresh,
		IP:        IPAddress,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.tokenManager.GetRefreshTTL()),
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
