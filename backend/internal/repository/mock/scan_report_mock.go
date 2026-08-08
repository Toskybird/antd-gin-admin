package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type MemoryObjectStorage struct {
	mu   sync.Mutex
	data map[string][]byte
}

func NewMemoryObjectStorage() *MemoryObjectStorage {
	return &MemoryObjectStorage{data: map[string][]byte{}}
}

func (s *MemoryObjectStorage) Put(_ context.Context, key string, content []byte, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]byte, len(content))
	copy(cp, content)
	s.data[key] = cp
	return nil
}

func (s *MemoryObjectStorage) Get(_ context.Context, key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	content, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("object not found: %s", key)
	}
	cp := make([]byte, len(content))
	copy(cp, content)
	return cp, nil
}

type ScanReportMockRepository struct {
	mu     sync.Mutex
	byCode map[string]*entity.ScanReport
	seq    int64
}

func NewScanReportMockRepository() *ScanReportMockRepository {
	return &ScanReportMockRepository{byCode: map[string]*entity.ScanReport{}}
}

func (r *ScanReportMockRepository) Create(_ context.Context, report *entity.ScanReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.seq++
	cp := *report
	cp.ID = r.seq
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.byCode[cp.ReportCode] = &cp
	return nil
}

func (r *ScanReportMockRepository) GetByCode(_ context.Context, reportCode string) (*entity.ScanReport, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[reportCode]
	if !ok {
		return nil, nil
	}
	cp := *existing
	return &cp, nil
}

func (r *ScanReportMockRepository) NextVersion(_ context.Context, jobCode string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	max := 0
	for _, report := range r.byCode {
		if report.JobCode == jobCode && report.Version > max {
			max = report.Version
		}
	}
	return max + 1, nil
}

func (r *ScanReportMockRepository) List(_ context.Context, page, pageSize int, filter interfaces.ScanReportListFilter) ([]*entity.ScanReport, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	allowed := map[string]struct{}{}
	if !filter.All {
		for _, code := range filter.JobCodes {
			allowed[code] = struct{}{}
		}
	}
	items := make([]*entity.ScanReport, 0)
	for _, report := range r.byCode {
		if !filter.All {
			if _, ok := allowed[report.JobCode]; !ok {
				continue
			}
		}
		if filter.JobCode != "" && report.JobCode != filter.JobCode {
			continue
		}
		cp := *report
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
		return []*entity.ScanReport{}, total, nil
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}
