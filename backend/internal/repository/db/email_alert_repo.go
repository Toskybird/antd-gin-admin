package db

import (
	"context"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	repo "antd-gin-admin-backend/internal/repository/interfaces"
)

const emailAlertSingletonID int64 = 1

type EmailAlertRepository struct {
	db *gorm.DB
}

func NewEmailAlertRepository(db *gorm.DB) repo.EmailAlertRepository {
	return &EmailAlertRepository{db: db}
}

func (r *EmailAlertRepository) Get(ctx context.Context) (*entity.EmailAlertConfig, error) {
	var cfg entity.EmailAlertConfig
	err := r.db.WithContext(ctx).Where("id = ?", emailAlertSingletonID).Limit(1).Find(&cfg).Error
	if err != nil {
		return nil, err
	}
	if cfg.ID == 0 {
		return nil, nil
	}
	return &cfg, nil
}

func (r *EmailAlertRepository) Save(ctx context.Context, cfg *entity.EmailAlertConfig) error {
	cfg.ID = emailAlertSingletonID
	return r.db.WithContext(ctx).Save(cfg).Error
}
