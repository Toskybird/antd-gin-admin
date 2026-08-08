package mock

import (
	"context"
	"strings"
	"sync"
	"time"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type FindingMockRepository struct {
	mu     sync.Mutex
	byCode map[string]*entity.Finding
	byJob  map[string][]string
	seq    int64
}

func NewFindingMockRepository() *FindingMockRepository {
	return &FindingMockRepository{
		byCode: map[string]*entity.Finding{},
		byJob:  map[string][]string{},
	}
}

func (r *FindingMockRepository) Create(_ context.Context, finding *entity.Finding) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.seq++
	cp := *finding
	cp.ID = r.seq
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.byCode[cp.FindingCode] = &cp
	r.byJob[cp.JobCode] = append(r.byJob[cp.JobCode], cp.FindingCode)
	return nil
}

func (r *FindingMockRepository) CountByJob(_ context.Context, jobCode string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return int64(len(r.byJob[jobCode])), nil
}

func (r *FindingMockRepository) ListByJob(_ context.Context, jobCode string) ([]*entity.Finding, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	codes := r.byJob[jobCode]
	out := make([]*entity.Finding, 0, len(codes))
	for _, code := range codes {
		if f, ok := r.byCode[code]; ok {
			cp := *f
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *FindingMockRepository) GetByCode(_ context.Context, findingCode string) (*entity.Finding, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[findingCode]
	if !ok {
		return nil, nil
	}
	cp := *existing
	return &cp, nil
}

func (r *FindingMockRepository) List(_ context.Context, page, pageSize int, filter interfaces.FindingListFilter) ([]*entity.Finding, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	allowed := map[string]struct{}{}
	if !filter.All {
		for _, code := range filter.JobCodes {
			allowed[code] = struct{}{}
		}
	}
	keyword := strings.TrimSpace(strings.ToLower(filter.Keyword))
	items := make([]*entity.Finding, 0)
	for _, f := range r.byCode {
		if !filter.All {
			if _, ok := allowed[f.JobCode]; !ok {
				continue
			}
		}
		if filter.JobCode != "" && f.JobCode != filter.JobCode {
			continue
		}
		if filter.RuleCode != "" && f.RuleCode != filter.RuleCode {
			continue
		}
		if filter.Severity != "" && f.Severity != filter.Severity {
			continue
		}
		if keyword != "" {
			hay := strings.ToLower(f.Title + " " + f.Description + " " + f.Location + " " + f.Evidence + " " + f.FindingCode)
			if !strings.Contains(hay, keyword) {
				continue
			}
		}
		cp := *f
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
		return []*entity.Finding{}, total, nil
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}

// FakeScanEngine returns deterministic findings for enabled rules.
type FakeScanEngine struct {
	Findings []interfaces.EngineFinding
	Err      error
}

func (e *FakeScanEngine) Scan(_ context.Context, req interfaces.ScanRequest) ([]interfaces.EngineFinding, error) {
	if e.Err != nil {
		return nil, e.Err
	}
	enabled := map[string]struct{}{}
	for _, code := range req.EnabledRules {
		enabled[code] = struct{}{}
	}
	out := make([]interfaces.EngineFinding, 0)
	for _, f := range e.Findings {
		if _, ok := enabled[f.RuleCode]; ok {
			out = append(out, f)
		}
	}
	return out, nil
}
