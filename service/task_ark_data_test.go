package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func markArkVideoTask(task *model.Task) {
	task.Action = "text_to_video"
	task.PrivateData.Execution = &model.TaskExecutionSnapshot{
		TaskPlugin: &model.TaskPluginSnapshot{Key: arkVideoTaskPluginKey},
	}
}

func decodeTaskDataObject(t *testing.T, data []byte) map[string]any {
	t.Helper()
	object := map[string]any{}
	require.NoError(t, common.Unmarshal(data, &object))
	return object
}

func TestIsArkVideoTask(t *testing.T) {
	task := &model.Task{}
	assert.False(t, IsArkVideoTask(task))
	markArkVideoTask(task)
	assert.True(t, IsArkVideoTask(task))
	task.Action = "image_to_video"
	assert.True(t, IsArkVideoTask(task))
	task.Action = "generate"
	assert.False(t, IsArkVideoTask(task), "image tasks served by the same plugin are not Ark video tasks")
	markArkVideoTask(task)
	task.PrivateData.Execution.TaskPlugin.Key = "kling"
	assert.False(t, IsArkVideoTask(task))
}

func TestArkTerminalTaskDataKeepsUpstreamFields(t *testing.T) {
	task := &model.Task{
		SubmitTime: 1786015600,
		Data:       []byte(`{"id":"cgt-up","model":"doubao-seedance-2-0-260128","status":"running","content":{"video_url":""},"created_at":1786015617,"error":{"code":"old","message":"old"}}`),
	}
	task.PrivateData.UpstreamTaskID = "cgt-up"
	object := decodeTaskDataObject(t, arkTerminalTaskData(task, ArkTaskStatusExpired, arkErrorCodeTaskExpired, "Task timed out", 1786099999))
	assert.Equal(t, "cgt-up", object["id"])
	assert.Equal(t, "doubao-seedance-2-0-260128", object["model"])
	assert.Equal(t, "expired", object["status"])
	assert.EqualValues(t, 1786015617, object["created_at"])
	assert.EqualValues(t, 1786099999, object["updated_at"])
	assert.Equal(t, map[string]any{"code": "task_expired", "message": "Task timed out"}, object["error"])
	assert.Contains(t, object, "content")
}

func TestArkTerminalTaskDataFillsMissingIdentity(t *testing.T) {
	task := &model.Task{SubmitTime: 1786015600, Data: []byte(`not json`)}
	task.PrivateData.UpstreamTaskID = "cgt-missing"
	task.Properties.UpstreamModelName = "doubao-seedance-2-5-260628"
	object := decodeTaskDataObject(t, arkTerminalTaskData(task, ArkTaskStatusFailed, arkErrorCodeUpstreamNotFound, "Upstream task not found", 1786015700))
	assert.Equal(t, "cgt-missing", object["id"])
	assert.Equal(t, "doubao-seedance-2-5-260628", object["model"])
	assert.EqualValues(t, 1786015600, object["created_at"])
	assert.Equal(t, "failed", object["status"])
}

func TestSweepTimedOutArkTaskPresentsExpired(t *testing.T) {
	truncate(t)
	const userID = 611
	seedUser(t, userID, 10_000)

	arkTask := makeTask(userID, 0, 1_000, 0, BillingSourceWallet, 0)
	arkTask.TaskID = "task_ark_timeout"
	arkTask.SubmitTime = 1786000000
	arkTask.Data = []byte(`{"id":"cgt-timeout","status":"running","model":"m"}`)
	arkTask.PrivateData.UpstreamTaskID = "cgt-timeout"
	markArkVideoTask(arkTask)
	require.NoError(t, model.DB.Create(arkTask).Error)

	otherTask := makeTask(userID, 0, 1_000, 0, BillingSourceWallet, 0)
	otherTask.TaskID = "task_other_timeout"
	otherTask.SubmitTime = 1786000000
	otherTask.Data = []byte(`{"state":"running"}`)
	require.NoError(t, model.DB.Create(otherTask).Error)

	previousTimeout := constant.TaskTimeoutMinutes
	constant.TaskTimeoutMinutes = 1
	t.Cleanup(func() { constant.TaskTimeoutMinutes = previousTimeout })

	sweepTimedOutTasks(context.Background())

	var reloadedArk, reloadedOther model.Task
	require.NoError(t, model.DB.First(&reloadedArk, arkTask.ID).Error)
	require.NoError(t, model.DB.First(&reloadedOther, otherTask.ID).Error)
	assert.EqualValues(t, model.TaskStatusFailure, reloadedArk.Status)
	arkData := decodeTaskDataObject(t, reloadedArk.Data)
	assert.Equal(t, "expired", arkData["status"])
	assert.Equal(t, map[string]any{"code": "task_expired", "message": "Task timed out"}, arkData["error"])
	assert.JSONEq(t, `{"state":"running"}`, string(reloadedOther.Data), "non-Ark tasks keep their data")
}

