package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
	authsvc "antd-gin-admin-backend/internal/service/auth"
	detectionrulesvc "antd-gin-admin-backend/internal/service/detectionrule"
)

func setupDetectionRuleRouter(t *testing.T) *gin.Engine {
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

	repo := repoMock.NewDetectionRuleMockRepository()
	source := &repoMock.StaticDetectionRuleSource{
		Rules: []entity.DetectionRule{
			{RuleCode: "exposed-panels", DisplayName: "Exposed Admin Panels"},
			{RuleCode: "tech-detect", DisplayName: "Technology Detection"},
		},
	}
	svc := detectionrulesvc.New(repo, source)

	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	scan.RegisterRoutes(api, authMiddleware, nil, nil, nil, svc, nil, nil)
	return router
}

func TestDetectionRuleListUnauthorized(t *testing.T) {
	router := setupDetectionRuleRouter(t)
	w := doJSON(t, router, http.MethodGet, "/api/v1/scan/detection-rules", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDetectionRuleSyncAndList(t *testing.T) {
	router := setupDetectionRuleRouter(t)
	token := loginAdminToken(t, router)

	sync := doJSON(t, router, http.MethodPost, "/api/v1/scan/detection-rules/sync", token, nil)
	require.Equal(t, http.StatusOK, sync.Code)
	assert.Equal(t, float64(2), parseData(t, sync)["synced"])

	list := doJSON(t, router, http.MethodGet, "/api/v1/scan/detection-rules", token, nil)
	require.Equal(t, http.StatusOK, list.Code)
	data := parseData(t, list)
	assert.Equal(t, float64(2), data["total"])
	items := data["items"].([]interface{})
	require.Len(t, items, 2)
	first := items[0].(map[string]interface{})
	assert.NotEmpty(t, first["rule_code"])
	assert.NotEmpty(t, first["display_name"])
}

func TestDetectionRuleDisableEnable(t *testing.T) {
	router := setupDetectionRuleRouter(t)
	token := loginAdminToken(t, router)
	require.Equal(t, http.StatusOK, doJSON(t, router, http.MethodPost, "/api/v1/scan/detection-rules/sync", token, nil).Code)

	w := doJSON(t, router, http.MethodPut, "/api/v1/scan/detection-rules/exposed-panels/status", token, map[string]bool{"enabled": false})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, false, parseData(t, w)["enabled"])

	w2 := doJSON(t, router, http.MethodPut, "/api/v1/scan/detection-rules/exposed-panels/status", token, map[string]bool{"enabled": true})
	require.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, true, parseData(t, w2)["enabled"])
}

func TestDetectionRuleSyncPreservesDisabled(t *testing.T) {
	router := setupDetectionRuleRouter(t)
	token := loginAdminToken(t, router)
	require.Equal(t, http.StatusOK, doJSON(t, router, http.MethodPost, "/api/v1/scan/detection-rules/sync", token, nil).Code)
	require.Equal(t, http.StatusOK, doJSON(t, router, http.MethodPut, "/api/v1/scan/detection-rules/tech-detect/status", token, map[string]bool{"enabled": false}).Code)

	sourceRenameSync := doJSON(t, router, http.MethodPost, "/api/v1/scan/detection-rules/sync", token, nil)
	require.Equal(t, http.StatusOK, sourceRenameSync.Code)

	getList := doJSON(t, router, http.MethodGet, "/api/v1/scan/detection-rules?keyword=tech-detect", token, nil)
	require.Equal(t, http.StatusOK, getList.Code)
	items := parseData(t, getList)["items"].([]interface{})
	require.Len(t, items, 1)
	assert.Equal(t, false, items[0].(map[string]interface{})["enabled"])
}

func TestDetectionRuleRejectsCustomCreate(t *testing.T) {
	router := setupDetectionRuleRouter(t)
	token := loginAdminToken(t, router)

	w := doJSON(t, router, http.MethodPost, "/api/v1/scan/detection-rules", token, map[string]string{
		"rule_code":    "custom-rule",
		"display_name": "自定义规则",
	})
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMethodNotAllowed || w.Code == http.StatusForbidden)

	body := bytes.NewBufferString(`{"rule_code":"x","display_name":"y"}`)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/scan/detection-rules/exposed-panels", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.True(t, w2.Code == http.StatusNotFound || w2.Code == http.StatusMethodNotAllowed)
}
