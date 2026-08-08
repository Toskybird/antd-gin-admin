package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// ScanJobListFilter filters scan job queries.
type ScanJobListFilter struct {
	Status    string
	Policy    string
	Keyword   string
	DeptCodes []string
	All       bool
	CreatedBy string // when self-scope without depts
}

// ScanJobQueue enqueues jobs for workers.
type ScanJobQueue interface {
	Enqueue(ctx context.Context, jobCode string) error
	ListEnqueued(ctx context.Context) ([]string, error)
	Dequeue(ctx context.Context) (string, bool, error)
}

// ScanJobRepository persists ScanJob records.
type ScanJobRepository interface {
	Create(ctx context.Context, job *entity.ScanJob) error
	Update(ctx context.Context, job *entity.ScanJob) error
	GetByCode(ctx context.Context, jobCode string) (*entity.ScanJob, error)
	List(ctx context.Context, page, pageSize int, filter ScanJobListFilter) ([]*entity.ScanJob, int64, error)
}
