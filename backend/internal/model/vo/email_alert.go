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

// EmailAlertRecordVO is a public 告警记录 (no pending send state).
type EmailAlertRecordVO struct {
	ID                int64  `json:"id"`
	OperationLogID    int64  `json:"operation_log_id"`
	Username          string `json:"username"`
	Module            string `json:"module"`
	Action            string `json:"action"`
	Path              string `json:"path"`
	ErrorMsg          string `json:"error_msg"`
	OccurredAt        int64  `json:"occurred_at"`
	SendStatus        string `json:"send_status"`
	SendFailureReason string `json:"send_failure_reason"`
}

func BuildEmailAlertRecordVO(rec *entity.EmailAlertRecord) *EmailAlertRecordVO {
	if rec == nil {
		return nil
	}
	return &EmailAlertRecordVO{
		ID:                rec.ID,
		OperationLogID:    rec.OperationLogID,
		Username:          rec.Username,
		Module:            rec.Module,
		Action:            rec.Action,
		Path:              rec.Path,
		ErrorMsg:          rec.ErrorMsg,
		OccurredAt:        rec.OccurredAt.Unix(),
		SendStatus:        rec.SendStatus,
		SendFailureReason: rec.SendFailureReason,
	}
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
