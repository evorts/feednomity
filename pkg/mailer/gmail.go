package mailer

type gmailManager struct {
	address  string
	identity string
	username string
	password string
	host     string
	sender   string
}

func NewGmail(address, identity, username, password, host, sender string) IMailer {
	return &gmailManager{
		address:  address,
		identity: identity,
		username: username,
		password: password,
		host:     host,
		sender:   sender,
	}
}

func (g *gmailManager) SendHtml() {
	panic("implement me")
}

