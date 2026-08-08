package finding

import (
	"context"
	"strings"

	"antd-gin-admin-backend/internal/model/entity"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

type Service struct {
	findings repointerfaces.FindingRepository
	jobs     repointerfaces.ScanJobRepository
	rules    repointerfaces.DetectionRuleRepository
}

func New(
	findings repointerfaces.FindingRepository,
	jobs repointerfaces.ScanJobRepository,
	rules repointerfaces.DetectionRuleRepository,
) serviceinterfaces.FindingService {
	return &Service{findings: findings, jobs: jobs, rules: rules}
}

func (s *Service) ListWithScope(ctx context.Context, page, pageSize int, jobCode, ruleCode, severity, keyword string, scope *serviceinterfaces.UserDataScope) ([]*entity.Finding, int64, error) {
	filter := repointerfaces.FindingListFilter{
		JobCode:  strings.TrimSpace(jobCode),
		RuleCode: strings.TrimSpace(ruleCode),
		Severity: strings.TrimSpace(severity),
		Keyword:  keyword,
		All:      scope == nil || scope.All,
	}
	if !filter.All && scope != nil {
		codes, err := s.visibleJobCodes(ctx, scope)
		if err != nil {
			return nil, 0, err
		}
		filter.JobCodes = codes
	}
	return s.findings.List(ctx, page, pageSize, filter)
}

func (s *Service) GetByCodeWithScope(ctx context.Context, findingCode string, scope *serviceinterfaces.UserDataScope) (*entity.Finding, error) {
	f, err := s.findings.GetByCode(ctx, findingCode)
	if err != nil || f == nil {
		return f, err
	}
	if scope == nil || scope.All {
		return f, nil
	}
	job, err := s.jobs.GetByCode(ctx, f.JobCode)
	if err != nil {
		return nil, err
	}
	if job == nil || !jobVisible(job, scope) {
		return nil, nil
	}
	return f, nil
}

func (s *Service) RuleDisplayName(ctx context.Context, ruleCode string) (string, error) {
	if s.rules == nil || strings.TrimSpace(ruleCode) == "" {
		return "", nil
	}
	rule, err := s.rules.GetByCode(ctx, ruleCode)
	if err != nil || rule == nil {
		return "", err
	}
	return rule.DisplayName, nil
}

func (s *Service) visibleJobCodes(ctx context.Context, scope *serviceinterfaces.UserDataScope) ([]string, error) {
	jobs, _, err := s.jobs.List(ctx, 1, 10000, repointerfaces.ScanJobListFilter{
		All:       false,
		DeptCodes: append([]string{}, scope.DeptCodes...),
		CreatedBy: scope.UserCode,
	})
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(jobs))
	for _, job := range jobs {
		codes = append(codes, job.JobCode)
	}
	return codes, nil
}

func jobVisible(job *entity.ScanJob, scope *serviceinterfaces.UserDataScope) bool {
	if scope == nil || scope.All {
		return true
	}
	for _, d := range scope.DeptCodes {
		if d == job.DeptCode {
			return true
		}
	}
	if scope.UserCode != "" && job.CreatedBy == scope.UserCode {
		return true
	}
	return false
}
