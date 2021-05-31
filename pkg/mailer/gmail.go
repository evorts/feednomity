package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

type gmailManager struct {
	user, pass, address string
	tpl                 *template.Template
	manager
}

func NewGmail(user, pass, address string) IMailer {
	return &gmailManager{
		user:    user,
		pass:    pass,
		address: address,
		tpl:     template.New("gmail_template"),
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

func (g *gmailManager) SendHtml(ctx context.Context, targets []Target, subject, html string, data map[string]string) ([]byte, error) {
	if err := validate(targets, subject, html); err != nil {
		return nil, err
	}
	tos := make([]string, 0)
	for _, v := range targets {
		tos = append(tos, fmt.Sprintf("%s", v.Email))
	}
	from := fmt.Sprintf("From: %s\n", g.senderEmail)
	to := fmt.Sprintf("To: %s\n", strings.Join(tos, ","))
	subject = fmt.Sprintf("Subject: %s\n", subject)
	mime := fmt.Sprintf("MIME Version: 1.0; \nContent-Type: text/html; charset=utf-8;\n\n")
	t, err := g.tpl.Parse(bindDataToTemplate(data, html))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return g.call([]byte(from + to + subject + mime + "\n" + buf.String()), map[string]interface{}{
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
