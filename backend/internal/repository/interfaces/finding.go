package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// EngineFinding is one finding emitted by a scan engine adapter.
type EngineFinding struct {
	RuleCode    string
	Severity    string
	Title       string
	Description string
	Evidence    string
	Location    string
}

// ScanRequest is input for a scan engine adapter.
type ScanRequest struct {
	JobCode      string
	EntryURL     string
	Policy       string
	MaxDepth     int
	MaxPages     int
	EnabledRules []string
}

// ScanEngineAdapter executes a scan against a target.
type ScanEngineAdapter interface {
	Scan(ctx context.Context, req ScanRequest) ([]EngineFinding, error)
}

// FindingListFilter filters Finding list queries.
type FindingListFilter struct {
	JobCode  string
	RuleCode string
	Severity string
	Keyword  string
	All      bool
	JobCodes []string // when !All, only findings whose JobCode is in this list
}

// FindingRepository persists Finding records.
type FindingRepository interface {
	Create(ctx context.Context, finding *entity.Finding) error
	CountByJob(ctx context.Context, jobCode string) (int64, error)
	ListByJob(ctx context.Context, jobCode string) ([]*entity.Finding, error)
	GetByCode(ctx context.Context, findingCode string) (*entity.Finding, error)
	List(ctx context.Context, page, pageSize int, filter FindingListFilter) ([]*entity.Finding, int64, error)
}
