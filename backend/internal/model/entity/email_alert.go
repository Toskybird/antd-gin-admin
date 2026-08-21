package entity

import "time"

// EmailAlertConfig is the singleton 邮件告警约定.
type EmailAlertConfig struct {
	ID             int64     `json:"id" gorm:"primaryKey;comment:主键，固定为 1"`
	Enabled        bool      `json:"enabled" gorm:"comment:是否启用"`
	Host           string    `json:"host" gorm:"size:255;comment:邮件服务器主机"`
	Port           int       `json:"port" gorm:"comment:邮件服务器端口"`
	Encryption     string    `json:"encryption" gorm:"size:32;comment:加密方式：none、starttls、ssl"`
	Username       string    `json:"username" gorm:"size:128;comment:认证用户名"`
	Password       string    `json:"-" gorm:"size:255;comment:SMTP 密码，明文存储"`
	SenderName     string    `json:"sender_name" gorm:"size:128;comment:发件人显示名"`
	SenderAddress  string    `json:"sender_address" gorm:"size:128;comment:发件人地址"`
	RecipientsJSON string    `json:"-" gorm:"type:text;comment:收件人邮箱 JSON 数组"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (EmailAlertConfig) TableName() string { return "sys_email_alert_config" }
