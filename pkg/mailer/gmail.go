package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type gmailManager struct {
	user, pass, address string
	manager
}

func NewGmail(user, pass, address string) IMailer {
	return &gmailManager{
		user:    user,
		pass:    pass,
		address: address,
	}
}

func (g *gmailManager) SetSender(name, email string) IMailer {
	g.senderName = name
	g.senderEmail = email
	return g
}

func (g *gmailManager) SetReplyTo(name, email string) IMailer {
	g.replyToName = name
	g.replyToEmail = email
	return g
}

func (g *gmailManager) SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) ([]byte, error) {
	if err := validate(to, subject, html); err != nil {
		return nil, err
	}
	message := make([]string, 0)
	tos := make([]string, 0)
	for _, v := range to {
		tos = append(tos, fmt.Sprintf("%s", v.Email))
	}
	message = append(message, fmt.Sprintf("From: %s", g.senderEmail))
	message = append(message, fmt.Sprintf("To: %s", strings.Join(tos, ",")))
	message = append(message, fmt.Sprintf("Subject: %s", subject))
	message = append(message, fmt.Sprintf("MIME Version: 1.0; Content-Type: text/html; charset=utf-8;\n"))
	message = append(message, bindDataToTemplate(data, html))
	return g.call([]byte(strings.Join(message, "\n")), map[string]interface{}{
		"dest": tos,
	})
}

func (g *gmailManager) call(payload []byte, args map[string]interface{}) ([]byte, error) {
	adr := strings.Split(g.address, ":")
	host := adr[0]
	err := smtp.SendMail(g.address,
		smtp.PlainAuth("", g.user, g.pass, host),
		g.user, args["dest"].([]string), payload,
	)
	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return nil, err
	}
	return nil, nil
}
