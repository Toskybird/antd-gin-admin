package db

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type DetectionRuleRepository struct {
	db *gorm.DB
}

func NewDetectionRuleRepository(db *gorm.DB) interfaces.DetectionRuleRepository {
	return &DetectionRuleRepository{db: db}
}

func (r *DetectionRuleRepository) UpsertByRuleCode(ctx context.Context, rule *entity.DetectionRule) error {
	existing, err := r.GetByCode(ctx, rule.RuleCode)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.db.WithContext(ctx).Create(rule).Error
	}
	return r.db.WithContext(ctx).Model(&entity.DetectionRule{}).
		Where("rule_code = ?", rule.RuleCode).
		Updates(map[string]interface{}{
			"display_name": rule.DisplayName,
			"updated_at":   time.Now(),
		}).Error
}

func (r *DetectionRuleRepository) GetByCode(ctx context.Context, ruleCode string) (*entity.DetectionRule, error) {
	var rule entity.DetectionRule
	if err := r.db.WithContext(ctx).Where("rule_code = ?", ruleCode).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *DetectionRuleRepository) UpdateEnabled(ctx context.Context, ruleCode string, enabled bool) error {
	return r.db.WithContext(ctx).Model(&entity.DetectionRule{}).
		Where("rule_code = ?", ruleCode).
		Updates(map[string]interface{}{
			"enabled":    enabled,
			"updated_at": time.Now(),
		}).Error
}

func (r *DetectionRuleRepository) List(ctx context.Context, page, pageSize int, filter interfaces.DetectionRuleListFilter) ([]*entity.DetectionRule, int64, error) {
	var rules []*entity.DetectionRule
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.DetectionRule{})
	if filter.Enabled != nil {
		query = query.Where("enabled = ?", *filter.Enabled)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(rule_code) LIKE ? OR LOWER(display_name) LIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.DetectionRule{}, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Limit(pageSize).Offset(offset).Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}
