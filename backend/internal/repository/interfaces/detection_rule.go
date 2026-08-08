package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// DetectionRuleListFilter filters detection rule queries.
type DetectionRuleListFilter struct {
	Keyword string
	Enabled *bool
}

// DetectionRuleSource provides template rules for sync upsert.
type DetectionRuleSource interface {
	ListTemplateRules(ctx context.Context) ([]entity.DetectionRule, error)
}

// DetectionRuleRepository persists DetectionRule records.
type DetectionRuleRepository interface {
	UpsertByRuleCode(ctx context.Context, rule *entity.DetectionRule) error
	GetByCode(ctx context.Context, ruleCode string) (*entity.DetectionRule, error)
	UpdateEnabled(ctx context.Context, ruleCode string, enabled bool) error
	List(ctx context.Context, page, pageSize int, filter DetectionRuleListFilter) ([]*entity.DetectionRule, int64, error)
}
