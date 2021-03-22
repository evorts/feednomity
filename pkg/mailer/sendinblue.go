package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type sibManager struct {
	apiKey string
	manager
}

type sibPayload struct {
	Sender      Target   `json:"sender"`
	To          []Target `json:"to"`
	Subject     string   `json:"subject"`
	HtmlContent string   `json:"htmlContent"`
}

func NewSendInBlue(url, apiKey string) IMailer {
	return &sibManager{
		apiKey: apiKey,
		manager: manager{
			apiUrl: url,
		},
	}
}

func (s *sibManager) SetSender(sender, email string) IMailer {
	s.senderName = sender
	s.senderEmail = email
	return s
}

func (s *sibManager) SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) error {
	if err := validate(to, subject, html); err != nil {
		return err
	}
	payload := sibPayload{
		Sender: Target{
			Name:  s.senderName,
			Email: s.senderEmail,
		},
		To:          to,
		Subject:     subject,
		HtmlContent: bindDataToTemplate(data, html),
	}
	if args, err := json.Marshal(payload); err == nil {
		_, err := s.call(args)
		return err
	}
	return errors.New("serializing payload failed")
}

func (s *sibManager) call(payload []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, s.apiUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}
