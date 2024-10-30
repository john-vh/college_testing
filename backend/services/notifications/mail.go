package notifications

import (
	"fmt"
	"net/smtp"
)

type MailClient struct {
	auth smtp.Auth
	addr string
	from string
}

func NewMailClient(host string, port string, from string, password string) *MailClient {
	auth := smtp.PlainAuth("", from, password, host)
	return &MailClient{
		auth: auth,
		addr: fmt.Sprintf("%v:%v", host, port),
		from: from,
	}
}

func (c *MailClient) SendMsg(toList []string, msg string) error {
	err := smtp.SendMail(c.addr, c.auth, c.from, toList, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
