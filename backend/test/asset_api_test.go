package test

import (
	"bytes"
	"context"
	"encoding/json"
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
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	assetsvc "antd-gin-admin-backend/internal/service/asset"
	authsvc "antd-gin-admin-backend/internal/service/auth"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

type fixedDataScope struct {
	scope *serviceinterfaces.UserDataScope
}

func (f *fixedDataScope) ResolveUserScope(_ context.Context, _ string, _ bool) (*serviceinterfaces.UserDataScope, error) {
	return f.scope, nil
}

func setupAssetRouter(t *testing.T, scope *serviceinterfaces.UserDataScope) *gin.Engine {
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

	resolved := scope
	if resolved == nil {
		resolved = &serviceinterfaces.UserDataScope{All: true}
	}
	dataScopeSvc := &fixedDataScope{scope: resolved}

	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	scan.RegisterRoutes(api, authMiddleware, nil, dataScopeSvc, assetSvc)
	return router
}

func loginAdminToken(t *testing.T, router *gin.Engine) string {
	t.Helper()
	loginData := map[string]string{"username": "admin", "password": "admin123"}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/system/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp["data"].(map[string]interface{})["token"].(string)
}

func doJSON(t *testing.T, router *gin.Engine, method, path, token string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var body *bytes.Buffer
	if payload != nil {
		b, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewBuffer(b)
	} else {
		body = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func parseData(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	if resp["data"] == nil {
		return nil
	}
	return resp["data"].(map[string]interface{})
}

func TestAssetListUnauthorized(t *testing.T) {
	router := setupAssetRouter(t, nil)
	w := doJSON(t, router, http.MethodGet, "/api/v1/scan/assets", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAssetCreateSuccess(t *testing.T) {
	router := setupAssetRouter(t, nil)
	token := loginAdminToken(t, router)

	w := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name":      "门户官网",
		"root_url":  "HTTPS://Portal.Example.com/",
		"dept_code": "DEPT-SEC",
		"remark":    "季度巡检",
	})
	require.Equal(t, http.StatusOK, w.Code)
	data := parseData(t, w)
	require.NotNil(t, data)
	assert.NotEmpty(t, data["asset_code"])
	assert.Equal(t, "门户官网", data["name"])
	assert.Equal(t, "https://portal.example.com", data["root_url"])
	assert.Equal(t, "DEPT-SEC", data["dept_code"])
	assert.Equal(t, "季度巡检", data["remark"])
	assert.Equal(t, "active", data["status"])
}

func TestAssetCreateDuplicateRootURLConflict(t *testing.T) {
	router := setupAssetRouter(t, nil)
	token := loginAdminToken(t, router)
	payload := map[string]string{
		"name":      "A",
		"root_url":  "https://dup.example.com",
		"dept_code": "DEPT-SEC",
	}
	w1 := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, payload)
	require.Equal(t, http.StatusOK, w1.Code)

	w2 := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name":      "B",
		"root_url":  "https://dup.example.com/",
		"dept_code": "DEPT-SEC",
	})
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestAssetUpdateFields(t *testing.T) {
	router := setupAssetRouter(t, nil)
	token := loginAdminToken(t, router)
	created := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name":      "旧名称",
		"root_url":  "https://old.example.com",
		"dept_code": "DEPT-A",
	})
	code := parseData(t, created)["asset_code"].(string)

	name := "新名称"
	root := "https://new.example.com/app"
	dept := "DEPT-B"
	remark := "备注"
	w := doJSON(t, router, http.MethodPut, "/api/v1/scan/assets/"+code, token, map[string]interface{}{
		"name":      name,
		"root_url":  root,
		"dept_code": dept,
		"remark":    remark,
	})
	require.Equal(t, http.StatusOK, w.Code)
	data := parseData(t, w)
	assert.Equal(t, "新名称", data["name"])
	assert.Equal(t, "https://new.example.com/app", data["root_url"])
	assert.Equal(t, "DEPT-B", data["dept_code"])
	assert.Equal(t, "备注", data["remark"])
}

func TestAssetDisableAndEnable(t *testing.T) {
	router := setupAssetRouter(t, nil)
	token := loginAdminToken(t, router)
	created := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name":      "可停用",
		"root_url":  "https://toggle.example.com",
		"dept_code": "DEPT-A",
	})
	code := parseData(t, created)["asset_code"].(string)

	disabled := "disabled"
	w := doJSON(t, router, http.MethodPut, "/api/v1/scan/assets/"+code, token, map[string]interface{}{"status": disabled})
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "disabled", parseData(t, w)["status"])

	active := "active"
	w2 := doJSON(t, router, http.MethodPut, "/api/v1/scan/assets/"+code, token, map[string]interface{}{"status": active})
	require.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "active", parseData(t, w2)["status"])
}

func TestAssetSoftDeleteAllowsReregister(t *testing.T) {
	router := setupAssetRouter(t, nil)
	token := loginAdminToken(t, router)
	payload := map[string]string{
		"name":      "软删",
		"root_url":  "https://soft.example.com",
		"dept_code": "DEPT-A",
	}
	created := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, payload)
	code := parseData(t, created)["asset_code"].(string)

	del := doJSON(t, router, http.MethodDelete, "/api/v1/scan/assets/"+code, token, nil)
	require.Equal(t, http.StatusOK, del.Code)

	list := doJSON(t, router, http.MethodGet, "/api/v1/scan/assets?keyword=soft.example.com", token, nil)
	require.Equal(t, http.StatusOK, list.Code)
	listData := parseData(t, list)
	assert.Equal(t, float64(0), listData["total"])

	recreate := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, payload)
	assert.Equal(t, http.StatusOK, recreate.Code)
	assert.NotEqual(t, code, parseData(t, recreate)["asset_code"])
}

func TestAssetListKeywordStatusAndDataScope(t *testing.T) {
	router := setupAssetRouter(t, &serviceinterfaces.UserDataScope{
		DeptCodes: []string{"DEPT-VISIBLE"},
	})
	token := loginAdminToken(t, router)

	doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name": "可见门户", "root_url": "https://visible.example.com", "dept_code": "DEPT-VISIBLE",
	})
	doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name": "隐藏门户", "root_url": "https://hidden.example.com", "dept_code": "DEPT-HIDDEN",
	})
	disabledName := "停用门户"
	created := doJSON(t, router, http.MethodPost, "/api/v1/scan/assets", token, map[string]string{
		"name": disabledName, "root_url": "https://disabled.example.com", "dept_code": "DEPT-VISIBLE",
	})
	code := parseData(t, created)["asset_code"].(string)
	doJSON(t, router, http.MethodPut, "/api/v1/scan/assets/"+code, token, map[string]interface{}{"status": "disabled"})

	list := doJSON(t, router, http.MethodGet, "/api/v1/scan/assets?keyword=门户&status=active", token, nil)
	require.Equal(t, http.StatusOK, list.Code)
	data := parseData(t, list)
	assert.Equal(t, float64(1), data["total"])
	items := data["items"].([]interface{})
	require.Len(t, items, 1)
	assert.Equal(t, "可见门户", items[0].(map[string]interface{})["name"])
}
