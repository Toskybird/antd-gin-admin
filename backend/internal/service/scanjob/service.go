package scanjob

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	assetsvc "antd-gin-admin-backend/internal/service/asset"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"

	SourceAsset   = "asset"
	SourceTempURL = "temp_url"

	PolicyQuick    = "quick"
	PolicyStandard = "standard"
	PolicyDeep     = "deep"
)

type Service struct {
	repo      repointerfaces.ScanJobRepository
	assetRepo repointerfaces.AssetRepository
	queue     repointerfaces.ScanJobQueue
}

func New(repo repointerfaces.ScanJobRepository, assetRepo repointerfaces.AssetRepository, queue repointerfaces.ScanJobQueue) serviceinterfaces.ScanJobService {
	return &Service{repo: repo, assetRepo: assetRepo, queue: queue}
}

func (s *Service) Create(ctx context.Context, req *dto.CreateScanJobRequest, creatorUserCode, creatorDeptCode string) (*entity.ScanJob, error) {
	policy := strings.TrimSpace(req.Policy)
	if policy != PolicyQuick && policy != PolicyStandard && policy != PolicyDeep {
		return nil, apperrors.New(http.StatusBadRequest, "invalid_policy", "扫描策略仅支持 quick/standard/deep")
	}
	maxDepth := req.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 3
	}
	maxPages := req.MaxPages
	if maxPages <= 0 {
		maxPages = 100
	}

	source := strings.TrimSpace(req.Source)
	var entryURL, assetCode, deptCode string
	tempURL := false

	switch source {
	case SourceAsset:
		assetCode = strings.TrimSpace(req.AssetCode)
		if assetCode == "" {
			return nil, apperrors.New(http.StatusBadRequest, "asset_code_required", "选择资产时必须提供资产编码")
		}
		asset, err := s.assetRepo.GetByCode(ctx, assetCode)
		if err != nil {
			return nil, err
		}
		if asset == nil {
			return nil, apperrors.New(http.StatusNotFound, "asset_not_found", "资产不存在")
		}
		if asset.Status != assetsvc.StatusActive {
			return nil, apperrors.New(http.StatusBadRequest, "asset_disabled", "停用资产不可作为新任务目标")
		}
		entryURL = asset.RootURL
		deptCode = asset.DeptCode
	case SourceTempURL:
		normalized, err := assetsvc.NormalizeRootURL(req.EntryURL)
		if err != nil {
			return nil, apperrors.New(http.StatusBadRequest, "invalid_entry_url", err.Error())
		}
		entryURL = normalized
		deptCode = strings.TrimSpace(creatorDeptCode)
		if deptCode == "" {
			deptCode = "UNASSIGNED"
		}
		tempURL = true
	default:
		return nil, apperrors.New(http.StatusBadRequest, "invalid_source", "目标来源仅支持 asset 或 temp_url")
	}

	code, err := s.generateJobCode(ctx)
	if err != nil {
		return nil, err
	}
	job := &entity.ScanJob{
		JobCode:   code,
		AssetCode: assetCode,
		EntryURL:  entryURL,
		DeptCode:  deptCode,
		CreatedBy: creatorUserCode,
		Policy:    policy,
		MaxDepth:  maxDepth,
		MaxPages:  maxPages,
		Status:    StatusQueued,
		TempURL:   tempURL,
	}
	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}
	if s.queue != nil {
		if err := s.queue.Enqueue(ctx, code); err != nil {
			return nil, err
		}
	}
	return s.repo.GetByCode(ctx, code)
}

func (s *Service) GetByCode(ctx context.Context, jobCode string) (*entity.ScanJob, error) {
	return s.repo.GetByCode(ctx, jobCode)
}

func (s *Service) ListWithScope(ctx context.Context, page, pageSize int, status, policy, keyword string, scope *serviceinterfaces.UserDataScope) ([]*entity.ScanJob, int64, error) {
	filter := repointerfaces.ScanJobListFilter{
		Status:  strings.TrimSpace(status),
		Policy:  strings.TrimSpace(policy),
		Keyword: keyword,
		All:     scope == nil || scope.All,
	}
	if !filter.All && scope != nil {
		filter.DeptCodes = append([]string{}, scope.DeptCodes...)
		filter.CreatedBy = scope.UserCode
	}
	return s.repo.List(ctx, page, pageSize, filter)
}

func (s *Service) Cancel(ctx context.Context, jobCode string) (*entity.ScanJob, error) {
	job, err := s.repo.GetByCode(ctx, jobCode)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}
	if job.Status != StatusQueued && job.Status != StatusRunning {
		return nil, apperrors.New(http.StatusBadRequest, "job_not_cancellable", "仅排队或执行中的任务可取消")
	}
	now := time.Now()
	job.Status = StatusCancelled
	job.FinishedAt = &now
	if err := s.repo.Update(ctx, job); err != nil {
		return nil, err
	}
	return s.repo.GetByCode(ctx, jobCode)
}

func (s *Service) generateJobCode(ctx context.Context) (string, error) {
	for i := 0; i < 8; i++ {
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			return "", err
		}
		code := fmt.Sprintf("JOB-%s-%s", time.Now().Format("150405"), hex.EncodeToString(b[:]))
		existing, err := s.repo.GetByCode(ctx, code)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return code, nil
		}
	}
	return "", apperrors.New(http.StatusInternalServerError, "job_code_generation_failed", "生成任务编码失败，请重试")
}
