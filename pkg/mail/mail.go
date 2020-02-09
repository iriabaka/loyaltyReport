package mail

import (
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
	"net/smtp"
)

type Server struct {
	Host     string
	Port     string
	Username string
	Password string
}

func Send(server Server, from string, to []string, subject string, text string, attachFiles []string) error {
	mail := email.NewEmail()
	mail.From = from
	mail.To = to
	mail.Subject = subject
	mail.Text = []byte(text)

	for _, file := range attachFiles {
		if _, err := mail.AttachFile(file); err != nil {
			return errors.WithStack(err)
		}
	}

	if server.Username != "" {
		auth := smtp.PlainAuth("", server.Username, server.Password, server.Host)

		if err := mail.Send(server.Host+":"+server.Port, auth); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if err := mail.Send(server.Host+":"+server.Port, nil); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
