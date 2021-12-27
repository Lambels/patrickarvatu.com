package pa

type EmailService interface {
	SendEmail(addr []byte, body, subject string) error
}
