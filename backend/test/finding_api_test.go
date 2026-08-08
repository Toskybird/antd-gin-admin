package test

import (
	"context"
	"net/http"
	"net/url"
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
)

type findingFixture struct {
	router      *gin.Engine
	findingRepo *repoMock.FindingMockRepository
	jobRepo     *repoMock.ScanJobMockRepository
	ruleRepo    *repoMock.DetectionRuleMockRepository
}

func setupFindingRouter(t *testing.T, scope *serviceinterfaces.UserDataScope) *findingFixture {
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

	resolved := scope
	if resolved == nil {
		resolved = &serviceinterfaces.UserDataScope{All: true}
	}
	dataScopeSvc := &fixedDataScope{scope: resolved}

	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	scan.RegisterRoutes(api, authMiddleware, nil, dataScopeSvc, assetSvc, nil, jobSvc, findingSvc)
	return &findingFixture{router: router, findingRepo: findingRepo, jobRepo: jobRepo, ruleRepo: ruleRepo}
}

func seedFindingJob(t *testing.T, f *findingFixture, jobCode, deptCode string) {
	t.Helper()
	require.NoError(t, f.jobRepo.Create(context.Background(), &entity.ScanJob{
		JobCode:  jobCode,
		EntryURL: "https://" + jobCode + ".example.com",
		DeptCode: deptCode,
		Policy:   "standard",
		Status:   "succeeded",
	}))
}

func seedRule(t *testing.T, f *findingFixture, ruleCode, displayName string) {
	t.Helper()
	require.NoError(t, f.ruleRepo.UpsertByRuleCode(context.Background(), &entity.DetectionRule{
		RuleCode:    ruleCode,
		DisplayName: displayName,
		Enabled:     true,
	}))
}

func seedFinding(t *testing.T, f *findingFixture, code, jobCode, ruleCode, severity, title, evidence string) {
	t.Helper()
	require.NoError(t, f.findingRepo.Create(context.Background(), &entity.Finding{
		FindingCode: code,
		JobCode:     jobCode,
		RuleCode:    ruleCode,
		Severity:    severity,
		Title:       title,
		Description: title + " desc",
		Evidence:    evidence,
		Location:    "https://example.com/" + code,
	}))
}

func TestFindingListUnauthorized(t *testing.T) {
	f := setupFindingRouter(t, nil)
	w := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFindingListFilters(t *testing.T) {
	f := setupFindingRouter(t, nil)
	token := loginAdminToken(t, f.router)
	seedFindingJob(t, f, "JOB-A", "DEPT-A")
	seedFindingJob(t, f, "JOB-B", "DEPT-A")
	seedRule(t, f, "rule-xss", "Cross Site Scripting")
	seedRule(t, f, "rule-sqli", "SQL Injection")
	seedFinding(t, f, "FND-1", "JOB-A", "rule-xss", "high", "XSS on login", "matched <script>")
	seedFinding(t, f, "FND-2", "JOB-A", "rule-sqli", "critical", "SQLi on search", "OR 1=1")
	seedFinding(t, f, "FND-3", "JOB-B", "rule-xss", "medium", "Reflected XSS", "alert(1)")

	byJob := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings?job_code=JOB-A", token, nil)
	require.Equal(t, http.StatusOK, byJob.Code)
	assert.Equal(t, float64(2), parseData(t, byJob)["total"])

	byRule := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings?rule_code=rule-xss", token, nil)
	require.Equal(t, http.StatusOK, byRule.Code)
	assert.Equal(t, float64(2), parseData(t, byRule)["total"])

	bySev := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings?severity=critical", token, nil)
	require.Equal(t, http.StatusOK, bySev.Code)
	assert.Equal(t, float64(1), parseData(t, bySev)["total"])
	assert.Equal(t, "FND-2", parseData(t, bySev)["items"].([]interface{})[0].(map[string]interface{})["finding_code"])

	byKW := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings?keyword="+url.QueryEscape("login"), token, nil)
	require.Equal(t, http.StatusOK, byKW.Code)
	assert.Equal(t, float64(1), parseData(t, byKW)["total"])
	item := parseData(t, byKW)["items"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "FND-1", item["finding_code"])
	assert.Equal(t, "Cross Site Scripting", item["rule_display_name"])
}

func TestFindingDetailReadonly(t *testing.T) {
	f := setupFindingRouter(t, nil)
	token := loginAdminToken(t, f.router)
	seedFindingJob(t, f, "JOB-D", "DEPT-A")
	seedRule(t, f, "rule-headers", "Missing Security Headers")
	seedFinding(t, f, "FND-D1", "JOB-D", "rule-headers", "low", "Missing CSP", "header: content-security-policy absent")

	w := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings/FND-D1", token, nil)
	require.Equal(t, http.StatusOK, w.Code)
	data := parseData(t, w)
	assert.Equal(t, "FND-D1", data["finding_code"])
	assert.Equal(t, "JOB-D", data["job_code"])
	assert.Equal(t, "rule-headers", data["rule_code"])
	assert.Equal(t, "Missing Security Headers", data["rule_display_name"])
	assert.Equal(t, "low", data["severity"])
	assert.Equal(t, "Missing CSP", data["title"])
	assert.Equal(t, "header: content-security-policy absent", data["evidence"])
	assert.Contains(t, data["description"], "desc")
}

func TestFindingListDataScopeFollowsJob(t *testing.T) {
	f := setupFindingRouter(t, &serviceinterfaces.UserDataScope{DeptCodes: []string{"DEPT-VIS"}})
	token := loginAdminToken(t, f.router)
	seedFindingJob(t, f, "JOB-VIS", "DEPT-VIS")
	seedFindingJob(t, f, "JOB-HID", "DEPT-HID")
	seedRule(t, f, "rule-x", "Rule X")
	seedFinding(t, f, "FND-VIS", "JOB-VIS", "rule-x", "high", "Visible finding", "ev-vis")
	seedFinding(t, f, "FND-HID", "JOB-HID", "rule-x", "high", "Hidden finding", "ev-hid")

	list := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings", token, nil)
	require.Equal(t, http.StatusOK, list.Code)
	assert.Equal(t, float64(1), parseData(t, list)["total"])
	assert.Equal(t, "FND-VIS", parseData(t, list)["items"].([]interface{})[0].(map[string]interface{})["finding_code"])

	ok := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings/FND-VIS", token, nil)
	require.Equal(t, http.StatusOK, ok.Code)

	denied := doJSON(t, f.router, http.MethodGet, "/api/v1/scan/findings/FND-HID", token, nil)
	assert.Equal(t, http.StatusNotFound, denied.Code)
}
