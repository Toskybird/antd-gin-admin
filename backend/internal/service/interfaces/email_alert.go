package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
)

// EmailAlertService is the 邮件告警 domain entry for the singleton 约定.
type EmailAlertService interface {
	Get(ctx context.Context) (*entity.EmailAlertConfig, error)
	Save(ctx context.Context, req *dto.SaveEmailAlertRequest) (*entity.EmailAlertConfig, error)
}
