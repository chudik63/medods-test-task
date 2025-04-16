package http

import (
	"medods-test-task/pkg/logger"
)

type AuthService interface {
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
