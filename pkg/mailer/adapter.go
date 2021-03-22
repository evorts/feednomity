package mailer

import "context"

type manager struct {
	senderName string
	senderEmail string
	apiUrl string // mail provider api url
}

type Target struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

type IMailer interface {
	SetSender(sender, email string) IMailer
	SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) error
	call(payload []byte) ([]byte, error)
}
