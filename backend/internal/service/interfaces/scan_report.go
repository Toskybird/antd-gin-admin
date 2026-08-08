package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// ScanReportService manages Scan Report generation and download.
type ScanReportService interface {
	Generate(ctx context.Context, jobCode, creatorUserCode string) (*entity.ScanReport, error)
	ListWithScope(ctx context.Context, page, pageSize int, jobCode string, scope *UserDataScope) ([]*entity.ScanReport, int64, error)
	Download(ctx context.Context, reportCode string, scope *UserDataScope) (content []byte, contentType string, filename string, err error)
}
