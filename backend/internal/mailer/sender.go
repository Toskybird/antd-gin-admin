package mailer

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// Message is one 邮件告警 delivery to all 收件人.
type Message struct {
	FromName    string
	FromAddress string
	To          []string
	Subject     string
	Body        string
}

// Sender is the 邮件发送端口 (SMTP in production, fake in tests).
type Sender interface {
	Send(ctx context.Context, cfg *entity.EmailAlertConfig, msg Message) error
}
