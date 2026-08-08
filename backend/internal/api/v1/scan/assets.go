package scan

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/vo"
	"antd-gin-admin-backend/internal/service/interfaces"
	"antd-gin-admin-backend/pkg/jwtutil"
	"antd-gin-admin-backend/pkg/response"
)

// RegisterRoutes wires scan module routes under /scan.
func RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, assetSvc interfaces.AssetService, detectionRuleSvc interfaces.DetectionRuleService, scanJobSvc interfaces.ScanJobService, findingSvc interfaces.FindingService, reportSvc interfaces.ScanReportService) {
	if assetSvc == nil && detectionRuleSvc == nil && scanJobSvc == nil && findingSvc == nil && reportSvc == nil {
		return
	}
	r := rg.Group("/scan")
	if authMiddleware != nil {
		r.Use(authMiddleware)
	}
	if assetSvc != nil {
		registerAssetRoutes(r, permissionSvc, dataScopeSvc, assetSvc)
	}
	if detectionRuleSvc != nil {
		registerDetectionRuleRoutes(r, permissionSvc, detectionRuleSvc)
	}
	if scanJobSvc != nil {
		registerScanJobRoutes(r, permissionSvc, dataScopeSvc, scanJobSvc)
	}
	if findingSvc != nil {
		registerFindingRoutes(r, permissionSvc, dataScopeSvc, findingSvc)
	}
	if reportSvc != nil {
		registerScanReportRoutes(r, permissionSvc, dataScopeSvc, reportSvc)
	}
}

func registerAssetRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, assetSvc interfaces.AssetService) {
	rg.GET("/assets", middleware.RequirePermission(permissionSvc, "scan:asset:list"), func(c *gin.Context) {
		listAssets(c, dataScopeSvc, assetSvc)
	})
	rg.GET("/assets/:code", middleware.RequirePermission(permissionSvc, "scan:asset:list"), func(c *gin.Context) {
		getAsset(c, assetSvc)
	})
	rg.POST("/assets", middleware.RequirePermission(permissionSvc, "scan:asset:create"), func(c *gin.Context) {
		createAsset(c, assetSvc)
	})
	rg.PUT("/assets/:code", middleware.RequirePermission(permissionSvc, "scan:asset:update"), func(c *gin.Context) {
		updateAsset(c, assetSvc)
	})
	rg.DELETE("/assets/:code", middleware.RequirePermission(permissionSvc, "scan:asset:delete"), func(c *gin.Context) {
		deleteAsset(c, assetSvc)
	})
}

func listAssets(c *gin.Context, dataScopeSvc interfaces.DataScopeService, assetSvc interfaces.AssetService) {
	var req dto.ListAssetRequest
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
	items, total, err := assetSvc.ListWithScope(c.Request.Context(), page, pageSize, req.Keyword, req.Status, scope)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	vos := make([]*vo.AssetVO, 0, len(items))
	for _, item := range items {
		vos = append(vos, vo.BuildAssetVO(item))
	}
	response.Success(c, vo.PageResult[vo.AssetVO]{Items: vos, Page: page, Size: pageSize, Total: total})
}

func getAsset(c *gin.Context, assetSvc interfaces.AssetService) {
	asset, err := assetSvc.GetByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if asset == nil {
		response.Error(c, http.StatusNotFound, "资产不存在")
		return
	}
	response.Success(c, vo.BuildAssetVO(asset))
}

func createAsset(c *gin.Context, assetSvc interfaces.AssetService) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	asset, err := assetSvc.Create(c.Request.Context(), &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, vo.BuildAssetVO(asset))
}

func updateAsset(c *gin.Context, assetSvc interfaces.AssetService) {
	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	asset, err := assetSvc.Update(c.Request.Context(), c.Param("code"), &req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if asset == nil {
		response.Error(c, http.StatusNotFound, "资产不存在")
		return
	}
	response.Success(c, vo.BuildAssetVO(asset))
}

func deleteAsset(c *gin.Context, assetSvc interfaces.AssetService) {
	if err := assetSvc.Delete(c.Request.Context(), c.Param("code")); err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func writeServiceError(c *gin.Context, err error) {
	if appErr, ok := apperrors.Extract(err); ok {
		response.ErrorWithCode(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}
	response.ErrorWithCode(c, http.StatusInternalServerError, "internal_server_error", err.Error())
}
