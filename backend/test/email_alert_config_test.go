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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/api/v1/system"
	"antd-gin-admin-backend/internal/bootstrap"
	"antd-gin-admin-backend/internal/middleware"
	repoDB "antd-gin-admin-backend/internal/repository/db"
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	authsvc "antd-gin-admin-backend/internal/service/auth"
	emailalertsvc "antd-gin-admin-backend/internal/service/email_alert"
	"antd-gin-admin-backend/pkg/jwtutil"
)

const (
	emailAlertTestSecret = "test-secret"
	emailAlertTestIssuer = "test"
	emailAlertOperator   = "U-OPERATOR"
)

type stubPermissionService struct {
	perms map[string][]string
}

func (s *stubPermissionService) GetUserPermissions(_ context.Context, userCode string) ([]string, error) {
	return s.perms[userCode], nil
}

func (s *stubPermissionService) CheckPermission(_ context.Context, userCode, perm string) (bool, error) {
	for _, p := range s.perms[userCode] {
		if p == perm {
			return true, nil
		}
	}
	return false, nil
}

func completeEmailAlertPayload(enabled bool) map[string]any {
	return map[string]any{
		"enabled":        enabled,
		"host":           "smtp.example.com",
		"port":           587,
		"encryption":     "starttls",
		"username":       "alert@example.com",
		"password":       "secret",
		"sender_name":    "Ops",
		"sender_address": "alert@example.com",
		"recipients":     []string{"ops@example.com"},
	}
}

func setupEmailAlertRouter(t *testing.T, perms []string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, bootstrap.AutoMigrate(db))

	authRepo := repoMock.NewAuthMockRepository()
	authSvc := authsvc.New(authRepo, emailAlertTestSecret, emailAlertTestIssuer, 3600)
	authMiddleware := middleware.Auth(emailAlertTestSecret)

	permSvc := &stubPermissionService{perms: map[string][]string{
		emailAlertOperator: perms,
	}}
	alertSvc := emailalertsvc.New(repoDB.NewEmailAlertRepository(db))

	router := gin.New()
	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, permSvc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, alertSvc)
	return router
}

func operatorToken(t *testing.T) string {
	t.Helper()
	token, err := jwtutil.GenerateToken(emailAlertTestSecret, emailAlertTestIssuer, emailAlertOperator, "operator", false, 3600)
	require.NoError(t, err)
	return token
}

func doJSON(router *gin.Engine, method, path, token string, body any) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func parseData(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var envelope map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	data, _ := envelope["data"].(map[string]any)
	return data
}

func TestEmailAlert_EnableWhenComplete(t *testing.T) {
	router := setupEmailAlertRouter(t, []string{"system:email-alert:update"})
	token := operatorToken(t)

	putW := doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, completeEmailAlertPayload(true))
	require.Equal(t, http.StatusOK, putW.Code, putW.Body.String())

	getW := doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	require.Equal(t, http.StatusOK, getW.Code, getW.Body.String())
	data := parseData(t, getW)
	assert.Equal(t, true, data["enabled"])
}

func TestEmailAlert_CannotEnableWhenIncomplete(t *testing.T) {
	router := setupEmailAlertRouter(t, []string{"system:email-alert:update"})
	token := operatorToken(t)

	payload := completeEmailAlertPayload(true)
	payload["host"] = ""
	putW := doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, payload)
	require.NotEqual(t, http.StatusOK, putW.Code, putW.Body.String())

	getW := doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	require.Equal(t, http.StatusOK, getW.Code, getW.Body.String())
	data := parseData(t, getW)
	assert.Equal(t, false, data["enabled"])
}

func TestEmailAlert_SaveIncompleteDraftWhenDisabled(t *testing.T) {
	router := setupEmailAlertRouter(t, []string{"system:email-alert:update"})
	token := operatorToken(t)

	putW := doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, map[string]any{
		"enabled": false,
		"host":    "smtp.draft.example.com",
	})
	require.Equal(t, http.StatusOK, putW.Code, putW.Body.String())

	getW := doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	require.Equal(t, http.StatusOK, getW.Code)
	data := parseData(t, getW)
	assert.Equal(t, false, data["enabled"])
	assert.Equal(t, "smtp.draft.example.com", data["host"])
}

func TestEmailAlert_RejectIncompleteSaveWhileEnabled(t *testing.T) {
	router := setupEmailAlertRouter(t, []string{"system:email-alert:update"})
	token := operatorToken(t)

	require.Equal(t, http.StatusOK, doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, completeEmailAlertPayload(true)).Code)

	payload := completeEmailAlertPayload(true)
	payload["recipients"] = []string{}
	putW := doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, payload)
	require.NotEqual(t, http.StatusOK, putW.Code, putW.Body.String())

	getW := doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	data := parseData(t, getW)
	assert.Equal(t, true, data["enabled"])
	recipients, _ := data["recipients"].([]any)
	require.Len(t, recipients, 1)
	assert.Equal(t, "ops@example.com", recipients[0])
}

func TestEmailAlert_GetHidesPasswordAndEmptyKeepsExisting(t *testing.T) {
	router := setupEmailAlertRouter(t, []string{"system:email-alert:update"})
	token := operatorToken(t)

	require.Equal(t, http.StatusOK, doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, completeEmailAlertPayload(true)).Code)

	getW := doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	data := parseData(t, getW)
	_, hasPassword := data["password"]
	assert.False(t, hasPassword)
	assert.Equal(t, true, data["password_configured"])

	payload := completeEmailAlertPayload(true)
	payload["password"] = ""
	payload["sender_name"] = "NightOps"
	putW := doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, payload)
	require.Equal(t, http.StatusOK, putW.Code, putW.Body.String())

	getW = doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	data = parseData(t, getW)
	assert.Equal(t, true, data["enabled"])
	assert.Equal(t, true, data["password_configured"])
	assert.Equal(t, "NightOps", data["sender_name"])
	_, hasPassword = data["password"]
	assert.False(t, hasPassword)
}

func TestEmailAlert_PutRequiresUpdatePermission(t *testing.T) {
	router := setupEmailAlertRouter(t, []string{"system:email-alert:list"})
	token := operatorToken(t)

	putW := doJSON(router, http.MethodPut, "/api/v1/system/email-alert", token, completeEmailAlertPayload(true))
	assert.Equal(t, http.StatusForbidden, putW.Code, putW.Body.String())

	getW := doJSON(router, http.MethodGet, "/api/v1/system/email-alert", token, nil)
	assert.Equal(t, http.StatusOK, getW.Code, getW.Body.String())
}
