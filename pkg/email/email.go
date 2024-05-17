package email

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"regexp"
	"text/template"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=email.go -destination=mocks/mock.go

type SMTPSender struct {
	Email string
	Pass  string
	Host  string
	Port  string
	auth  smtp.Auth
}

type EmailSender interface {
	Send(templatePath string, emailTo string, subject string,
		data interface{}) error
}

var (
	ErrInvalidEmail = errors.New("email is invalid")
)

type EmailSenderConfig struct {
	Email string
	Pass  string
	Host  string
	Port  string
}

func NewEmailSender(cfg EmailSenderConfig) (*SMTPSender, error) {
	if err := validateEmail(cfg.Email); err != nil {
		return nil, ErrInvalidEmail
	}
	auth := smtp.PlainAuth("", cfg.Email, cfg.Pass, cfg.Host)
	return &SMTPSender{
		Email: cfg.Email,
		Pass:  cfg.Pass,
		Host:  cfg.Host,
		Port:  cfg.Port,
		auth:  auth,
	}, nil
}

func validateEmail(email string) error {
	if len(email) < 3 || len(email) > 1024 {
		return ErrInvalidEmail
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		email)
	if !matched {
		return ErrInvalidEmail
	}

	return nil
}

func (e *SMTPSender) Send(templatePath string, emailTo string, subject string,
	data interface{}) error {

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		logrus.Error("error parse html template: %w", err)
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		logrus.Error("error execute html template: %w", err)
		return err
	}
	msg := fmt.Sprintf("To: %s\nSubject: %s\n", emailTo, subject) + mime + buf.String()

	to := []string{
		emailTo,
	}
	err = smtp.SendMail(e.Host+":"+e.Port, e.auth, e.Email, to, []byte(msg))
	if err != nil {
		logrus.Error("error send email message: %w", err)
		return err
	}

	return nil
}
