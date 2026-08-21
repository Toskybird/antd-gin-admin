package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/vo"
	"antd-gin-admin-backend/internal/service/interfaces"
	"antd-gin-admin-backend/pkg/response"
)

func registerEmailAlertRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, alertSvc interfaces.EmailAlertService) {
	if alertSvc == nil {
		return
	}
	rg.GET("/email-alert", middleware.RequireAnyPermission(permissionSvc, []string{"system:email-alert:update", "system:email-alert:list"}), func(c *gin.Context) {
		getEmailAlert(c, alertSvc)
	})
	rg.PUT("/email-alert", middleware.RequirePermission(permissionSvc, "system:email-alert:update"), func(c *gin.Context) {
		saveEmailAlert(c, alertSvc)
	})
	rg.GET("/email-alert/records", middleware.RequirePermission(permissionSvc, "system:email-alert:list"), func(c *gin.Context) {
		listEmailAlertRecords(c, alertSvc)
	})
}

func getEmailAlert(c *gin.Context, alertSvc interfaces.EmailAlertService) {
	cfg, err := alertSvc.Get(c)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildEmailAlertVO(cfg))
}

func saveEmailAlert(c *gin.Context, alertSvc interfaces.EmailAlertService) {
	var req dto.SaveEmailAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	cfg, err := alertSvc.Save(c, &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildEmailAlertVO(cfg))
}

func listEmailAlertRecords(c *gin.Context, alertSvc interfaces.EmailAlertService) {
	var req dto.EmailAlertRecordPageRequest
	req.Page = 1
	req.PageSize = 10
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	result, err := alertSvc.PageRecords(c, &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if result == nil {
		result = &vo.PageResult[vo.EmailAlertRecordVO]{
			Items: []*vo.EmailAlertRecordVO{},
			Page:  req.Page,
			Size:  req.PageSize,
			Total: 0,
		}
	}
	response.Success(c, result)
}
