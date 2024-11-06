package notifications

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
)

type MailClient struct {
	auth   smtp.Auth
	addr   string
	from   string
	logger *slog.Logger
}

type MailInfo struct {
	ToList  []string
	Subject string
	Body    string
}

func NewMailClient(host string, port string, from string, password string, logger *slog.Logger) *MailClient {
	auth := smtp.PlainAuth("", from, password, host)
	return &MailClient{
		logger: logger,
		auth:   auth,
		addr:   fmt.Sprintf("%v:%v", host, port),
		from:   from,
	}
}

func (d *MailInfo) toMsg() []byte {
	return []byte(
		fmt.Sprintf("To: %v\r\n", strings.Join(d.ToList, " ")) +
			fmt.Sprintf("Subject: %v\r\n", d.Subject) +
			"\r\n" +
			fmt.Sprintf("%v\r\n", d.Body))
}

func (c *MailClient) SendMsg(toList []string, msg *MailInfo) error {
	err := smtp.SendMail(c.addr, c.auth, c.from, toList, []byte(msg.toMsg()))
	if err != nil {
		return err
	}

	return nil
}
