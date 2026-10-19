package controller

import (
	"context"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTaskCallbackEnqueueDB(t *testing.T) *gorm.DB {
	t.Helper()
	previousDB := model.DB
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := database.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, database.AutoMigrate(&model.Task{}, &model.TaskCallbackSubscription{}, &model.TaskCallbackEvent{}))
	model.DB = database
	setting := system_setting.GetTaskCallbackSetting()
	previousSetting := *setting
	setting.Enabled = true
	t.Cleanup(func() {
		model.DB = previousDB
		*setting = previousSetting
		_ = sqlDB.Close()
	})
	return database
}

func newSeedanceCallbackTask(t *testing.T, database *gorm.DB, taskID string, withSubscription bool) *model.Task {
	t.Helper()
	task := &model.Task{
		TaskID:     taskID,
		UserId:     21,
		Platform:   constant.TaskPlatform("doubao"),
		Action:     "text_to_video",
		Status:     model.TaskStatusSubmitted,
		Progress:   "0%",
		SubmitTime: 1760000000,
		Data:       []byte(`{"id":"cgt-` + taskID + `"}`),
		PrivateData: model.TaskPrivateData{
			UpstreamTaskID: "cgt-" + taskID,
			Execution:      &model.TaskExecutionSnapshot{TaskPlugin: &model.TaskPluginSnapshot{Key: "doubao"}},
		},
	}
	require.NoError(t, database.Create(task).Error)
	if withSubscription {
		require.NoError(t, model.CreateTaskCallbackSubscription(&model.TaskCallbackSubscription{
			TaskID:          taskID,
			UserID:          21,
			CallbackURL:     "https://hooks.example.com/cb",
			CallbackURLHash: service.TaskCallbackURLHash("https://hooks.example.com/cb"),
			CallbackHost:    "hooks.example.com",
		}))
	}
	return task
}

func transitionSeedanceTask(t *testing.T, task *model.Task, status model.TaskStatus, data string) bool {
	t.Helper()
	from := task.Status
	task.Status = status
	task.Data = []byte(data)
	won, err := task.UpdateWithStatus(from)
	require.NoError(t, err)
	return won
}

func taskCallbackEvents(t *testing.T, database *gorm.DB, taskID string) []model.TaskCallbackEvent {
	t.Helper()
	var events []model.TaskCallbackEvent
	require.NoError(t, database.Where("task_id = ?", taskID).Order("seq").Find(&events).Error)
	return events
}

func TestTaskCallbackEnqueueOnStatusTransitions(t *testing.T) {
	database := setupTaskCallbackEnqueueDB(t)
	task := newSeedanceCallbackTask(t, database, "task_cb_flow", true)

	// Status missing before the first poll: nothing is enqueued.
	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusSubmitted, `{"id":"cgt-task_cb_flow","model":"m"}`) || true)
	assert.Empty(t, taskCallbackEvents(t, database, task.TaskID))

	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusQueued, `{"id":"cgt-task_cb_flow","model":"m","status":"queued"}`))
	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusInProgress, `{"id":"cgt-task_cb_flow","model":"m","status":"running"}`))
	// Same public status with changed data: no new event.
	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusInProgress, `{"id":"cgt-task_cb_flow","model":"m","status":"running","seed":1}`))
	require.True(t, transitionSeedanceTask(t, task, model.TaskStatusSuccess, `{"id":"cgt-task_cb_flow","model":"m","status":"succeeded","content":{"video_url":"https://cdn.example/v.mp4?a=1&b=2"}}`))

	events := taskCallbackEvents(t, database, task.TaskID)
	require.Len(t, events, 3)
	statuses := []string{events[0].Status, events[1].Status, events[2].Status}
	assert.Equal(t, []string{"queued", "running", "succeeded"}, statuses)
	assert.Equal(t, []int{1, 1, 4}, []int{events[0].MaxAttempts, events[1].MaxAttempts, events[2].MaxAttempts})
	for index, event := range events {
		assert.Equal(t, index+1, event.Seq)
		assert.NotContains(t, event.Payload, "cgt-task_cb_flow")
		assert.Contains(t, event.Payload, `"id":"task_cb_flow"`)
		assert.Len(t, event.PayloadSHA256, 64)
		assert.Regexp(t, `^evt_[0-9a-f]{32}$`, event.EventID)
	}

	// Undelivered queued/running events are superseded by the newer status.
	assert.Equal(t, []string{model.TaskCallbackStateSuperseded, model.TaskCallbackStateSuperseded, model.TaskCallbackStatePending},
		[]string{events[0].State, events[1].State, events[2].State})

	// The snapshot is byte-identical to GET at that moment.
	var stored model.Task
	require.NoError(t, database.Where("task_id = ?", task.TaskID).First(&stored).Error)
	rendered, err := service.RenderSeedanceTaskStatus(context.Background(), &stored)
	require.NoError(t, err)
	assert.Equal(t, string(rendered), events[2].Payload)

	subscription, found, err := model.GetTaskCallbackSubscriptionByTaskID(task.TaskID)
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "succeeded", subscription.LastStatus)
	assert.Equal(t, 3, subscription.LastSeq)
}

