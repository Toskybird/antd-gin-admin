package interfaces

import (
	"context"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/vo"
)

// AuthService defines authentication behaviors.
type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*vo.TokenVO, error)
}

