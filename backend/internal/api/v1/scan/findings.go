package scan

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/vo"
	"antd-gin-admin-backend/internal/service/interfaces"
	"antd-gin-admin-backend/pkg/jwtutil"
	"antd-gin-admin-backend/pkg/response"
)

func registerFindingRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, findingSvc interfaces.FindingService) {
	rg.GET("/findings", middleware.RequirePermission(permissionSvc, "scan:finding:list"), func(c *gin.Context) {
		listFindings(c, dataScopeSvc, findingSvc)
	})
	rg.GET("/findings/:code", middleware.RequirePermission(permissionSvc, "scan:finding:list"), func(c *gin.Context) {
		getFinding(c, dataScopeSvc, findingSvc)
	})
}

func listFindings(c *gin.Context, dataScopeSvc interfaces.DataScopeService, findingSvc interfaces.FindingService) {
	var req dto.ListFindingRequest
	_ = c.ShouldBindQuery(&req)
	page, pageSize := req.Page, req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, _ := claimsVal.(*jwtutil.Claims)
	var scope *interfaces.UserDataScope
	if dataScopeSvc != nil && claims != nil {
		resolved, err := dataScopeSvc.ResolveUserScope(c.Request.Context(), claims.UserCode, claims.IsSuperAdmin)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		scope = resolved
	}
	items, total, err := findingSvc.ListWithScope(c.Request.Context(), page, pageSize, req.JobCode, req.RuleCode, req.Severity, req.Keyword, scope)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	vos := make([]*vo.FindingVO, 0, len(items))
	for _, item := range items {
		displayName, _ := findingSvc.RuleDisplayName(c.Request.Context(), item.RuleCode)
		vos = append(vos, vo.BuildFindingVO(item, displayName))
	}
	response.Success(c, vo.PageResult[vo.FindingVO]{Items: vos, Page: page, Size: pageSize, Total: total})
}

func getFinding(c *gin.Context, dataScopeSvc interfaces.DataScopeService, findingSvc interfaces.FindingService) {
	claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, _ := claimsVal.(*jwtutil.Claims)
	var scope *interfaces.UserDataScope
	if dataScopeSvc != nil && claims != nil {
		resolved, err := dataScopeSvc.ResolveUserScope(c.Request.Context(), claims.UserCode, claims.IsSuperAdmin)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		scope = resolved
	}
	finding, err := findingSvc.GetByCodeWithScope(c.Request.Context(), c.Param("code"), scope)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if finding == nil {
		response.Error(c, http.StatusNotFound, "发现项不存在")
		return
	}
	displayName, _ := findingSvc.RuleDisplayName(c.Request.Context(), finding.RuleCode)
	response.Success(c, vo.BuildFindingVO(finding, displayName))
}
