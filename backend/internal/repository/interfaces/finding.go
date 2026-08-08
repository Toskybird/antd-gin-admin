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

// FindingRepository persists Finding records.
type FindingRepository interface {
	Create(ctx context.Context, finding *entity.Finding) error
	CountByJob(ctx context.Context, jobCode string) (int64, error)
	ListByJob(ctx context.Context, jobCode string) ([]*entity.Finding, error)
}
