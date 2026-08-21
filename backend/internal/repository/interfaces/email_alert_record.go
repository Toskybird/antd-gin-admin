package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// EmailAlertRecordRepository persists 告警记录.
type EmailAlertRecordRepository interface {
	Create(ctx context.Context, rec *entity.EmailAlertRecord) error
	GetByOperationLogID(ctx context.Context, operationLogID int64) (*entity.EmailAlertRecord, error)
	List(ctx context.Context, filter *EmailAlertRecordFilter) ([]*entity.EmailAlertRecord, int64, error)
}

// EmailAlertRecordFilter defines list filters for 告警记录.
type EmailAlertRecordFilter struct {
	SendStatus string
	Module     string
	StartTime  *int64
	EndTime    *int64
	Page       int
	PageSize   int
}
