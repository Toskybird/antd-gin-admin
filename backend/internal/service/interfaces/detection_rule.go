package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// DetectionRuleService manages Detection Rule catalog operations.
type DetectionRuleService interface {
	Sync(ctx context.Context) (int, error)
	List(ctx context.Context, page, pageSize int, keyword string, enabled *bool) ([]*entity.DetectionRule, int64, error)
	SetEnabled(ctx context.Context, ruleCode string, enabled bool) (*entity.DetectionRule, error)
	ListEnabledCodes(ctx context.Context) ([]string, error)
}
