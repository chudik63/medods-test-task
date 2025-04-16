package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Config interface {
	GetAuthJWTSecret() string
	GetAccessTokenExpiration() time.Duration
	GetRefreshTokenExpiration() time.Duration
}

type Claims struct {
	UserID    uint64
	IPAddress string
	jwt.StandardClaims
}

type TokenManager interface {
	NewJWT(userId uint64, email, ipAddress, tokenID, role string) (string, string, error)
	SignToken(claims Claims) (string, error)
	ParseJWT(accessToken string) (*Claims, error)
	HashToken(password string) (string, error)
	ValidateToken(password, hashedPassword string) error
	GetAccessTTL() time.Duration
	GetRefreshTTL() time.Duration
}

type Manager struct {
	secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		secret:     cfg.GetAuthJWTSecret(),
		AccessTTL:  cfg.GetAccessTokenExpiration(),
		RefreshTTL: cfg.GetRefreshTokenExpiration(),
	}
}

func (m *Manager) NewJWT(userId uint64, email, ipAddress, tokenID, role string) (string, string, error) {
	accessClaims := Claims{
		UserID:    userId,
		IPAddress: ipAddress,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.AccessTTL).Unix(),
		},
	}

	accessToken, err := m.SignToken(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("could not sign token: %w", err)
	}

	refreshClaims := Claims{
		UserID:    userId,
		IPAddress: ipAddress,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.RefreshTTL).Unix(),
		},
	}

	refreshToken, err := m.SignToken(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("could not sign token: %w", err)
	}

	return accessToken, refreshToken, nil

}

func (m *Manager) SignToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(m.secret)
}

func (m *Manager) ParseJWT(accessToken string) (*Claims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (m *Manager) HashToken(token string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashed), nil
}

func (m *Manager) ValidateToken(token, hashedToken string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
}

func (m *Manager) GetAccessTTL() time.Duration {
	return m.AccessTTL
}

func (m *Manager) GetRefreshTTL() time.Duration {
	return m.RefreshTTL
}
