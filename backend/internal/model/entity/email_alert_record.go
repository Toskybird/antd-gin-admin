package entity

import "time"

const (
	EmailAlertSendSuccess = "success"
	EmailAlertSendFailed  = "failed"
)

// EmailAlertRecord is the 告警记录 for one failed 操作日志.
type EmailAlertRecord struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	OperationLogID    int64     `json:"operation_log_id" gorm:"uniqueIndex;comment:操作日志 ID"`
	Username          string    `json:"username" gorm:"size:64;comment:快照用户名"`
	Module            string    `json:"module" gorm:"size:64;comment:快照模块"`
	Action            string    `json:"action" gorm:"size:64;comment:快照动作"`
	Path              string    `json:"path" gorm:"size:255;comment:快照路径"`
	ErrorMsg          string    `json:"error_msg" gorm:"size:512;comment:快照错误信息"`
	OccurredAt        time.Time `json:"occurred_at" gorm:"comment:操作发生时间"`
	SendStatus        string    `json:"send_status" gorm:"size:32;index;comment:发送成功或发送失败"`
	SendFailureReason string    `json:"send_failure_reason" gorm:"size:512;comment:最近一次发送失败原因"`
	CreatedAt         time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (EmailAlertRecord) TableName() string { return "sys_email_alert_records" }
