package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"antd-gin-admin-backend/internal/model/entity"
)

// SMTPSender delivers mail using the 邮件告警约定 SMTP settings.
type SMTPSender struct{}

func NewSMTPSender() Sender {
	return &SMTPSender{}
}

func (s *SMTPSender) Send(ctx context.Context, cfg *entity.EmailAlertConfig, msg Message) error {
	if cfg == nil {
		return fmt.Errorf("邮件告警约定缺失")
	}
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	from := msg.FromAddress
	if from == "" {
		from = cfg.SenderAddress
	}
	body := composeRFC822(msg, from)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	enc := strings.ToLower(strings.TrimSpace(cfg.Encryption))
	switch enc {
	case "ssl":
		return sendSSL(ctx, addr, cfg.Host, auth, from, msg.To, body)
	default:
		// none and starttls both use smtp.SendMail; starttls is handled by the server EHLO.
		return smtp.SendMail(addr, auth, from, msg.To, []byte(body))
	}
}

func sendSSL(ctx context.Context, addr, host string, auth smtp.Auth, from string, to []string, body string) error {
	dialer := &tls.Dialer{Config: &tls.Config{ServerName: host}}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := client.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(body)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func composeRFC822(msg Message, from string) string {
	name := msg.FromName
	headerFrom := from
	if name != "" {
		headerFrom = fmt.Sprintf("%s <%s>", name, from)
	}
	return fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		headerFrom, strings.Join(msg.To, ", "), msg.Subject, msg.Body)
}
