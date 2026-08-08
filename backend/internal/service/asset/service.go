package asset

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

type Service struct {
	repo repointerfaces.AssetRepository
}

func New(repo repointerfaces.AssetRepository) serviceinterfaces.AssetService {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *dto.CreateAssetRequest) (*entity.Asset, error) {
	rootURL, err := NormalizeRootURL(req.RootURL)
	if err != nil {
		return nil, apperrors.New(http.StatusBadRequest, "invalid_root_url", err.Error())
	}
	deptCode := strings.TrimSpace(req.DeptCode)
	if existing, err := s.repo.GetByDeptAndRootURL(ctx, deptCode, rootURL); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, apperrors.NewConflict("duplicate_asset_root_url", "同部门下该根 URL 已存在，请更换后重试。")
	}
	code, err := s.generateAssetCode(ctx)
	if err != nil {
		return nil, err
	}
	asset := &entity.Asset{
		AssetCode: code,
		Name:      strings.TrimSpace(req.Name),
		RootURL:   rootURL,
		DeptCode:  deptCode,
		Remark:    strings.TrimSpace(req.Remark),
		Status:    StatusActive,
	}
	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}
	return s.repo.GetByCode(ctx, code)
}

func (s *Service) Update(ctx context.Context, assetCode string, req *dto.UpdateAssetRequest) (*entity.Asset, error) {
	asset, err := s.repo.GetByCode(ctx, assetCode)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, nil
	}
	if req.Name != nil {
		asset.Name = strings.TrimSpace(*req.Name)
	}
	if req.Remark != nil {
		asset.Remark = strings.TrimSpace(*req.Remark)
	}
	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)
		if status != StatusActive && status != StatusDisabled {
			return nil, apperrors.New(http.StatusBadRequest, "invalid_status", "状态仅支持 active 或 disabled")
		}
		asset.Status = status
	}
	deptCode := asset.DeptCode
	if req.DeptCode != nil {
		deptCode = strings.TrimSpace(*req.DeptCode)
		asset.DeptCode = deptCode
	}
	rootURL := asset.RootURL
	if req.RootURL != nil {
		normalized, err := NormalizeRootURL(*req.RootURL)
		if err != nil {
			return nil, apperrors.New(http.StatusBadRequest, "invalid_root_url", err.Error())
		}
		rootURL = normalized
		asset.RootURL = rootURL
	}
	if existing, err := s.repo.GetByDeptAndRootURL(ctx, deptCode, rootURL); err != nil {
		return nil, err
	} else if existing != nil && existing.AssetCode != asset.AssetCode {
		return nil, apperrors.NewConflict("duplicate_asset_root_url", "同部门下该根 URL 已存在，请更换后重试。")
	}
	if err := s.repo.Update(ctx, asset); err != nil {
		return nil, err
	}
	return s.repo.GetByCode(ctx, assetCode)
}

func (s *Service) Delete(ctx context.Context, assetCode string) error {
	return s.repo.DeleteByCode(ctx, assetCode)
}

func (s *Service) GetByCode(ctx context.Context, assetCode string) (*entity.Asset, error) {
	return s.repo.GetByCode(ctx, assetCode)
}

func (s *Service) ListWithScope(ctx context.Context, page, pageSize int, keyword, status string, scope *serviceinterfaces.UserDataScope) ([]*entity.Asset, int64, error) {
	filter := repointerfaces.AssetListFilter{
		Keyword: keyword,
		Status:  strings.TrimSpace(status),
		All:     scope == nil || scope.All,
	}
	if !filter.All && scope != nil {
		filter.DeptCodes = append([]string{}, scope.DeptCodes...)
	}
	return s.repo.List(ctx, page, pageSize, filter)
}

func (s *Service) generateAssetCode(ctx context.Context) (string, error) {
	for i := 0; i < 8; i++ {
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			return "", err
		}
		code := fmt.Sprintf("AST-%s-%s", time.Now().Format("150405"), hex.EncodeToString(b[:]))
		existing, err := s.repo.GetByCode(ctx, code)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return code, nil
		}
	}
	return "", apperrors.New(http.StatusInternalServerError, "asset_code_generation_failed", "生成资产编码失败，请重试")
}

// NormalizeRootURL normalizes a Web entry URL to scheme+host+port+optional base path.
func NormalizeRootURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("根 URL 不能为空")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("根 URL 格式无效")
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("根 URL 仅支持 http/https")
	}
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
	}
	hostPort := host
	if port != "" {
		hostPort = host + ":" + port
	}
	path := u.EscapedPath()
	if path == "/" {
		path = ""
	} else {
		path = strings.TrimRight(path, "/")
	}
	return scheme + "://" + hostPort + path, nil
}
