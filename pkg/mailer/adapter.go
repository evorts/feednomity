package mailer

import "context"

type manager struct {
	senderName   string
	senderEmail  string
	replyToName  string
	replyToEmail string
	apiUrl       string // mail provider api url
}

type Target struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type IMailer interface {
	SetSender(name, email string) IMailer
	SetReplyTo(name, email string) IMailer
	SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) ([]byte, error)
	call(payload []byte, args map[string]interface{}) ([]byte, error)
}
