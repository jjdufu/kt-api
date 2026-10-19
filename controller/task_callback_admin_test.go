package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type taskCallbackAPIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		TaskID       string                        `json:"task_id"`
		Subscription *taskCallbackSubscriptionView `json:"subscription"`
		Events       []taskCallbackEventView       `json:"events"`
		Items        []taskCallbackEventView       `json:"items"`
		Total        int                           `json:"total"`
	} `json:"data"`
}

func TestTaskCallbackAPIs(t *testing.T) {
	database := setupTaskCallbackEnqueueDB(t)
	require.NoError(t, database.AutoMigrate(&model.User{}, &model.AuditLog{}))
	previousLogDB, previousRedis := model.LOG_DB, common.RedisEnabled
	model.LOG_DB, common.RedisEnabled = database, false
	t.Cleanup(func() { model.LOG_DB, common.RedisEnabled = previousLogDB, previousRedis })
	require.NoError(t, database.Create(&model.User{Id: 1, Username: "admin", Role: common.RoleAdminUser}).Error)

	task := newSeedanceCallbackTask(t, database, "task_cb_api", true)
	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusQueued, `{"id":"cgt-task_cb_api","status":"queued"}`))
	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusFailure, `{"id":"cgt-task_cb_api","status":"failed","error":{"code":"x","message":"y"}}`))
	events := taskCallbackEvents(t, database, task.TaskID)
	require.Len(t, events, 2)
	failed := events[1]
	require.NoError(t, database.Model(&model.TaskCallbackEvent{}).Where("id = ?", failed.ID).
		Updates(map[string]any{"state": model.TaskCallbackStateFailed, "attempts": 4, "last_error": "http_5xx"}).Error)

	engine := gin.New()
	as := func(userID, role int) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("id", userID)
			c.Set("role", role)
		}
	}
	engine.GET("/user/api/task/:task_id/callbacks", as(21, common.RoleCommonUser), GetTaskCallbacks)
	engine.GET("/other/api/task/:task_id/callbacks", as(99, common.RoleCommonUser), GetTaskCallbacks)
	engine.GET("/admin/api/task/:task_id/callbacks", as(1, common.RoleAdminUser), GetTaskCallbacks)
	engine.GET("/api/admin/task-callbacks", as(1, common.RoleAdminUser), AdminListTaskCallbackEvents)
	engine.POST("/api/admin/task-callbacks/:event_id/redeliver", as(1, common.RoleAdminUser), AdminRedeliverTaskCallbackEvent)
	call := func(method, path string) (*httptest.ResponseRecorder, taskCallbackAPIResponse) {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
		var response taskCallbackAPIResponse
		_ = common.Unmarshal(recorder.Body.Bytes(), &response)
		return recorder, response
	}

	t.Run("owner sees host, hash prefix and events, never the URL", func(t *testing.T) {
		recorder, response := call(http.MethodGet, "/user/api/task/task_cb_api/callbacks")
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		require.NotNil(t, response.Data.Subscription)
		assert.Equal(t, "hooks.example.com", response.Data.Subscription.CallbackHost)
		assert.Len(t, response.Data.Subscription.CallbackURLHash, 12)
		assert.NotContains(t, recorder.Body.String(), "/cb")
		require.Len(t, response.Data.Events, 2)
		assert.Equal(t, "failed", response.Data.Events[1].Status)
		assert.True(t, response.Data.Events[1].Redeliverable)
		assert.Contains(t, response.Data.Events[1].Payload, `"status":"failed"`)
		assert.NotContains(t, response.Data.Events[1].Payload, "cgt-task_cb_api")
	})
	t.Run("other users cannot see the task; admins can", func(t *testing.T) {
		recorder, _ := call(http.MethodGet, "/other/api/task/task_cb_api/callbacks")
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		recorder, response := call(http.MethodGet, "/admin/api/task/task_cb_api/callbacks")
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Data.Events, 2)
	})
	t.Run("admin list filters by state", func(t *testing.T) {
		recorder, response := call(http.MethodGet, "/api/admin/task-callbacks?state=failed&p=1&page_size=10")
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		assert.Equal(t, 1, response.Data.Total)
		require.Len(t, response.Data.Items, 1)
		assert.Equal(t, failed.EventID, response.Data.Items[0].EventID)
		assert.Equal(t, "hooks.example.com", response.Data.Items[0].CallbackHost)
		assert.Empty(t, response.Data.Items[0].Payload, "the list does not carry payloads")
	})
	t.Run("redeliver re-arms failed events only and is audited", func(t *testing.T) {
		recorder, _ := call(http.MethodPost, "/api/admin/task-callbacks/"+failed.EventID+"/redeliver")
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		var event model.TaskCallbackEvent
		require.NoError(t, database.First(&event, failed.ID).Error)
		assert.Equal(t, model.TaskCallbackStatePending, event.State)
		assert.Zero(t, event.Attempts)
		assert.Equal(t, failed.EventID, event.EventID)

		recorder, _ = call(http.MethodPost, "/api/admin/task-callbacks/"+failed.EventID+"/redeliver")
		assert.Equal(t, http.StatusConflict, recorder.Code)
		recorder, _ = call(http.MethodPost, "/api/admin/task-callbacks/evt_missing/redeliver")
		assert.Equal(t, http.StatusConflict, recorder.Code)

		var audits []model.AuditLog
		require.NoError(t, database.Where("action = ?", "task_callback.redeliver").Order("id").Find(&audits).Error)
		require.Len(t, audits, 3)
		assert.True(t, audits[0].Success)
		assert.False(t, audits[1].Success)
		assert.Contains(t, audits[0].Content, failed.EventID)
	})
}
