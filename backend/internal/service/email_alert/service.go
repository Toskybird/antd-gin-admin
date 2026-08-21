package email_alert

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	repo "antd-gin-admin-backend/internal/repository/interfaces"
	"antd-gin-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repo.EmailAlertRepository
}

func New(repo repo.EmailAlertRepository) interfaces.EmailAlertService {
	return &Service{repo: repo}
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
