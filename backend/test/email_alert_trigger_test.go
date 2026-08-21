package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"antd-gin-admin-backend/internal/api/v1/system"
	"antd-gin-admin-backend/internal/bootstrap"
	"antd-gin-admin-backend/internal/mailer"
	"antd-gin-admin-backend/internal/middleware"
	"antd-gin-admin-backend/internal/model/entity"
	repoDB "antd-gin-admin-backend/internal/repository/db"
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	authsvc "antd-gin-admin-backend/internal/service/auth"
	emailalertsvc "antd-gin-admin-backend/internal/service/email_alert"
	operationlogsvc "antd-gin-admin-backend/internal/service/operation_log"
)

type fakeMailSender struct {
	mu      sync.Mutex
	calls   []mailer.Message
	err     error
	gate    chan struct{}
	started chan struct{}
	once    sync.Once
}

func (f *fakeMailSender) Send(_ context.Context, _ *entity.EmailAlertConfig, msg mailer.Message) error {
	f.once.Do(func() {
		if f.started != nil {
			close(f.started)
		}
	})
	if f.gate != nil {
		<-f.gate
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, msg)
	return f.err
}

func (f *fakeMailSender) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func (f *fakeMailSender) lastTo() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.calls) == 0 {
		return nil
	}
	return f.calls[len(f.calls)-1].To
}

type triggerHarness struct {
	router *gin.Engine
	logSvc interface {
		Create(context.Context, *entity.OperationLog) error
	}
	alertSvc interface {
		NotifyFailedOperation(context.Context, *entity.OperationLog) error
	}
	sender *fakeMailSender
}

func setupTriggerHarness(t *testing.T, sender *fakeMailSender, perms []string) *triggerHarness {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared&_busy_timeout=5000"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, bootstrap.AutoMigrate(db))

	authRepo := repoMock.NewAuthMockRepository()
	authSvc := authsvc.New(authRepo, emailAlertTestSecret, emailAlertTestIssuer, 3600)
	permSvc := &stubPermissionService{perms: map[string][]string{emailAlertOperator: perms}}
	alertSvc := emailalertsvc.New(
		repoDB.NewEmailAlertRepository(db),
		repoDB.NewEmailAlertRecordRepository(db),
		sender,
	)
	logSvc := operationlogsvc.New(repoDB.NewOperationLogRepository(db), alertSvc)
	router := gin.New()
	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, middleware.Auth(emailAlertTestSecret), permSvc, nil, nil, nil, nil, nil, nil, nil, nil, nil, logSvc, nil, nil, nil, alertSvc)
	return &triggerHarness{router: router, logSvc: logSvc, alertSvc: alertSvc, sender: sender}
}

func enableAlert(t *testing.T, h *triggerHarness, recipients []string) {
	t.Helper()
	payload := completeEmailAlertPayload(true)
	if recipients != nil {
		payload["recipients"] = recipients
	}
	token := operatorToken(t)
	w := doJSON(h.router, http.MethodPut, "/api/v1/system/email-alert", token, payload)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
}

func failedLog(module, errMsg string) *entity.OperationLog {
	return &entity.OperationLog{
		UserCode:  "U-OTHER",
		Username:  "other",
		Module:    module,
		Action:    "update",
		Method:    "PUT",
		Path:      "/api/v1/system/" + module,
		Status:    0,
		ErrorMsg:  errMsg,
		CreatedAt: time.Now(),
	}
}

func listRecordItems(t *testing.T, router *gin.Engine, query string) []map[string]any {
	t.Helper()
	path := "/api/v1/system/email-alert/records"
	if query != "" {
		path += "?" + query
	}
	w := doJSON(router, http.MethodGet, path, operatorToken(t), nil)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	var envelope map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	data, _ := envelope["data"].(map[string]any)
	raw, _ := data["items"].([]any)
	items := make([]map[string]any, 0, len(raw))
	for _, it := range raw {
		m, _ := it.(map[string]any)
		items = append(items, m)
	}
	return items
}

func waitRecords(t *testing.T, router *gin.Engine, n int) []map[string]any {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	var items []map[string]any
	for time.Now().Before(deadline) {
		items = listRecordItems(t, router, "")
		if len(items) == n {
			return items
		}
		time.Sleep(15 * time.Millisecond)
	}
	require.Len(t, items, n, "timed out waiting for %d 告警记录", n)
	return items
}

func TestEmailAlert_FailedLogCreatesSuccessRecord(t *testing.T) {
	sender := &fakeMailSender{}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, []string{"ops@example.com", "oncall@example.com"})

	log := failedLog("system-user", "业务失败：权限不足")
	require.NoError(t, h.logSvc.Create(context.Background(), log))
	require.NotZero(t, log.ID)

	items := waitRecords(t, h.router, 1)
	assert.Equal(t, entity.EmailAlertSendSuccess, items[0]["send_status"])
	assert.Equal(t, float64(log.ID), items[0]["operation_log_id"])
	assert.Equal(t, "other", items[0]["username"])
	assert.Equal(t, "system-user", items[0]["module"])
	assert.Equal(t, "update", items[0]["action"])
	assert.Equal(t, log.Path, items[0]["path"])
	assert.Equal(t, "业务失败：权限不足", items[0]["error_msg"])
	assert.Equal(t, 1, sender.callCount())
	assert.Equal(t, []string{"ops@example.com", "oncall@example.com"}, sender.lastTo())
}

