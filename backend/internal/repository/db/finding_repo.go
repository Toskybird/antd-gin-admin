package db

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type FindingRepository struct {
	db *gorm.DB
}

func NewFindingRepository(db *gorm.DB) interfaces.FindingRepository {
	return &FindingRepository{db: db}
}

func (r *FindingRepository) Create(ctx context.Context, finding *entity.Finding) error {
	return r.db.WithContext(ctx).Create(finding).Error
}

func (r *FindingRepository) CountByJob(ctx context.Context, jobCode string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.Finding{}).Where("job_code = ?", jobCode).Count(&total).Error
	return total, err
}

func (r *FindingRepository) ListByJob(ctx context.Context, jobCode string) ([]*entity.Finding, error) {
	var items []*entity.Finding
	err := r.db.WithContext(ctx).Where("job_code = ?", jobCode).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *FindingRepository) GetByCode(ctx context.Context, findingCode string) (*entity.Finding, error) {
	var item entity.Finding
	err := r.db.WithContext(ctx).Where("finding_code = ?", findingCode).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *FindingRepository) List(ctx context.Context, page, pageSize int, filter interfaces.FindingListFilter) ([]*entity.Finding, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Finding{})
	if !filter.All {
		if len(filter.JobCodes) == 0 {
			return []*entity.Finding{}, 0, nil
		}
		q = q.Where("job_code IN ?", filter.JobCodes)
	}
	if filter.JobCode != "" {
		q = q.Where("job_code = ?", filter.JobCode)
	}
	if filter.RuleCode != "" {
		q = q.Where("rule_code = ?", filter.RuleCode)
	}
	if filter.Severity != "" {
		q = q.Where("severity = ?", filter.Severity)
	}
	if kw := strings.TrimSpace(filter.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ? OR location ILIKE ? OR evidence ILIKE ? OR finding_code ILIKE ?", like, like, like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	var items []*entity.Finding
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}
