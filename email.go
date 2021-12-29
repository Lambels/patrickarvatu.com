package pa

type EmailService interface {
	SendEmail(to []string, body, subject string) error
}
