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

func registerScanJobRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, scanJobSvc interfaces.ScanJobService) {
	rg.GET("/scan-jobs", middleware.RequirePermission(permissionSvc, "scan:job:list"), func(c *gin.Context) {
		listScanJobs(c, dataScopeSvc, scanJobSvc)
	})
	rg.GET("/scan-jobs/:code", middleware.RequirePermission(permissionSvc, "scan:job:list"), func(c *gin.Context) {
		getScanJob(c, scanJobSvc)
	})
	rg.POST("/scan-jobs", middleware.RequirePermission(permissionSvc, "scan:job:create"), func(c *gin.Context) {
		createScanJob(c, scanJobSvc)
	})
	rg.POST("/scan-jobs/:code/cancel", middleware.RequirePermission(permissionSvc, "scan:job:cancel"), func(c *gin.Context) {
		cancelScanJob(c, scanJobSvc)
	})
}

func listScanJobs(c *gin.Context, dataScopeSvc interfaces.DataScopeService, scanJobSvc interfaces.ScanJobService) {
	var req dto.ListScanJobRequest
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
	items, total, err := scanJobSvc.ListWithScope(c.Request.Context(), page, pageSize, req.Status, req.Policy, req.Keyword, scope)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	vos := make([]*vo.ScanJobVO, 0, len(items))
	for _, item := range items {
		vos = append(vos, vo.BuildScanJobVO(item))
	}
	response.Success(c, vo.PageResult[vo.ScanJobVO]{Items: vos, Page: page, Size: pageSize, Total: total})
}

func getScanJob(c *gin.Context, scanJobSvc interfaces.ScanJobService) {
	job, err := scanJobSvc.GetByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if job == nil {
		response.Error(c, http.StatusNotFound, "扫描任务不存在")
		return
	}
	response.Success(c, vo.BuildScanJobVO(job))
}

func createScanJob(c *gin.Context, scanJobSvc interfaces.ScanJobService) {
	var req dto.CreateScanJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, _ := claimsVal.(*jwtutil.Claims)
	userCode, deptCode := "", ""
	if claims != nil {
		userCode = claims.UserCode
	}
	job, err := scanJobSvc.Create(c.Request.Context(), &req, userCode, deptCode)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildScanJobVO(job))
}

func cancelScanJob(c *gin.Context, scanJobSvc interfaces.ScanJobService) {
	job, err := scanJobSvc.Cancel(c.Request.Context(), c.Param("code"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if job == nil {
		response.Error(c, http.StatusNotFound, "扫描任务不存在")
		return
	}
	response.Success(c, vo.BuildScanJobVO(job))
}
