package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/model/vo"
)

// EmailAlertService is the 邮件告警 domain entry.
type EmailAlertService interface {
	Get(ctx context.Context) (*entity.EmailAlertConfig, error)
	Save(ctx context.Context, req *dto.SaveEmailAlertRequest) (*entity.EmailAlertConfig, error)
	NotifyFailedOperation(ctx context.Context, log *entity.OperationLog) error
	PageRecords(ctx context.Context, req *dto.EmailAlertRecordPageRequest) (*vo.PageResult[vo.EmailAlertRecordVO], error)
}
