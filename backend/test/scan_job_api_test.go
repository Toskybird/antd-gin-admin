package test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"antd-gin-admin-backend/internal/api/v1/scan"
	"antd-gin-admin-backend/internal/api/v1/system"
	"antd-gin-admin-backend/internal/config"
	"antd-gin-admin-backend/internal/middleware"
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	assetsvc "antd-gin-admin-backend/internal/service/asset"
	authsvc "antd-gin-admin-backend/internal/service/auth"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
	scanjobsvc "antd-gin-admin-backend/internal/service/scanjob"
)

type scanJobFixture struct {
	router *gin.Engine
	queue  *repoMock.MemoryScanJobQueue
}

func setupScanJobRouter(t *testing.T, scope *serviceinterfaces.UserDataScope) *scanJobFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()

	authRepo := repoMock.NewAuthMockRepository()
	cfg := &config.Config{}
	cfg.JWT.Secret = "test-secret"
	cfg.JWT.ExpireTime = 3600
	cfg.JWT.Issuer = "test"
	authSvc := authsvc.New(authRepo, cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime)
	authMiddleware := middleware.Auth(cfg.JWT.Secret)

	assetRepo := repoMock.NewAssetMockRepository()
	assetSvc := assetsvc.New(assetRepo)
	jobRepo := repoMock.NewScanJobMockRepository()
	queue := repoMock.NewMemoryScanJobQueue()
	jobSvc := scanjobsvc.New(jobRepo, assetRepo, queue)

	resolved := scope
	if resolved == nil {
		resolved = &serviceinterfaces.UserDataScope{All: true}
	}
	dataScopeSvc := &fixedDataScope{scope: resolved}

	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	scan.RegisterRoutes(api, authMiddleware, nil, dataScopeSvc, assetSvc, nil, jobSvc)
	return &scanJobFixture{router: router, queue: queue}
}

func createActiveAsset(t *testing.T, f *scanJobFixture, token, name, url, dept string) string {
	t.Helper()
	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name": name, "root_url": url, "dept_code": dept,
	})
	require.Equal(t, http.StatusOK, w.Code)
	return parseData(t, w)["asset_code"].(string)
}

func TestScanJobUnauthorized(t *testing.T) {
	f := setupScanJobRouter(t, nil)
	w := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/scan-jobs", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScanJobCreateFromAssetQueuedAndEnqueued(t *testing.T) {
	f := setupScanJobRouter(t, nil)
	token := loginAdminToken(t, f.router)
	assetCode := createActiveAsset(t, f, token, "门户", "https://portal.example.com", "DEPT-SEC")

	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs", token, map[string]interface{}{
		"source":     "asset",
		"asset_code": assetCode,
		"policy":     "standard",
		"max_depth":  2,
		"max_pages":  50,
	})
	require.Equal(t, http.StatusOK, w.Code)
	data := parseData(t, w)
	assert.Equal(t, "queued", data["status"])
	assert.Equal(t, "https://portal.example.com", data["entry_url"])
	assert.Equal(t, "standard", data["policy"])
	assert.Equal(t, float64(2), data["max_depth"])
	assert.Equal(t, float64(50), data["max_pages"])
	assert.Equal(t, false, data["temp_url"])
	assert.Equal(t, float64(0), data["finding_count"])

	enqueued, err := f.queue.ListEnqueued(nil)
	require.NoError(t, err)
	require.Len(t, enqueued, 1)
	assert.Equal(t, data["job_code"], enqueued[0])
}

func TestScanJobCreateTempURL(t *testing.T) {
	f := setupScanJobRouter(t, nil)
	token := loginAdminToken(t, f.router)

	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs", token, map[string]interface{}{
		"source":    "temp_url",
		"entry_url": "https://temp.example.com/app/",
		"policy":    "quick",
	})
	require.Equal(t, http.StatusOK, w.Code)
	data := parseData(t, w)
	assert.Equal(t, true, data["temp_url"])
	assert.Equal(t, "https://temp.example.com/app", data["entry_url"])
	assert.Equal(t, "queued", data["status"])
	_, hasAsset := data["asset_code"]
	assert.False(t, hasAsset)
}

func TestScanJobRejectDisabledAsset(t *testing.T) {
	f := setupScanJobRouter(t, nil)
	token := loginAdminToken(t, f.router)
	assetCode := createActiveAsset(t, f, token, "旧站", "https://old.example.com", "DEPT-A")
	require.Equal(t, http.StatusOK, doJSON(t, f.router, http.MethodPut, "/api/v1/scan/assets/"+assetCode, token, map[string]string{"status": "disabled"}).Code)

	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs", token, map[string]interface{}{
		"source": "asset", "asset_code": assetCode, "policy": "standard",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScanJobListAndDetailWithScope(t *testing.T) {
	f := setupScanJobRouter(t, &serviceinterfaces.UserDataScope{DeptCodes: []string{"DEPT-VIS"}})
	token := loginAdminToken(t, f.router)
	a1 := createActiveAsset(t, f, token, "可见", "https://vis.example.com", "DEPT-VIS")
	a2 := createActiveAsset(t, f, token, "隐藏", "https://hid.example.com", "DEPT-HID")
	doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs", token, map[string]interface{}{
		"source": "asset", "asset_code": a1, "policy": "deep",
	})
	doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs", token, map[string]interface{}{
		"source": "asset", "asset_code": a2, "policy": "quick",
	})

	list := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/scan-jobs", token, nil)
	require.Equal(t, http.StatusOK, list.Code)
	assert.Equal(t, float64(1), parseData(t, list)["total"])
	jobCode := parseData(t, list)["items"].([]interface{})[0].(map[string]interface{})["job_code"].(string)

	detail := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/scan-jobs/"+jobCode, token, nil)
	require.Equal(t, http.StatusOK, detail.Code)
	assert.Equal(t, "deep", parseData(t, detail)["policy"])
}

func TestScanJobCancelQueued(t *testing.T) {
	f := setupScanJobRouter(t, nil)
	token := loginAdminToken(t, f.router)
	created := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs", token, map[string]interface{}{
		"source": "temp_url", "entry_url": "https://cancel.example.com", "policy": "standard",
	})
	code := parseData(t, created)["job_code"].(string)

	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-jobs/"+code+"/cancel", token, nil)
	require.Equal(t, http.StatusOK, w.Code)
	data := parseData(t, w)
	assert.Equal(t, "cancelled", data["status"])
	assert.NotNil(t, data["finished_at"])
	assert.Equal(t, float64(0), data["finding_count"])
}
