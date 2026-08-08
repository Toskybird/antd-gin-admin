package mock

import (
	"context"
	"strings"
	"sync"
	"time"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type DetectionRuleMockRepository struct {
	mu    sync.Mutex
	byCode map[string]*entity.DetectionRule
	seq   int64
}

func NewDetectionRuleMockRepository() *DetectionRuleMockRepository {
	return &DetectionRuleMockRepository{byCode: map[string]*entity.DetectionRule{}}
}

func (r *DetectionRuleMockRepository) UpsertByRuleCode(_ context.Context, rule *entity.DetectionRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	existing, ok := r.byCode[rule.RuleCode]
	if ok {
		existing.DisplayName = rule.DisplayName
		existing.UpdatedAt = now
		return nil
	}
	r.seq++
	cp := *rule
	cp.ID = r.seq
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.byCode[cp.RuleCode] = &cp
	return nil
}

func (r *DetectionRuleMockRepository) GetByCode(_ context.Context, ruleCode string) (*entity.DetectionRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[ruleCode]
	if !ok {
		return nil, nil
	}
	cp := *existing
	return &cp, nil
}

func (r *DetectionRuleMockRepository) UpdateEnabled(_ context.Context, ruleCode string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[ruleCode]
	if !ok {
		return nil
	}
	existing.Enabled = enabled
	existing.UpdatedAt = time.Now()
	return nil
}

func (r *DetectionRuleMockRepository) List(_ context.Context, page, pageSize int, filter interfaces.DetectionRuleListFilter) ([]*entity.DetectionRule, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := make([]*entity.DetectionRule, 0)
	keyword := strings.TrimSpace(strings.ToLower(filter.Keyword))
	for _, rule := range r.byCode {
		if filter.Enabled != nil && rule.Enabled != *filter.Enabled {
			continue
		}
		if keyword != "" {
			code := strings.ToLower(rule.RuleCode)
			name := strings.ToLower(rule.DisplayName)
			if !strings.Contains(code, keyword) && !strings.Contains(name, keyword) {
				continue
			}
		}
		cp := *rule
		items = append(items, &cp)
	}
	total := int64(len(items))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []*entity.DetectionRule{}, total, nil
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}

// StaticDetectionRuleSource is a fixture template source for tests/sync.
type StaticDetectionRuleSource struct {
	Rules []entity.DetectionRule
}

func (s *StaticDetectionRuleSource) ListTemplateRules(_ context.Context) ([]entity.DetectionRule, error) {
	out := make([]entity.DetectionRule, len(s.Rules))
	copy(out, s.Rules)
	return out, nil
}
