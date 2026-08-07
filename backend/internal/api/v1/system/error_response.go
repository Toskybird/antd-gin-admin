package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/pkg/response"
)

func writeServiceError(c *gin.Context, err error) {
	if appErr, ok := apperrors.Extract(err); ok {
		response.ErrorWithCode(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}
	response.ErrorWithCode(c, http.StatusInternalServerError, "internal_server_error", err.Error())
}
