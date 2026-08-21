package operation_log

import (
	"context"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/model/vo"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

// Service implements OperationLogService.
type Service struct {
	repo    repointerfaces.OperationLogRepository
	alerter serviceinterfaces.EmailAlertService
}

// New creates a new OperationLogService.
func New(repo repointerfaces.OperationLogRepository, alerter serviceinterfaces.EmailAlertService) serviceinterfaces.OperationLogService {
	return &Service{repo: repo, alerter: alerter}
}

// Create writes a new operation log, then asynchronously notifies 邮件告警 (ADR-0005).
func (s *Service) Create(ctx context.Context, log *entity.OperationLog) error {
	if log == nil {
		return nil
	}
	if err := s.repo.Create(ctx, log); err != nil {
		return err
	}
	if log.Status != 0 || s.alerter == nil {
		return nil
	}
	entry := *log
	go func() {
		_ = s.alerter.NotifyFailedOperation(context.WithoutCancel(ctx), &entry)
	}()
	return nil
}

// Page queries operation logs with pagination.
func (s *Service) Page(ctx context.Context, req *dto.OperationLogPageRequest) (*vo.PageResult[vo.OperationLogVO], error) {
	return s.PageWithScope(ctx, req, nil)
}

func (s *Service) PageWithScope(ctx context.Context, req *dto.OperationLogPageRequest, scope *serviceinterfaces.UserDataScope) (*vo.PageResult[vo.OperationLogVO], error) {
	if req == nil {
		req = &dto.OperationLogPageRequest{}
	}
	filter := &repointerfaces.OperationLogFilter{
		UserCode:  req.UserCode,
		Username:  req.Username,
		Module:    req.Module,
		Action:    req.Action,
		Status:    req.Status,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	if scope != nil && !scope.All {
		filter.ScopeRestricted = true
		filter.AllowedUserCode = scope.UserCode
		filter.AllowedDeptCodes = scope.DeptCodes
	}
	logs, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*vo.OperationLogVO, 0, len(logs))
	for _, l := range logs {
		items = append(items, vo.BuildOperationLogVO(l))
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return &vo.PageResult[vo.OperationLogVO]{
		Items: items,
		Page:  req.Page,
		Size:  req.PageSize,
		Total: total,
	}, nil
}