func TestEmailAlert_DisabledDoesNotRecordOrSend(t *testing.T) {
	sender := &fakeMailSender{}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	payload := completeEmailAlertPayload(false)
	require.Equal(t, http.StatusOK, doJSON(h.router, http.MethodPut, "/api/v1/system/email-alert", operatorToken(t), payload).Code)

	require.NoError(t, h.logSvc.Create(context.Background(), failedLog("system-user", "x")))
	time.Sleep(80 * time.Millisecond)
	assert.Empty(t, listRecordItems(t, h.router, ""))
	assert.Equal(t, 0, sender.callCount())
}

func TestEmailAlert_SuccessLogDoesNotAlert(t *testing.T) {
	sender := &fakeMailSender{}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, nil)
	okLog := failedLog("system-user", "")
	okLog.Status = 1
	require.NoError(t, h.logSvc.Create(context.Background(), okLog))
	time.Sleep(80 * time.Millisecond)
	assert.Empty(t, listRecordItems(t, h.router, ""))
	assert.Equal(t, 0, sender.callCount())
}

func TestEmailAlert_OneToOneDoesNotDuplicate(t *testing.T) {
	sender := &fakeMailSender{}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, nil)
	log := failedLog("system-user", "x")
	require.NoError(t, h.logSvc.Create(context.Background(), log))
	waitRecords(t, h.router, 1)
	require.NoError(t, h.alertSvc.NotifyFailedOperation(context.Background(), log))
	time.Sleep(50 * time.Millisecond)
	assert.Len(t, listRecordItems(t, h.router, ""), 1)
	assert.Equal(t, 1, sender.callCount())
}

func TestEmailAlert_SendFailureReasonDiffersFromOpError(t *testing.T) {
	sender := &fakeMailSender{err: errors.New("smtp timeout")}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, nil)
	log := failedLog("system-role", "业务失败：名称重复")
	require.NoError(t, h.logSvc.Create(context.Background(), log))
	items := waitRecords(t, h.router, 1)
	assert.Equal(t, entity.EmailAlertSendFailed, items[0]["send_status"])
	assert.NotEmpty(t, items[0]["send_failure_reason"])
	assert.NotEqual(t, items[0]["error_msg"], items[0]["send_failure_reason"])
	assert.Equal(t, "业务失败：名称重复", items[0]["error_msg"])
}

func TestEmailAlert_MultipleRecipientsOneRecordOnFailure(t *testing.T) {
	sender := &fakeMailSender{err: errors.New("relay denied")}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, []string{"a@example.com", "b@example.com", "c@example.com"})
	require.NoError(t, h.logSvc.Create(context.Background(), failedLog("system-dept", "x")))
	items := waitRecords(t, h.router, 1)
	assert.Equal(t, entity.EmailAlertSendFailed, items[0]["send_status"])
	assert.Equal(t, 1, sender.callCount())
	assert.Len(t, sender.lastTo(), 3)
}

func TestEmailAlert_WriteReturnsBeforeSMTPFinishes(t *testing.T) {
	started := make(chan struct{})
	gate := make(chan struct{})
	sender := &fakeMailSender{started: started, gate: gate}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, nil)

	done := make(chan struct{})
	go func() {
		require.NoError(t, h.logSvc.Create(context.Background(), failedLog("system-menu", "x")))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Create blocked on SMTP")
	}
	select {
	case <-started:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("sender was not called")
	}
	assert.Empty(t, listRecordItems(t, h.router, ""))
	close(gate)
	items := waitRecords(t, h.router, 1)
	assert.Equal(t, entity.EmailAlertSendSuccess, items[0]["send_status"])
	assert.NotEqual(t, "pending", items[0]["send_status"])
}

func TestEmailAlert_RecordsPageFilters(t *testing.T) {
	sender := &fakeMailSender{}
	h := setupTriggerHarness(t, sender, []string{"system:email-alert:update", "system:email-alert:list"})
	enableAlert(t, h, nil)
	require.NoError(t, h.logSvc.Create(context.Background(), failedLog("system-user", "u")))
	waitRecords(t, h.router, 1)
	sender.err = errors.New("boom")
	require.NoError(t, h.logSvc.Create(context.Background(), failedLog("system-role", "r")))
	waitRecords(t, h.router, 2)

	byModule := listRecordItems(t, h.router, "module=system-user")
	require.Len(t, byModule, 1)
	assert.Equal(t, "system-user", byModule[0]["module"])

	byStatus := listRecordItems(t, h.router, "send_status=failed")
	require.Len(t, byStatus, 1)
	assert.Equal(t, entity.EmailAlertSendFailed, byStatus[0]["send_status"])

	future := time.Now().Add(time.Hour).Unix()
	past := listRecordItems(t, h.router, fmt.Sprintf("end_time=%d", future))
	assert.Len(t, past, 2)
	none := listRecordItems(t, h.router, fmt.Sprintf("start_time=%d", future))
	assert.Empty(t, none)

	denied := setupTriggerHarness(t, &fakeMailSender{}, []string{"system:email-alert:update"})
	w := doJSON(denied.router, http.MethodGet, "/api/v1/system/email-alert/records", operatorToken(t), nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
