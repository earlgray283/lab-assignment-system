package lib

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type SmtpClient struct {
	addr  string
	email string
	auth  smtp.Auth
}

func NewSmtpClient(email, password, host, port string) *SmtpClient {
	return &SmtpClient{
		addr:  fmt.Sprintf("%s:%s", host, port),
		email: email,
		auth:  smtp.PlainAuth("", email, password, host),
	}
}

func (c *SmtpClient) SendMail(to string, msg []byte) error {
	conn, err := smtp.Dial(c.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		ServerName:         c.addr,
		InsecureSkipVerify: true,
	}
	if err := conn.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err := conn.Auth(c.auth); err != nil {
		return err
	}
	if err := conn.Mail(c.email); err != nil {
		return err
	}
	if err := conn.Rcpt(to); err != nil {
		return err
	}
	w, err := conn.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	w.Close()

	return nil
}

func MakeMailBody(subject string, lines []string) []byte {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "Subject: %s\r\n\r\n", subject)
	for _, line := range lines {
		fmt.Fprintf(builder, "%s\r\n", line)
	}
	return []byte(builder.String())
}
