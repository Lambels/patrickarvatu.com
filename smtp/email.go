package smtp

import (
	"net/smtp"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/jordan-wright/email"
)

var _ pa.EmailService = (*EmailService)(nil)

type EmailService struct {
	auth smtp.Auth
	addr string
}

func NewEmailService(addr, identity, username, password, host string) *EmailService {
	auth := smtp.PlainAuth(identity, username, password, host)

	return &EmailService{
		auth: auth,
		addr: addr,
	}
}

func (e *EmailService) SendEmail(to []string, body, subject string) error {
	email := email.NewEmail()
	email.Subject = subject
	email.To = to
	email.Text = []byte(body)

	return email.Send(e.addr, e.auth)
}
