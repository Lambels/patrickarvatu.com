package pa

// EmailService represents a service which manages auth in the system.
type EmailService interface {
	// SendEmail will send a emails to the adresses provided in to.
	SendEmail(to []string, body, subject string) error
}
