package detectionrule

import (
	"context"
	"net/http"
	"strings"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/model/entity"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo   repointerfaces.DetectionRuleRepository
	source repointerfaces.DetectionRuleSource
}

func New(repo repointerfaces.DetectionRuleRepository, source repointerfaces.DetectionRuleSource) serviceinterfaces.DetectionRuleService {
	return &Service{repo: repo, source: source}
}

func (s *Service) Sync(ctx context.Context) (int, error) {
	if s.source == nil {
		return 0, apperrors.New(http.StatusServiceUnavailable, "rule_source_unavailable", "检测规则同步源不可用")
	}
	templates, err := s.source.ListTemplateRules(ctx)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, tpl := range templates {
		code := strings.TrimSpace(tpl.RuleCode)
		name := strings.TrimSpace(tpl.DisplayName)
		if code == "" || name == "" {
			continue
		}
		existing, err := s.repo.GetByCode(ctx, code)
		if err != nil {
			return count, err
		}
		rule := &entity.DetectionRule{
			RuleCode:    code,
			DisplayName: name,
			Enabled:     true,
		}
		if existing != nil {
			// preserve enabled; only refresh display name via upsert
			rule.Enabled = existing.Enabled
		}
		if err := s.repo.UpsertByRuleCode(ctx, rule); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *Service) List(ctx context.Context, page, pageSize int, keyword string, enabled *bool) ([]*entity.DetectionRule, int64, error) {
	return s.repo.List(ctx, page, pageSize, repointerfaces.DetectionRuleListFilter{
		Keyword: keyword,
		Enabled: enabled,
	})
}

func (s *Service) SetEnabled(ctx context.Context, ruleCode string, enabled bool) (*entity.DetectionRule, error) {
	ruleCode = strings.TrimSpace(ruleCode)
	existing, err := s.repo.GetByCode(ctx, ruleCode)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	if err := s.repo.UpdateEnabled(ctx, ruleCode, enabled); err != nil {
		return nil, err
	}
	return s.repo.GetByCode(ctx, ruleCode)
}

func (s *Service) ListEnabledCodes(ctx context.Context) ([]string, error) {
	enabled := true
	items, _, err := s.repo.List(ctx, 1, 10000, repointerfaces.DetectionRuleListFilter{Enabled: &enabled})
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(items))
	for _, item := range items {
		codes = append(codes, item.RuleCode)
	}
	return codes, nil
}
