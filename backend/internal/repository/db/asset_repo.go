package db

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) interfaces.AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(ctx context.Context, asset *entity.Asset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *AssetRepository) Update(ctx context.Context, asset *entity.Asset) error {
	return r.db.WithContext(ctx).
		Model(&entity.Asset{}).
		Where("asset_code = ?", asset.AssetCode).
		Select("name", "root_url", "dept_code", "remark", "status", "updated_at").
		Updates(asset).Error
}

func (r *AssetRepository) DeleteByCode(ctx context.Context, assetCode string) error {
	return r.db.WithContext(ctx).Where("asset_code = ?", assetCode).Delete(&entity.Asset{}).Error
}

func (r *AssetRepository) GetByCode(ctx context.Context, assetCode string) (*entity.Asset, error) {
	var asset entity.Asset
	if err := r.db.WithContext(ctx).Where("asset_code = ?", assetCode).First(&asset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) GetByDeptAndRootURL(ctx context.Context, deptCode, rootURL string) (*entity.Asset, error) {
	var asset entity.Asset
	if err := r.db.WithContext(ctx).
		Where("dept_code = ? AND root_url = ?", deptCode, rootURL).
		First(&asset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) List(ctx context.Context, page, pageSize int, filter interfaces.AssetListFilter) ([]*entity.Asset, int64, error) {
	var assets []*entity.Asset
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.Asset{})
	if !filter.All {
		if len(filter.DeptCodes) == 0 {
			return []*entity.Asset{}, 0, nil
		}
		query = query.Where("dept_code IN ?", filter.DeptCodes)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(root_url) LIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.Asset{}, 0, nil
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}
