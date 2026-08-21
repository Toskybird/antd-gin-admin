package vo

import (
	"encoding/json"

	"antd-gin-admin-backend/internal/model/entity"
)

// EmailAlertVO is the public 邮件告警约定 (no SMTP password).
type EmailAlertVO struct {
	Enabled            bool     `json:"enabled"`
	Host               string   `json:"host"`
	Port               int      `json:"port"`
	Encryption         string   `json:"encryption"`
	Username           string   `json:"username"`
	PasswordConfigured bool     `json:"password_configured"`
	SenderName         string   `json:"sender_name"`
	SenderAddress      string   `json:"sender_address"`
	Recipients         []string `json:"recipients"`
}

func BuildEmailAlertVO(cfg *entity.EmailAlertConfig) *EmailAlertVO {
	if cfg == nil {
		return &EmailAlertVO{Recipients: []string{}}
	}
	recipients := []string{}
	if cfg.RecipientsJSON != "" {
		_ = json.Unmarshal([]byte(cfg.RecipientsJSON), &recipients)
		if recipients == nil {
			recipients = []string{}
		}
	}
	return &EmailAlertVO{
		Enabled:            cfg.Enabled,
		Host:               cfg.Host,
		Port:               cfg.Port,
		Encryption:         cfg.Encryption,
		Username:           cfg.Username,
		PasswordConfigured: cfg.Password != "",
		SenderName:         cfg.SenderName,
		SenderAddress:      cfg.SenderAddress,
		Recipients:         recipients,
	}
}
