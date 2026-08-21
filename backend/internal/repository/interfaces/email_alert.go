package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// EmailAlertRepository persists the singleton 邮件告警约定.
type EmailAlertRepository interface {
	Get(ctx context.Context) (*entity.EmailAlertConfig, error)
	Save(ctx context.Context, cfg *entity.EmailAlertConfig) error
}
