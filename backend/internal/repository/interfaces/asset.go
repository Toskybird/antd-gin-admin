package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
)

// AssetListFilter filters asset queries.
type AssetListFilter struct {
	Keyword   string
	Status    string
	DeptCodes []string
	All       bool
}

// AssetRepository persists Asset records.
type AssetRepository interface {
	Create(ctx context.Context, asset *entity.Asset) error
	Update(ctx context.Context, asset *entity.Asset) error
	DeleteByCode(ctx context.Context, assetCode string) error
	GetByCode(ctx context.Context, assetCode string) (*entity.Asset, error)
	GetByDeptAndRootURL(ctx context.Context, deptCode, rootURL string) (*entity.Asset, error)
	List(ctx context.Context, page, pageSize int, filter AssetListFilter) ([]*entity.Asset, int64, error)
}