func TestFailTaskFromPollArkTaskPresentsFailed(t *testing.T) {
	truncate(t)
	const userID, channelID = 612, 612
	seedUser(t, userID, 10_000)
	ch := &model.Channel{Id: channelID, Type: constant.ChannelTypeDoubaoVideo, Name: "ark", Key: "sk-test", Status: common.ChannelStatusEnabled}

	task := makeTask(userID, channelID, 1_000, 0, BillingSourceWallet, 0)
	task.TaskID = "task_ark_not_found"
	task.Data = []byte(`{"id":"cgt-gone","status":"queued"}`)
	task.PrivateData.UpstreamTaskID = "cgt-gone"
	markArkVideoTask(task)
	require.NoError(t, model.DB.Create(task).Error)

	adaptor := &scriptedPollingAdaptor{statusCode: http.StatusNotFound}
	require.NoError(t, updateVideoSingleTask(context.Background(), adaptor, ch, "cgt-gone", map[string]*model.Task{"cgt-gone": task}))

	var persisted model.Task
	require.NoError(t, model.DB.First(&persisted, task.ID).Error)
	assert.EqualValues(t, model.TaskStatusFailure, persisted.Status)
	data := decodeTaskDataObject(t, persisted.Data)
	assert.Equal(t, "failed", data["status"])
	assert.Equal(t, map[string]any{"code": "upstream_task_not_found", "message": "Upstream task not found"}, data["error"])
}

func TestUpdateVideoSingleTaskSupplierFailedEnvelopePersistsArkShape(t *testing.T) {
	truncate(t)
	const userID, channelID = 613, 613
	seedUser(t, userID, 10_000)
	ch := &model.Channel{Id: channelID, Type: constant.ChannelTypeDoubaoVideo, Name: "ark", Key: "sk-test", Status: common.ChannelStatusEnabled}

	task := makeTask(userID, channelID, 1_000, 0, BillingSourceWallet, 0)
	task.TaskID = "task_ark_supplier_failed"
	task.Data = []byte(`{"id":"cgt-failed","model":"doubao-seedance-2-0-260128","status":"running","created_at":1786015617}`)
	task.PrivateData.UpstreamTaskID = "cgt-failed"
	markArkVideoTask(task)
	require.NoError(t, model.DB.Create(task).Error)

	adaptor := &scriptedPollingAdaptor{
		statusCode: http.StatusBadRequest,
		body:       []byte(`{"error":{"code":"task_not_exist","message":"task_not_exist","type":"yisu_api_error","request_id":"r"}}`),
		parse:      &relaycommon.TaskInfo{Status: model.TaskStatusFailure, Reason: "task_not_exist", Progress: "100%"},
	}
	require.NoError(t, updateVideoSingleTask(context.Background(), adaptor, ch, "cgt-failed", map[string]*model.Task{"cgt-failed": task}))

	var persisted model.Task
	require.NoError(t, model.DB.First(&persisted, task.ID).Error)
	assert.EqualValues(t, model.TaskStatusFailure, persisted.Status)
	assert.Equal(t, "task_not_exist", persisted.FailReason)
	assert.Zero(t, persisted.PrivateData.PollFailures)
	data := decodeTaskDataObject(t, persisted.Data)
	assert.Equal(t, "failed", data["status"])
	assert.Equal(t, "cgt-failed", data["id"])
	assert.Equal(t, "doubao-seedance-2-0-260128", data["model"])
	assert.Equal(t, map[string]any{"code": "task_not_exist", "message": "task_not_exist"}, data["error"])
	assert.Equal(t, 10_000+1_000, getUserQuota(t, userID), "failed task is refunded once")
}

func TestUpdateVideoSingleTaskArkFailedBodyKeepsUpstreamShape(t *testing.T) {
	truncate(t)
	const userID, channelID = 614, 614
	seedUser(t, userID, 10_000)
	ch := &model.Channel{Id: channelID, Type: constant.ChannelTypeDoubaoVideo, Name: "ark", Key: "sk-test", Status: common.ChannelStatusEnabled}

	task := makeTask(userID, channelID, 1_000, 0, BillingSourceWallet, 0)
	task.TaskID = "task_ark_failed_body"
	task.PrivateData.UpstreamTaskID = "cgt-body"
	markArkVideoTask(task)
	require.NoError(t, model.DB.Create(task).Error)

	upstream := `{"id":"cgt-body","status":"failed","error":{"code":"InternalServiceError","message":"boom"}}`
	adaptor := &scriptedPollingAdaptor{body: []byte(upstream), parse: &relaycommon.TaskInfo{Status: model.TaskStatusFailure, Reason: "boom"}}
	require.NoError(t, updateVideoSingleTask(context.Background(), adaptor, ch, "cgt-body", map[string]*model.Task{"cgt-body": task}))

	var persisted model.Task
	require.NoError(t, model.DB.First(&persisted, task.ID).Error)
	assert.JSONEq(t, upstream, string(persisted.Data))
}