func TestTaskCallbackEnqueueSkipsLosersUnsubscribedAndCancelled(t *testing.T) {
	database := setupTaskCallbackEnqueueDB(t)

	t.Run("CAS loser does not enqueue", func(t *testing.T) {
		task := newSeedanceCallbackTask(t, database, "task_cb_loser", true)
		stale := *task
		require.True(t, transitionSeedanceTask(t, task, model.TaskStatusQueued, `{"status":"queued"}`))
		stale.Status = model.TaskStatusFailure
		stale.Data = []byte(`{"status":"failed","error":{"code":"x","message":"y"}}`)
		won, err := stale.UpdateWithStatus(model.TaskStatusSubmitted)
		require.NoError(t, err)
		assert.False(t, won)
		events := taskCallbackEvents(t, database, task.TaskID)
		require.Len(t, events, 1)
		assert.Equal(t, "queued", events[0].Status)
	})

	t.Run("task without subscription", func(t *testing.T) {
		task := newSeedanceCallbackTask(t, database, "task_cb_none", false)
		require.True(t, transitionSeedanceTask(t, task, model.TaskStatusQueued, `{"status":"queued"}`))
		assert.Empty(t, taskCallbackEvents(t, database, task.TaskID))
	})

	t.Run("cancelled is recorded but not delivered", func(t *testing.T) {
		task := newSeedanceCallbackTask(t, database, "task_cb_cancel", true)
		require.True(t, transitionSeedanceTask(t, task, model.TaskStatusFailure, `{"status":"cancelled"}`))
		assert.Empty(t, taskCallbackEvents(t, database, task.TaskID))
		subscription, _, err := model.GetTaskCallbackSubscriptionByTaskID(task.TaskID)
		require.NoError(t, err)
		assert.Equal(t, "cancelled", subscription.LastStatus)
	})

	t.Run("disabled setting does not enqueue", func(t *testing.T) {
		system_setting.GetTaskCallbackSetting().Enabled = false
		defer func() { system_setting.GetTaskCallbackSetting().Enabled = true }()
		task := newSeedanceCallbackTask(t, database, "task_cb_disabled", true)
		require.True(t, transitionSeedanceTask(t, task, model.TaskStatusQueued, `{"status":"queued"}`))
		assert.Empty(t, taskCallbackEvents(t, database, task.TaskID))
	})

	t.Run("expired and failed get their attempt budgets", func(t *testing.T) {
		expired := newSeedanceCallbackTask(t, database, "task_cb_expired", true)
		require.True(t, transitionSeedanceTask(t, expired, model.TaskStatusFailure, `{"status":"expired","error":{"code":"task_expired","message":"m"}}`))
		failed := newSeedanceCallbackTask(t, database, "task_cb_failed", true)
		require.True(t, transitionSeedanceTask(t, failed, model.TaskStatusFailure, `{"status":"failed","error":{"code":"c","message":"m"}}`))
		assert.Equal(t, 1, taskCallbackEvents(t, database, expired.TaskID)[0].MaxAttempts)
		assert.Equal(t, 4, taskCallbackEvents(t, database, failed.TaskID)[0].MaxAttempts)
	})
}

func TestTaskCallbackReconcileCatchesBulkWrites(t *testing.T) {
	database := setupTaskCallbackEnqueueDB(t)
	task := newSeedanceCallbackTask(t, database, "task_cb_reconcile", true)
	fresh := newSeedanceCallbackTask(t, database, "task_cb_reconcile_fresh", true)

	// A bulk write bypasses UpdateWithStatus (no observer).
	old := int64(1700000000)
	require.NoError(t, database.Model(&model.Task{}).Where("id = ?", task.ID).
		UpdateColumns(map[string]any{"status": "FAILURE", "data": common.RawMessage(`{"status":"failed","error":{"code":"task_failed","message":"m"}}`), "updated_at": old}).Error)
	require.NoError(t, database.Model(&model.Task{}).Where("id = ?", fresh.ID).
		UpdateColumns(map[string]any{"status": "FAILURE", "data": common.RawMessage(`{"status":"failed","error":{"code":"task_failed","message":"m"}}`)}).Error)
	assert.Empty(t, taskCallbackEvents(t, database, task.TaskID))
	assert.True(t, service.HasOpenTaskCallbackSubscriptions())

	summary, err := service.ReconcileTaskCallbacks(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, summary.Scanned)
	assert.Equal(t, 1, summary.Enqueued)
	events := taskCallbackEvents(t, database, task.TaskID)
	require.Len(t, events, 1)
	assert.Equal(t, "failed", events[0].Status)
	// Recently updated tasks are left to the transactional path.
	assert.Empty(t, taskCallbackEvents(t, database, fresh.TaskID))

	// A second pass is idempotent.
	summary, err = service.ReconcileTaskCallbacks(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, summary.Enqueued)
}
