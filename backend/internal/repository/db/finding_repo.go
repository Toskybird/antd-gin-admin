package db

import (
	"context"

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
