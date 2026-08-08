package mock

import (
	"context"
	"strings"
	"sync"
	"time"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type ScanJobMockRepository struct {
	mu     sync.Mutex
	byCode map[string]*entity.ScanJob
	seq    int64
}

func NewScanJobMockRepository() *ScanJobMockRepository {
	return &ScanJobMockRepository{byCode: map[string]*entity.ScanJob{}}
}

func (r *ScanJobMockRepository) Create(_ context.Context, job *entity.ScanJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.seq++
	cp := *job
	cp.ID = r.seq
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.byCode[cp.JobCode] = &cp
	return nil
}

func (r *ScanJobMockRepository) Update(_ context.Context, job *entity.ScanJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[job.JobCode]
	if !ok {
		return nil
	}
	cp := *job
	cp.ID = existing.ID
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.byCode[cp.JobCode] = &cp
	return nil
}

func (r *ScanJobMockRepository) GetByCode(_ context.Context, jobCode string) (*entity.ScanJob, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[jobCode]
	if !ok {
		return nil, nil
	}
	cp := *existing
	return &cp, nil
}

func (r *ScanJobMockRepository) List(_ context.Context, page, pageSize int, filter interfaces.ScanJobListFilter) ([]*entity.ScanJob, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := make([]*entity.ScanJob, 0)
	keyword := strings.TrimSpace(strings.ToLower(filter.Keyword))
	for _, job := range r.byCode {
		if !filter.All {
			matched := false
			for _, d := range filter.DeptCodes {
				if d == job.DeptCode {
					matched = true
					break
				}
			}
			if !matched && filter.CreatedBy != "" && job.CreatedBy == filter.CreatedBy {
				matched = true
			}
			if !matched {
				continue
			}
		}
		if filter.Status != "" && job.Status != filter.Status {
			continue
		}
		if filter.Policy != "" && job.Policy != filter.Policy {
			continue
		}
		if keyword != "" {
			url := strings.ToLower(job.EntryURL)
			code := strings.ToLower(job.AssetCode)
			if !strings.Contains(url, keyword) && !strings.Contains(code, keyword) && !strings.Contains(strings.ToLower(job.JobCode), keyword) {
				continue
			}
		}
		cp := *job
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
		return []*entity.ScanJob{}, total, nil
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}

type MemoryScanJobQueue struct {
	mu    sync.Mutex
	codes []string
}

func NewMemoryScanJobQueue() *MemoryScanJobQueue {
	return &MemoryScanJobQueue{codes: make([]string, 0)}
}

func (q *MemoryScanJobQueue) Enqueue(_ context.Context, jobCode string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.codes = append(q.codes, jobCode)
	return nil
}

func (q *MemoryScanJobQueue) ListEnqueued(_ context.Context) ([]string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]string, len(q.codes))
	copy(out, q.codes)
	return out, nil
}

func (q *MemoryScanJobQueue) Dequeue(_ context.Context) (string, bool, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.codes) == 0 {
		return "", false, nil
	}
	code := q.codes[len(q.codes)-1]
	q.codes = q.codes[:len(q.codes)-1]
	return code, true, nil
}
