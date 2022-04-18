package mailer

import (
	"net/smtp"
)

type Mailer struct {
	auth      smtp.Auth
	sender    string
	host      string
	port      string
	isDiscard bool
}

func New(host, port, username, password, sender string, isDiscard bool) Mailer {
	auth := smtp.PlainAuth("", sender, password, host)

	return Mailer{
		auth:      auth,
		sender:    sender,
		host:      host,
		port:      port,
		isDiscard: isDiscard,
	}
}

func (m Mailer) Send(recipient []string, data []byte) error {
	if !m.isDiscard {
		return smtp.SendMail(m.host+":"+m.port, m.auth, m.sender, recipient, data)
	} else {
		return nil
	}
}
