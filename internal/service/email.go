package service

import (
	"context"
	"medods-test-task/pkg/email/smtp"
	"medods-test-task/pkg/logger"

	"go.uber.org/zap"
)

type SMTPSender interface {
	Send(input smtp.SendEmailInput) error
}

type emailService struct {
	sender          SMTPSender
	warningSubject  string
	warningTemplate string

	domain string
	logger logger.Logger
}

func NewEmailService(s SMTPSender, wsub, wtempl, dom string) *emailService {
	return &emailService{
		sender:          s,
		warningSubject:  wsub,
		warningTemplate: wtempl,
		domain:          dom,
	}
}

func (s *emailService) SendWarningEmail(ctx context.Context, email string) {
	sendInput := smtp.SendEmailInput{Subject: s.warningSubject, To: email}

	if err := sendInput.GenerateBodyFromHTML(s.warningTemplate, nil); err != nil {
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
