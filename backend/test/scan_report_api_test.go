package test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"antd-gin-admin-backend/internal/api/v1/scan"
	"antd-gin-admin-backend/internal/api/v1/system"
	"antd-gin-admin-backend/internal/config"
	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/entity"
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	assetsvc "antd-gin-admin-backend/internal/service/asset"
	authsvc "antd-gin-admin-backend/internal/service/auth"
	findingsvc "antd-gin-admin-backend/internal/service/finding"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
	scanjobsvc "antd-gin-admin-backend/internal/service/scanjob"
	scanreportsvc "antd-gin-admin-backend/internal/service/scanreport"
)

type scanReportFixture struct {
	router   *gin.Engine
	jobRepo  *repoMock.ScanJobMockRepository
	findings *repoMock.FindingMockRepository
	reports  *repoMock.ScanReportMockRepository
	storage  *repoMock.MemoryObjectStorage
}

func setupScanReportRouter(t *testing.T) *scanReportFixture {
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

	findingRepo := repoMock.NewFindingMockRepository()
	ruleRepo := repoMock.NewDetectionRuleMockRepository()
	findingSvc := findingsvc.New(findingRepo, jobRepo, ruleRepo)

	reportRepo := repoMock.NewScanReportMockRepository()
	storage := repoMock.NewMemoryObjectStorage()
	reportSvc := scanreportsvc.New(reportRepo, jobRepo, findingRepo, storage)

	dataScopeSvc := &fixedDataScope{scope: &serviceinterfaces.UserDataScope{All: true}}

	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	scan.RegisterRoutes(api, authMiddleware, nil, dataScopeSvc, assetSvc, nil, jobSvc, findingSvc, reportSvc)
	return &scanReportFixture{
		router:   router,
		jobRepo:  jobRepo,
		findings: findingRepo,
		reports:  reportRepo,
		storage:  storage,
	}
}

func seedReportJob(t *testing.T, f *scanReportFixture, jobCode, status, dept string) {
	t.Helper()
	require.NoError(t, f.jobRepo.Create(context.Background(), &entity.ScanJob{
		JobCode:  jobCode,
		EntryURL: "https://" + jobCode + ".example.com",
		DeptCode: dept,
		Policy:   "standard",
		Status:   status,
	}))
}

func TestScanReportGenerateMultiVersionOnSucceeded(t *testing.T) {
	f := setupScanReportRouter(t)
	token := loginAdminToken(t, f.router)
	seedReportJob(t, f, "JOB-OK", "succeeded", "DEPT-A")
	require.NoError(t, f.findings.Create(context.Background(), &entity.Finding{
		FindingCode: "FND-R1",
		JobCode:     "JOB-OK",
		RuleCode:    "rule-x",
		Severity:    "high",
		Title:       "Demo finding",
		Evidence:    "ev",
	}))

	w1 := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-reports", token, map[string]string{"job_code": "JOB-OK"})
	require.Equal(t, http.StatusOK, w1.Code)
	r1 := parseData(t, w1)
	assert.Equal(t, float64(1), r1["version"])
	assert.Equal(t, "JOB-OK", r1["job_code"])
	assert.Equal(t, "html", r1["format"])
	assert.NotEmpty(t, r1["report_code"])
	assert.NotEmpty(t, r1["object_key"])

	w2 := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-reports", token, map[string]string{"job_code": "JOB-OK"})
	require.Equal(t, http.StatusOK, w2.Code)
	r2 := parseData(t, w2)
	assert.Equal(t, float64(2), r2["version"])
	assert.NotEqual(t, r1["report_code"], r2["report_code"])

	list := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/scan-reports?job_code=JOB-OK", token, nil)
	require.Equal(t, http.StatusOK, list.Code)
	assert.Equal(t, float64(2), parseData(t, list)["total"])
}

func TestScanReportGenerateOnCancelled(t *testing.T) {
	f := setupScanReportRouter(t)
	token := loginAdminToken(t, f.router)
	seedReportJob(t, f, "JOB-CANCEL", "cancelled", "DEPT-A")

	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-reports", token, map[string]string{"job_code": "JOB-CANCEL"})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, float64(1), parseData(t, w)["version"])
}

func TestScanReportGenerateRejectedOnFailed(t *testing.T) {
	f := setupScanReportRouter(t)
	token := loginAdminToken(t, f.router)
	seedReportJob(t, f, "JOB-FAIL", "failed", "DEPT-A")

	w := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-reports", token, map[string]string{"job_code": "JOB-FAIL"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScanReportDownloadViaFakeS3(t *testing.T) {
	f := setupScanReportRouter(t)
	token := loginAdminToken(t, f.router)
	seedReportJob(t, f, "JOB-DL", "succeeded", "DEPT-A")
	require.NoError(t, f.findings.Create(context.Background(), &entity.Finding{
		FindingCode: "FND-DL",
		JobCode:     "JOB-DL",
		RuleCode:    "rule-y",
		Severity:    "medium",
		Title:       "Downloadable finding",
		Evidence:    "body-match",
	}))

	created := doJSON(t, f.router, http.MethodPost, "/api/v1/scan/scan-reports", token, map[string]string{"job_code": "JOB-DL"})
	require.Equal(t, http.StatusOK, created.Code)
	code := parseData(t, created)["report_code"].(string)
	objectKey := parseData(t, created)["object_key"].(string)

	stored, err := f.storage.Get(context.Background(), objectKey)
	require.NoError(t, err)
	assert.Contains(t, string(stored), "JOB-DL")
	assert.Contains(t, string(stored), "Downloadable finding")

	w := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/scan-reports/"+code+"/download", token, nil)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, w.Body.String(), "扫描报告")
	assert.Contains(t, w.Body.String(), "Downloadable finding")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "JOB-DL-v1.html")
}
