package email_alert

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/mailer"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/model/vo"
	repo "antd-gin-admin-backend/internal/repository/interfaces"
	"antd-gin-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo    repo.EmailAlertRepository
	records repo.EmailAlertRecordRepository
	sender  mailer.Sender
}

func New(repo repo.EmailAlertRepository, records repo.EmailAlertRecordRepository, sender mailer.Sender) interfaces.EmailAlertService {
	return &Service{repo: repo, records: records, sender: sender}
}

func (s *Service) Get(ctx context.Context) (*entity.EmailAlertConfig, error) {
	return s.repo.Get(ctx)
}

func (s *Service) Save(ctx context.Context, req *dto.SaveEmailAlertRequest) (*entity.EmailAlertConfig, error) {
	if req == nil {
		req = &dto.SaveEmailAlertRequest{}
	}
	existing, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	cfg := &entity.EmailAlertConfig{}
	if existing != nil {
		*cfg = *existing
	}
	cfg.Enabled = req.Enabled
	cfg.Host = req.Host
	cfg.Port = req.Port
	cfg.Encryption = req.Encryption
	cfg.Username = req.Username
	cfg.SenderName = req.SenderName
	cfg.SenderAddress = req.SenderAddress
	if req.Password != nil && *req.Password != "" {
		cfg.Password = *req.Password
	}
	recipients := req.Recipients
	if recipients == nil {
		recipients = []string{}
	}
	raw, err := json.Marshal(recipients)
	if err != nil {
		return nil, err
	}
	cfg.RecipientsJSON = string(raw)
	if req.Enabled && !isReadyToEnable(cfg, recipients) {
		return nil, apperrors.New(http.StatusBadRequest, "email_alert_incomplete", "须配齐邮件服务器、发件人与至少一个收件人后才能启用")
	}
	if err := s.repo.Save(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *Service) NotifyFailedOperation(ctx context.Context, log *entity.OperationLog) error {
	if log == nil || log.ID == 0 || log.Status != 0 {
		return nil
	}
	if s.records == nil || s.sender == nil {
		return nil
	}
	cfg, err := s.repo.Get(ctx)
	if err != nil {
		return err
	}
	recipients := parseRecipients(cfg)
	if cfg == nil || !cfg.Enabled || !isReadyToEnable(cfg, recipients) {
		return nil
	}
	existing, err := s.records.GetByOperationLogID(ctx, log.ID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}
	occurred := log.CreatedAt
	if occurred.IsZero() {
		occurred = time.Now()
	}
	rec := &entity.EmailAlertRecord{
		OperationLogID: log.ID,
		Username:       log.Username,
		Module:         log.Module,
		Action:         log.Action,
		Path:           log.Path,
		ErrorMsg:       log.ErrorMsg,
		OccurredAt:     occurred,
	}
	sendErr := s.sender.Send(ctx, cfg, mailer.Message{
		FromName:    cfg.SenderName,
		FromAddress: cfg.SenderAddress,
		To:          recipients,
		Subject:     fmt.Sprintf("操作失败告警：%s %s", log.Module, log.Action),
		Body:        formatAlertBody(log),
	})
	if sendErr != nil {
		rec.SendStatus = entity.EmailAlertSendFailed
		rec.SendFailureReason = sendErr.Error()
	} else {
		rec.SendStatus = entity.EmailAlertSendSuccess
		rec.SendFailureReason = ""
	}
	if err := s.records.Create(ctx, rec); err != nil {
		if apperrors.IsUniqueConstraint(err) {
			return nil
		}
		return err
	}
	return nil
}

func (s *Service) PageRecords(ctx context.Context, req *dto.EmailAlertRecordPageRequest) (*vo.PageResult[vo.EmailAlertRecordVO], error) {
	if s.records == nil {
		return &vo.PageResult[vo.EmailAlertRecordVO]{Items: []*vo.EmailAlertRecordVO{}, Page: 1, Size: 10}, nil
	}
	if req == nil {
		req = &dto.EmailAlertRecordPageRequest{}
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
	items, total, err := s.records.List(ctx, &repo.EmailAlertRecordFilter{
		SendStatus: req.SendStatus,
		Module:     req.Module,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	vos := make([]*vo.EmailAlertRecordVO, 0, len(items))
	for _, rec := range items {
		vos = append(vos, vo.BuildEmailAlertRecordVO(rec))
	}
	return &vo.PageResult[vo.EmailAlertRecordVO]{
		Items: vos,
		Page:  req.Page,
		Size:  req.PageSize,
		Total: total,
	}, nil
}

func parseRecipients(cfg *entity.EmailAlertConfig) []string {
	if cfg == nil || cfg.RecipientsJSON == "" {
		return []string{}
	}
	var recipients []string
	_ = json.Unmarshal([]byte(cfg.RecipientsJSON), &recipients)
	if recipients == nil {
		return []string{}
	}
	return recipients
}

func formatAlertBody(log *entity.OperationLog) string {
	return fmt.Sprintf("用户：%s\n模块：%s\n动作：%s\n路径：%s\n错误信息：%s\n发生时间：%s\n",
		log.Username, log.Module, log.Action, log.Path, log.ErrorMsg, log.CreatedAt.Format(time.RFC3339))
}

func isReadyToEnable(cfg *entity.EmailAlertConfig, recipients []string) bool {
	if cfg.Host == "" || cfg.Port <= 0 || strings.TrimSpace(cfg.Encryption) == "" {
		return false
	}
	if strings.TrimSpace(cfg.Username) == "" || cfg.Password == "" {
		return false
	}
	if strings.TrimSpace(cfg.SenderAddress) == "" {
		return false
	}
	for _, r := range recipients {
		if strings.TrimSpace(r) != "" {
			return true
		}
	}
	return false
}
