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

func registerScanReportRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, reportSvc interfaces.ScanReportService) {
	rg.GET("/scan-reports", middleware.RequirePermission(permissionSvc, "scan:report:list"), func(c *gin.Context) {
		listScanReports(c, dataScopeSvc, reportSvc)
	})
	rg.POST("/scan-reports", middleware.RequirePermission(permissionSvc, "scan:report:create"), func(c *gin.Context) {
		createScanReport(c, reportSvc)
	})
	rg.GET("/scan-reports/:code/download", middleware.RequirePermission(permissionSvc, "scan:report:list"), func(c *gin.Context) {
		downloadScanReport(c, dataScopeSvc, reportSvc)
	})
}

func listScanReports(c *gin.Context, dataScopeSvc interfaces.DataScopeService, reportSvc interfaces.ScanReportService) {
	var req dto.ListScanReportRequest
	_ = c.ShouldBindQuery(&req)
	page, pageSize := req.Page, req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	scope := resolveScope(c, dataScopeSvc)
	if c.IsAborted() {
		return
	}
	items, total, err := reportSvc.ListWithScope(c.Request.Context(), page, pageSize, req.JobCode, scope)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	vos := make([]*vo.ScanReportVO, 0, len(items))
	for _, item := range items {
		vos = append(vos, vo.BuildScanReportVO(item))
	}
	response.Success(c, vo.PageResult[vo.ScanReportVO]{Items: vos, Page: page, Size: pageSize, Total: total})
}

func createScanReport(c *gin.Context, reportSvc interfaces.ScanReportService) {
	var req dto.CreateScanReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, _ := claimsVal.(*jwtutil.Claims)
	userCode := ""
	if claims != nil {
		userCode = claims.UserCode
	}
	report, err := reportSvc.Generate(c.Request.Context(), req.JobCode, userCode)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildScanReportVO(report))
}

func downloadScanReport(c *gin.Context, dataScopeSvc interfaces.DataScopeService, reportSvc interfaces.ScanReportService) {
	scope := resolveScope(c, dataScopeSvc)
	if c.IsAborted() {
		return
	}
	content, contentType, filename, err := reportSvc.Download(c.Request.Context(), c.Param("code"), scope)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if content == nil {
		response.Error(c, http.StatusNotFound, "扫描报告不存在")
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, content)
}

func resolveScope(c *gin.Context, dataScopeSvc interfaces.DataScopeService) *interfaces.UserDataScope {
	claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
	claims, _ := claimsVal.(*jwtutil.Claims)
	if dataScopeSvc == nil || claims == nil {
		return nil
	}
	resolved, err := dataScopeSvc.ResolveUserScope(c.Request.Context(), claims.UserCode, claims.IsSuperAdmin)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		c.Abort()
		return nil
	}
	return resolved
}
