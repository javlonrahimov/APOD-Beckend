package mailer

import (
	"net/smtp"
)

type Mailer struct {
	auth smtp.Auth
	sender string
	host string
	port string
}

func New(host, port, username, password, sender string) Mailer {
	auth := smtp.PlainAuth("", sender, password, host)

	return Mailer{
		auth: auth,
		sender: sender,
		host: host,
		port: port,
	}
}

func (m Mailer) Send(recipient []string, data []byte) error {
	return smtp.SendMail(m.host + ":" + m.port, m.auth, m.sender, recipient, data)
}
