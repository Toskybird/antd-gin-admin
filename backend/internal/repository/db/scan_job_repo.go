package db

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type ScanJobRepository struct {
	db *gorm.DB
}

func NewScanJobRepository(db *gorm.DB) interfaces.ScanJobRepository {
	return &ScanJobRepository{db: db}
}

func (r *ScanJobRepository) Create(ctx context.Context, job *entity.ScanJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *ScanJobRepository) Update(ctx context.Context, job *entity.ScanJob) error {
	return r.db.WithContext(ctx).
		Model(&entity.ScanJob{}).
		Where("job_code = ?", job.JobCode).
		Select("status", "finding_count", "finished_at", "updated_at").
		Updates(job).Error
}

func (r *ScanJobRepository) GetByCode(ctx context.Context, jobCode string) (*entity.ScanJob, error) {
	var job entity.ScanJob
	if err := r.db.WithContext(ctx).Where("job_code = ?", jobCode).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *ScanJobRepository) List(ctx context.Context, page, pageSize int, filter interfaces.ScanJobListFilter) ([]*entity.ScanJob, int64, error) {
	var jobs []*entity.ScanJob
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.ScanJob{})
	if !filter.All {
		conds := make([]string, 0, 2)
		args := make([]interface{}, 0, 2)
		if len(filter.DeptCodes) > 0 {
			conds = append(conds, "dept_code IN ?")
			args = append(args, filter.DeptCodes)
		}
		if filter.CreatedBy != "" {
			conds = append(conds, "created_by = ?")
			args = append(args, filter.CreatedBy)
		}
		if len(conds) == 0 {
			return []*entity.ScanJob{}, 0, nil
		}
		query = query.Where(strings.Join(conds, " OR "), args...)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	if policy := strings.TrimSpace(filter.Policy); policy != "" {
		query = query.Where("policy = ?", policy)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(entry_url) LIKE ? OR LOWER(asset_code) LIKE ? OR LOWER(job_code) LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.ScanJob{}, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
}
