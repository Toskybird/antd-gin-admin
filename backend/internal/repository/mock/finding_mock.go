package mock

import (
	"context"
	"sync"
	"time"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type FindingMockRepository struct {
	mu       sync.Mutex
	byCode   map[string]*entity.Finding
	byJob    map[string][]string
	seq      int64
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
