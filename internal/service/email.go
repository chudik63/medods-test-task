package service

import (
	"context"
	"medods-test-task/config"
	"medods-test-task/pkg/email/smtp"
	"medods-test-task/pkg/logger"

	"go.uber.org/zap"
)

type SMTPSender interface {
	Send(input smtp.SendEmailInput) error
}

type emailService struct {
	sender      SMTPSender
	emailConfig *config.EmailConfig
	logger      logger.Logger
}

func NewEmailService(s SMTPSender, logger logger.Logger, emailConf *config.EmailConfig) *emailService {
	return &emailService{
		sender:      s,
		emailConfig: emailConf,
		logger:      logger,
	}
}

func (s *emailService) SendIPWarningEmail(ctx context.Context, email string) {
	sendInput := smtp.SendEmailInput{Subject: s.emailConfig.IPWarningSubject, To: email}

	if err := sendInput.GenerateBodyFromHTML(s.emailConfig.IPWarningTemplate, nil); err != nil {
		s.logger.Error(ctx, "failed generate body from html template", zap.Error(err))

		return
	}

	err := s.sender.Send(sendInput)
	if err != nil {
		s.logger.Error(ctx, "failed send verification email", zap.Error(err))

		return
	}

	s.logger.Debug(ctx, "warning message sent", zap.String("to", email))
}
