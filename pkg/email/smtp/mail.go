package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"medods-test-task/internal/models"
	"medods-test-task/pkg/email"
)

type SendEmailInput struct {
	To      string
	Subject string
	Body    string
}

func (e *SendEmailInput) GenerateBodyFromHTML(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return fmt.Errorf("failed to parse file %s:%w", templateFileName, err)
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	e.Body = buf.String()

	return nil
}

func (e *SendEmailInput) Validate() error {
	if e.To == "" {
		return models.ErrSMTPEmptyTo
	}

	if e.Subject == "" || e.Body == "" {
		return models.ErrSMTPEmptyMail
	}

	if !email.IsValid(e.To) {
		return models.ErrSMTPInvalidToEmail
	}

	return nil
}
