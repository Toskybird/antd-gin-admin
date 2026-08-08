package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
)

// AssetService manages Asset ledger operations.
type AssetService interface {
	Create(ctx context.Context, req *dto.CreateAssetRequest) (*entity.Asset, error)
	Update(ctx context.Context, assetCode string, req *dto.UpdateAssetRequest) (*entity.Asset, error)
	Delete(ctx context.Context, assetCode string) error
	GetByCode(ctx context.Context, assetCode string) (*entity.Asset, error)
	ListWithScope(ctx context.Context, page, pageSize int, keyword, status string, scope *UserDataScope) ([]*entity.Asset, int64, error)
}
