package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
)

// ScanJobService manages Scan Job operations.
type ScanJobService interface {
	Create(ctx context.Context, req *dto.CreateScanJobRequest, creatorUserCode, creatorDeptCode string) (*entity.ScanJob, error)
	GetByCode(ctx context.Context, jobCode string) (*entity.ScanJob, error)
	ListWithScope(ctx context.Context, page, pageSize int, status, policy, keyword string, scope *UserDataScope) ([]*entity.ScanJob, int64, error)
	Cancel(ctx context.Context, jobCode string) (*entity.ScanJob, error)
}
