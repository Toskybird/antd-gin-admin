package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	repo "antd-gin-admin-backend/internal/repository/interfaces"
)

type EmailAlertRecordRepository struct {
	db *gorm.DB
}

func NewEmailAlertRecordRepository(db *gorm.DB) repo.EmailAlertRecordRepository {
	return &EmailAlertRecordRepository{db: db}
}

func (r *EmailAlertRecordRepository) Create(ctx context.Context, rec *entity.EmailAlertRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *EmailAlertRecordRepository) GetByOperationLogID(ctx context.Context, operationLogID int64) (*entity.EmailAlertRecord, error) {
	var rec entity.EmailAlertRecord
	err := r.db.WithContext(ctx).Where("operation_log_id = ?", operationLogID).Limit(1).Find(&rec).Error
	if err != nil {
		return nil, err
	}
	if rec.ID == 0 {
		return nil, nil
	}
	return &rec, nil
}

func (r *EmailAlertRecordRepository) List(ctx context.Context, filter *repo.EmailAlertRecordFilter) ([]*entity.EmailAlertRecord, int64, error) {
	var items []*entity.EmailAlertRecord
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.EmailAlertRecord{})
	if filter != nil {
		if filter.SendStatus != "" {
			query = query.Where("send_status = ?", filter.SendStatus)
		}
		if filter.Module != "" {
			query = query.Where("module = ?", filter.Module)
		}
		if filter.StartTime != nil {
			query = query.Where("occurred_at >= ?", time.Unix(*filter.StartTime, 0))
		}
		if filter.EndTime != nil {
			query = query.Where("occurred_at <= ?", time.Unix(*filter.EndTime, 0))
		}
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, size := 1, 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.PageSize > 0 {
			size = filter.PageSize
		}
	}
	err := query.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = []*entity.EmailAlertRecord{}
	}
	return items, total, nil
}
