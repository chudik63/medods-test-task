package http

import (
	"context"
	"medods-test-task/pkg/logger"
)

type AuthService interface {
	NewSession(ctx context.Context, userID, IPAddress string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken, IPAdress string) (string, string, error)
}

type AppController struct {
	serv   AuthService
	logger logger.Logger
}

func NewAppController(serv AuthService, logger logger.Logger) *AppController {
	return &AppController{
		serv:   serv,
		logger: logger,
	}
}
