package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupSeedanceRenderDB(t *testing.T) *gorm.DB {
	t.Helper()
	previousDB := model.DB
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.Task{}))
	model.DB = database
	t.Cleanup(func() { model.DB = previousDB })
	return database
}

func seedanceStatusGET(t *testing.T, userID int, taskID string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	generation := pluginruntime.DefaultRegistry.Generation()
	require.NotNil(t, generation)
	binding, found := generation.LookupDeclaredRoute(http.MethodGet, "/doubao/api/v3/contents/generations/tasks/:task_id")
	require.True(t, found)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, service.SeedanceTaskStatusPath(taskID), nil)
	c.Params = gin.Params{{Key: "task_id", Value: taskID}}
	c.Set(pluginruntime.ContextKeyPinnedRoute, pluginruntime.PinnedRoute{Generation: generation, Plugin: binding.Plugin, Route: binding.Route})
	common.SetContextKey(c, constant.ContextKeyUserId, userID)
	middleware.PrepareTaskPluginRoute()(c)
	require.True(t, c.IsAborted())
	return recorder
}

// TestSeedanceTaskStatusRendererMatchesGET locks the GET bytes and proves the
// shared renderer used for callbacks produces exactly the same bytes.
func TestSeedanceTaskStatusRendererMatchesGET(t *testing.T) {
	database := setupSeedanceRenderDB(t)
	task := &model.Task{
		TaskID:     "task_seedance_public",
		UserId:     11,
		Platform:   constant.TaskPlatform("doubao"),
		Action:     "text_to_video",
		Status:     model.TaskStatusSuccess,
		Progress:   "100%",
		SubmitTime: 1760000000,
		CreatedAt:  1760000000,
		Data: []byte(`{"id":"cgt-private-1","model":"doubao-seedance-2-0-260128","status":"succeeded",` +
			`"content":{"video_url":"https://cdn.example/v.mp4?X-Tos-Expires=86400&X-Tos-Signature=a<b>"},` +
			`"usage":{"completion_tokens":108900,"total_tokens":108900},"created_at":1760000000,"updated_at":1760000100,` +
			`"seed":42,"resolution":"720p","ratio":"16:9","duration":5,"framespersecond":24}`),
		PrivateData: model.TaskPrivateData{UpstreamTaskID: "cgt-private-1"},
	}
	require.NoError(t, database.Create(task).Error)

	recorder := seedanceStatusGET(t, 11, task.TaskID)
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
	const expected = `{"content":{"video_url":"https://cdn.example/v.mp4?X-Tos-Expires=86400\u0026X-Tos-Signature=a\u003cb\u003e"},` +
		`"created_at":1760000000,"duration":5,"framespersecond":24,"id":"task_seedance_public",` +
		`"model":"doubao-seedance-2-0-260128","ratio":"16:9","resolution":"720p","seed":42,"status":"succeeded",` +
		`"updated_at":1760000100,"usage":{"completion_tokens":108900,"total_tokens":108900}}`
	assert.Equal(t, expected, recorder.Body.String())

	var stored model.Task
	require.NoError(t, database.Where("task_id = ?", task.TaskID).First(&stored).Error)
	rendered, err := service.RenderSeedanceTaskStatus(context.Background(), &stored)
	require.NoError(t, err)
	assert.Equal(t, recorder.Body.Bytes(), rendered)
	assert.NotContains(t, string(rendered), "cgt-private-1")
}

func TestSeedanceTaskStatusRendererMatchesGETForLocalFailure(t *testing.T) {
	database := setupSeedanceRenderDB(t)
	task := &model.Task{
		TaskID:      "task_seedance_failed",
		UserId:      12,
		Platform:    constant.TaskPlatform("doubao"),
		Status:      model.TaskStatusFailure,
		Progress:    "100%",
		SubmitTime:  1760000000,
		CreatedAt:   1760000000,
		Data:        []byte(`{"id":"cgt-private-2","model":"m","status":"expired","error":{"code":"task_expired","message":"task exceeded its execution deadline"},"created_at":1760000000,"updated_at":1760003600}`),
		PrivateData: model.TaskPrivateData{UpstreamTaskID: "cgt-private-2"},
	}
	require.NoError(t, database.Create(task).Error)
	recorder := seedanceStatusGET(t, 12, task.TaskID)
	require.Equal(t, http.StatusOK, recorder.Code)
	rendered, err := service.RenderSeedanceTaskStatus(context.Background(), task)
	require.NoError(t, err)
	assert.Equal(t, recorder.Body.Bytes(), rendered)
	assert.Contains(t, string(rendered), `"status":"expired"`)
	assert.Contains(t, string(rendered), `"id":"task_seedance_failed"`)
	assert.NotContains(t, string(rendered), "cgt-private-2")
}
