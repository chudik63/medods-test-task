package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"medods-test-task/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenLength = 16
)

type Config interface {
	GetAuthJWTSecret() string
	GetAccessTokenExpiration() time.Duration
	GetRefreshTokenExpiration() time.Duration
}

type TokenManager interface {
	NewTokenPair(userID uuid.UUID, IPAddress string) (string, string, error)
	SignToken(claims Claims) (string, error)
	ParseJWT(token string) (*Claims, error)
	ParseRefreshToken(refreshToken string) (uuid.UUID, error)
	HashToken(password string) (string, error)
	ValidateToken(token, hashedToken string) error
	GetAccessTTL() time.Duration
	GetRefreshTTL() time.Duration
}

type Claims struct {
	UserID    uuid.UUID
	IPAddress string
	Subject   string
	jwt.StandardClaims
}

type Manager struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		secret:     cfg.GetAuthJWTSecret(),
		accessTTL:  cfg.GetAccessTokenExpiration(),
		refreshTTL: cfg.GetRefreshTokenExpiration(),
	}
}

func (m *Manager) NewTokenPair(userID uuid.UUID, IPAddress string) (string, string, error) {
	accessClaims := Claims{
		UserID:    userID,
		IPAddress: IPAddress,
		Subject:   "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken, err := m.SignToken(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("could not sign token: %w", err)
	}

	refreshToken := make([]byte, refreshTokenLength)
	_, err = rand.Read(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	refreshToken = append(refreshToken, userID[:]...)

	return accessToken, base64.URLEncoding.EncodeToString(refreshToken), nil
}

func (m *Manager) SignToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(m.secret))
}

func (m *Manager) ParseJWT(accessToken string) (*Claims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, models.ErrInvalidToken
	}

	return claims, nil
}

func (m *Manager) ParseRefreshToken(refreshToken string) (uuid.UUID, error) {
	decoded, err := base64.URLEncoding.DecodeString(refreshToken)
	if err != nil {
		return uuid.UUID{}, err
	}

	if len(decoded) != refreshTokenLength+len(uuid.UUID{}) {
		return uuid.UUID{}, models.ErrInvalidToken
	}

	userID, err := uuid.FromBytes(decoded[refreshTokenLength:])
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to retrieve user uuid from refresh token: %w", err)
	}

	return userID, nil
}

func (m *Manager) HashToken(token string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashed), nil
}

func (m *Manager) ValidateToken(token, hashedToken string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return models.ErrMismatchedHashAndToken
	}

	return nil
}

func (m *Manager) GetAccessTTL() time.Duration {
	return m.accessTTL
}

func (m *Manager) GetRefreshTTL() time.Duration {
	return m.refreshTTL
}
