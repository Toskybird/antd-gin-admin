package role_dept

import (
	"context"

	"antd-gin-admin-backend/internal/model/entity"
	repoInterfaces "antd-gin-admin-backend/internal/repository/interfaces"
	serviceInterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo repoInterfaces.RoleDeptRepository
}

func New(repo repoInterfaces.RoleDeptRepository) serviceInterfaces.RoleDeptService {
	return &Service{repo: repo}
}

func (s *Service) GetRoleDepts(ctx context.Context, roleCode string) ([]*entity.Dept, error) {
	return s.repo.GetRoleDepts(ctx, roleCode)
}

func (s *Service) SetRoleDepts(ctx context.Context, roleCode string, deptCodes []string) error {
	return s.repo.SetRoleDepts(ctx, roleCode, deptCodes)
}

