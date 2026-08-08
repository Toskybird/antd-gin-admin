package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// FindingService provides read-only Finding queries.
type FindingService interface {
	ListWithScope(ctx context.Context, page, pageSize int, jobCode, ruleCode, severity, keyword string, scope *UserDataScope) ([]*entity.Finding, int64, error)
	GetByCodeWithScope(ctx context.Context, findingCode string, scope *UserDataScope) (*entity.Finding, error)
	RuleDisplayName(ctx context.Context, ruleCode string) (string, error)
}
