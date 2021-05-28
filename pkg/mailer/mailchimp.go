package mailer

import "context"

type mcManager struct {
	apiKey string
	manager
}

type mcPayload struct {

}

func NewMailChimp(apiKey string) IMailer {
	return &mcManager{ apiKey: apiKey }
}

func (m *mcManager) SetSender(name, email string) IMailer {
	m.senderName = name
	m.senderEmail = email
	return m
}

func (m *mcManager) SetReplyTo(name, email string) IMailer {
	m.replyToName = name
	m.replyToEmail = email
	return m
}

func (m *mcManager) SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) ([]byte, error) {
	if err := validate(to, subject, html); err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *mcManager) call(payload []byte, args map[string]interface{}) ([]byte, error) {
	panic("implement me")
}

