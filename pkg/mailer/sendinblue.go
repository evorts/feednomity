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
	Name        string   `json:"name"`
	Sender      Target   `json:"sender"`
	ReplyTo     Target   `json:"reply_to"`
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

func (s *sibManager) SetReplyTo(name, email string) IMailer {
	s.replyToName = email
	s.replyToEmail = email
	return s
}

func (s *sibManager) SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) ([]byte, error) {
	if err := validate(to, subject, html); err != nil {
		return nil, err
	}
	payload := sibPayload{
		Name: "Mail Sender",
		Sender: Target{
			Name:  s.senderName,
			Email: s.senderEmail,
		},
		ReplyTo: Target{
			Name:  s.replyToName,
			Email: s.replyToEmail,
		},
		To:          to,
		Subject:     subject,
		HtmlContent: bindDataToTemplate(data, html),
	}
	if args, err := json.Marshal(payload); err == nil {
		body, err2 := s.call(args)
		return body, err2
	}
	return nil, errors.New("serializing payload failed")
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
	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		err = errors.New("http status are not 2xx, thus something must have been wrong!")
	}
	return body, err
}
