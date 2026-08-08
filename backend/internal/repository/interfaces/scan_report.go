package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// ObjectStorage stores report bodies (S3-compatible).
type ObjectStorage interface {
	Put(ctx context.Context, key string, content []byte, contentType string) error
	Get(ctx context.Context, key string) ([]byte, error)
}

// ScanReportListFilter filters Scan Report list queries.
type ScanReportListFilter struct {
	JobCode  string
	All      bool
	JobCodes []string
}

// ScanReportRepository persists Scan Report metadata.
type ScanReportRepository interface {
	Create(ctx context.Context, report *entity.ScanReport) error
	GetByCode(ctx context.Context, reportCode string) (*entity.ScanReport, error)
	NextVersion(ctx context.Context, jobCode string) (int, error)
	List(ctx context.Context, page, pageSize int, filter ScanReportListFilter) ([]*entity.ScanReport, int64, error)
}
