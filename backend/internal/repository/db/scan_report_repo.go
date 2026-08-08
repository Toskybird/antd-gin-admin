package db

import (
	"context"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type ScanReportRepository struct {
	db *gorm.DB
}

func NewScanReportRepository(db *gorm.DB) interfaces.ScanReportRepository {
	return &ScanReportRepository{db: db}
}

func (r *ScanReportRepository) Create(ctx context.Context, report *entity.ScanReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *ScanReportRepository) GetByCode(ctx context.Context, reportCode string) (*entity.ScanReport, error) {
	var item entity.ScanReport
	err := r.db.WithContext(ctx).Where("report_code = ?", reportCode).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ScanReportRepository) NextVersion(ctx context.Context, jobCode string) (int, error) {
	var max *int
	err := r.db.WithContext(ctx).Model(&entity.ScanReport{}).
		Where("job_code = ?", jobCode).
		Select("MAX(version)").
		Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if max == nil {
		return 1, nil
	}
	return *max + 1, nil
}

func (r *ScanReportRepository) List(ctx context.Context, page, pageSize int, filter interfaces.ScanReportListFilter) ([]*entity.ScanReport, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.ScanReport{})
	if !filter.All {
		if len(filter.JobCodes) == 0 {
			return []*entity.ScanReport{}, 0, nil
		}
		q = q.Where("job_code IN ?", filter.JobCodes)
	}
	if filter.JobCode != "" {
		q = q.Where("job_code = ?", filter.JobCode)
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
	var items []*entity.ScanReport
	err := q.Order("version DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}
