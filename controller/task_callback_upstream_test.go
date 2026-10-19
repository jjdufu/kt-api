package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleTaskUpstreamCallbackEndpoint(t *testing.T) {
	database := setupTaskCallbackEnqueueDB(t)
	task := newSeedanceCallbackTask(t, database, "task_upstream_endpoint", false)
	token, err := service.NewTaskUpstreamCallbackToken()
	require.NoError(t, err)
	hash := service.TaskUpstreamCallbackTokenHash(token)
	require.NoError(t, model.CreateTaskCallbackSubscription(&model.TaskCallbackSubscription{
		TaskID: task.TaskID, UserID: task.UserId, UpstreamTokenHash: &hash, UpstreamTokenExpiresAt: common.GetTimestamp() + 3600,
	}))
	// Mark the task terminal so the triggered refresh is a no-op.
	require.NoError(t, database.Model(&model.Task{}).Where("id = ?", task.ID).Update("status", model.TaskStatusSuccess).Error)

	engine := gin.New()
	engine.POST("/api/task-callback/upstream/:token", HandleTaskUpstreamCallback)
	post := func(path, body string) *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, strings.NewReader(body)))
		return recorder
	}

	// Fast path off (no public base): the endpoint does not exist.
	assert.Equal(t, http.StatusNotFound, post("/api/task-callback/upstream/"+token, `{}`).Code)

	system_setting.GetTaskCallbackSetting().UpstreamPublicBase = "https://kt.example"
	recorder := post("/api/task-callback/upstream/"+token, `{"id":"cgt-x","status":"succeeded"}`)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{}`, recorder.Body.String())
	assert.Equal(t, http.StatusNotFound, post("/api/task-callback/upstream/wrong", `{}`).Code)
	assert.Equal(t, http.StatusOK, post("/api/task-callback/upstream/"+token, `not json`).Code, "the body is not trusted or required")
	assert.Equal(t, http.StatusRequestEntityTooLarge, post("/api/task-callback/upstream/"+token, strings.Repeat("a", (1<<20)+1)).Code)
}
