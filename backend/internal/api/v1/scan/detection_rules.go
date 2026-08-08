package scan

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/vo"
	"antd-gin-admin-backend/internal/service/interfaces"
	"antd-gin-admin-backend/pkg/response"
)

func registerDetectionRuleRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, svc interfaces.DetectionRuleService) {
	rg.GET("/detection-rules", middleware.RequirePermission(permissionSvc, "scan:detection-rule:list"), func(c *gin.Context) {
		listDetectionRules(c, svc)
	})
	rg.POST("/detection-rules/sync", middleware.RequirePermission(permissionSvc, "scan:detection-rule:sync"), func(c *gin.Context) {
		syncDetectionRules(c, svc)
	})
	rg.PUT("/detection-rules/:code/status", middleware.RequirePermission(permissionSvc, "scan:detection-rule:update"), func(c *gin.Context) {
		updateDetectionRuleStatus(c, svc)
	})
}

func listDetectionRules(c *gin.Context, svc interfaces.DetectionRuleService) {
	var req dto.ListDetectionRuleRequest
	_ = c.ShouldBindQuery(&req)
	page, pageSize := req.Page, req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	items, total, err := svc.List(c.Request.Context(), page, pageSize, req.Keyword, req.Enabled)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	vos := make([]*vo.DetectionRuleVO, 0, len(items))
	for _, item := range items {
		vos = append(vos, vo.BuildDetectionRuleVO(item))
	}
	response.Success(c, vo.PageResult[vo.DetectionRuleVO]{Items: vos, Page: page, Size: pageSize, Total: total})
}

func syncDetectionRules(c *gin.Context, svc interfaces.DetectionRuleService) {
	count, err := svc.Sync(c.Request.Context())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, gin.H{"synced": count})
}

func updateDetectionRuleStatus(c *gin.Context, svc interfaces.DetectionRuleService) {
	var req dto.UpdateDetectionRuleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	rule, err := svc.SetEnabled(c.Request.Context(), c.Param("code"), req.Enabled)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if rule == nil {
		response.Error(c, http.StatusNotFound, "检测规则不存在")
		return
	}
	response.Success(c, vo.BuildDetectionRuleVO(rule))
}
