package mock

import (
	"context"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
)

type AssetMockRepository struct {
	mu     sync.Mutex
	byCode map[string]*entity.Asset
	seq    int64
}

func NewAssetMockRepository() *AssetMockRepository {
	return &AssetMockRepository{byCode: map[string]*entity.Asset{}}
}

func (r *AssetMockRepository) Create(_ context.Context, asset *entity.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.seq++
	cp := *asset
	cp.ID = r.seq
	cp.CreatedAt = now
	cp.UpdatedAt = now
	r.byCode[cp.AssetCode] = &cp
	return nil
}

func (r *AssetMockRepository) Update(_ context.Context, asset *entity.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[asset.AssetCode]
	if !ok || existing.DeletedAt.Valid {
		return gorm.ErrRecordNotFound
	}
	cp := *asset
	cp.ID = existing.ID
	cp.CreatedAt = existing.CreatedAt
	cp.UpdatedAt = time.Now()
	r.byCode[cp.AssetCode] = &cp
	return nil
}

func (r *AssetMockRepository) DeleteByCode(_ context.Context, assetCode string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[assetCode]
	if !ok || existing.DeletedAt.Valid {
		return nil
	}
	existing.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return nil
}

func (r *AssetMockRepository) GetByCode(_ context.Context, assetCode string) (*entity.Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byCode[assetCode]
	if !ok || existing.DeletedAt.Valid {
		return nil, nil
	}
	cp := *existing
	return &cp, nil
}

func (r *AssetMockRepository) GetByDeptAndRootURL(_ context.Context, deptCode, rootURL string) (*entity.Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, a := range r.byCode {
		if a.DeletedAt.Valid {
			continue
		}
		if a.DeptCode == deptCode && a.RootURL == rootURL {
			cp := *a
			return &cp, nil
		}
	}
	return nil, nil
}

func (r *AssetMockRepository) List(_ context.Context, page, pageSize int, filter interfaces.AssetListFilter) ([]*entity.Asset, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := make([]*entity.Asset, 0)
	keyword := strings.TrimSpace(strings.ToLower(filter.Keyword))
	for _, a := range r.byCode {
		if a.DeletedAt.Valid {
			continue
		}
		if !filter.All {
			if len(filter.DeptCodes) == 0 {
				continue
			}
			matched := false
			for _, d := range filter.DeptCodes {
				if d == a.DeptCode {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		if filter.Status != "" && a.Status != filter.Status {
			continue
		}
		if keyword != "" {
			name := strings.ToLower(a.Name)
			url := strings.ToLower(a.RootURL)
			if !strings.Contains(name, keyword) && !strings.Contains(url, keyword) {
				continue
			}
		}
		cp := *a
		items = append(items, &cp)
	}
	total := int64(len(items))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []*entity.Asset{}, total, nil
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}
