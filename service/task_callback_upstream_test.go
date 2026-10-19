package service

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withTaskCallbackSetting(t *testing.T, mutate func(*system_setting.TaskCallbackSetting)) {
	t.Helper()
	setting := system_setting.GetTaskCallbackSetting()
	previous := *setting
	t.Cleanup(func() { *setting = previous })
	mutate(setting)
}

func TestPrepareTaskUpstreamCallbackModes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newContext := func() *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		return c
	}
	enabledChannel := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{ChannelOtherSettings: dto.ChannelOtherSettings{UpstreamTaskCallbackEnabled: true}}}
	disabledChannel := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{}}

	t.Run("off by default", func(t *testing.T) {
		c := newContext()
		require.NoError(t, PrepareTaskUpstreamCallback(c, true))
		assert.Empty(t, TaskUpstreamCallbackURLForChannel(c, enabledChannel))
	})
	t.Run("only when the customer supplied callback_url", func(t *testing.T) {
		withTaskCallbackSetting(t, func(s *system_setting.TaskCallbackSetting) {
			s.Enabled = true
			s.UpstreamPublicBase = "https://kt.example/"
			s.UpstreamInjectMode = system_setting.TaskCallbackUpstreamInjectOnlyWhenCustomerCallback
		})
		c := newContext()
		require.NoError(t, PrepareTaskUpstreamCallback(c, false))
		assert.Empty(t, TaskUpstreamCallbackURLForChannel(c, enabledChannel))

		c = newContext()
		require.NoError(t, PrepareTaskUpstreamCallback(c, true))
		url := TaskUpstreamCallbackURLForChannel(c, enabledChannel)
		token := c.GetString(ContextKeyTaskUpstreamCallbackToken)
		assert.Len(t, token, 43)
		assert.Equal(t, "https://kt.example/api/task-callback/upstream/"+token, url)
		assert.Empty(t, TaskUpstreamCallbackURLForChannel(c, disabledChannel), "the channel must opt in")
		assert.Equal(t, token, taskUpstreamCallbackTokenFromContext(c, enabledChannel))
		assert.Empty(t, taskUpstreamCallbackTokenFromContext(c, disabledChannel))
	})
	t.Run("always and never", func(t *testing.T) {
		withTaskCallbackSetting(t, func(s *system_setting.TaskCallbackSetting) {
			s.Enabled = true
			s.UpstreamPublicBase = "https://kt.example"
			s.UpstreamInjectMode = system_setting.TaskCallbackUpstreamInjectAlways
		})
		c := newContext()
		require.NoError(t, PrepareTaskUpstreamCallback(c, false))
		assert.NotEmpty(t, TaskUpstreamCallbackURLForChannel(c, enabledChannel))
		system_setting.GetTaskCallbackSetting().UpstreamInjectMode = system_setting.TaskCallbackUpstreamInjectNever
		c = newContext()
		require.NoError(t, PrepareTaskUpstreamCallback(c, true))
		assert.Empty(t, TaskUpstreamCallbackURLForChannel(c, enabledChannel))
	})
}

func TestTaskUpstreamCallbackTriggersRefreshThroughPolling(t *testing.T) {
	truncate(t)
	const userID, channelID = 640, 640
	seedUser(t, userID, 10_000)
	channel := &model.Channel{Id: channelID, Type: constant.ChannelTypeDoubaoVideo, Name: "ark", Key: "sk-test", Status: common.ChannelStatusEnabled}
	require.NoError(t, model.DB.Create(channel).Error)
	task := makeTask(userID, channelID, 1_000, 0, BillingSourceWallet, 0)
	task.TaskID = "task_upstream_hint"
	task.Status = model.TaskStatusQueued
	task.Data = []byte(`{"id":"cgt-hint","status":"queued"}`)
	task.PrivateData.UpstreamTaskID = "cgt-hint"
	markArkVideoTask(task)
	require.NoError(t, model.DB.Create(task).Error)

	token, err := NewTaskUpstreamCallbackToken()
	require.NoError(t, err)
	hash := TaskUpstreamCallbackTokenHash(token)
	require.NoError(t, model.CreateTaskCallbackSubscription(&model.TaskCallbackSubscription{
		TaskID: task.TaskID, UserID: userID, UpstreamTokenHash: &hash, UpstreamTokenExpiresAt: common.GetTimestamp() + 3600,
	}))

	adaptor := &scriptedPollingAdaptor{
		body:  []byte(`{"id":"cgt-hint","status":"running"}`),
		parse: &relaycommon.TaskInfo{Status: model.TaskStatusInProgress, Progress: "50%"},
	}
	previousFactory := GetTaskAdaptorFunc
	GetTaskAdaptorFunc = func(constant.TaskPlatform) TaskPollingAdaptor { return adaptor }
	t.Cleanup(func() { GetTaskAdaptorFunc = previousFactory })

	assert.False(t, HandleTaskUpstreamCallback("unknown-token", "running"))
	assert.False(t, HandleTaskUpstreamCallback("", "running"))

	// The body's claimed status is ignored; the task is re-queried.
	require.NoError(t, RefreshTaskNow(context.Background(), task.TaskID))
	var persisted model.Task
	require.NoError(t, model.DB.First(&persisted, task.ID).Error)
	assert.EqualValues(t, model.TaskStatusInProgress, persisted.Status)
	assert.Equal(t, "running", decodeTaskDataObject(t, persisted.Data)["status"])

	assert.True(t, HandleTaskUpstreamCallback(token, "succeeded"))
	assert.False(t, TriggerTaskRefresh(task.TaskID), "repeated hints within 2s are coalesced")

	// Expired tokens are rejected.
	expiredToken, err := NewTaskUpstreamCallbackToken()
	require.NoError(t, err)
	expiredHash := TaskUpstreamCallbackTokenHash(expiredToken)
	require.NoError(t, model.CreateTaskCallbackSubscription(&model.TaskCallbackSubscription{
		TaskID: "task_upstream_expired", UserID: userID, UpstreamTokenHash: &expiredHash, UpstreamTokenExpiresAt: common.GetTimestamp() - 1,
	}))
	assert.False(t, HandleTaskUpstreamCallback(expiredToken, ""))

	// Let the background refresh finish before the DB is truncated.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		taskRefreshState.Lock()
		running := taskRefreshState.running[task.TaskID]
		taskRefreshState.Unlock()
		if !running {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRecordTaskCallbackSubscriptionStoresUpstreamTokenHash(t *testing.T) {
	truncate(t)
	withTaskCallbackSetting(t, func(s *system_setting.TaskCallbackSetting) {
		s.Enabled = true
		s.UpstreamPublicBase = "https://kt.example"
		s.UpstreamInjectMode = system_setting.TaskCallbackUpstreamInjectAlways
	})
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	require.NoError(t, PrepareTaskUpstreamCallback(c, false))
	token := c.GetString(ContextKeyTaskUpstreamCallbackToken)
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{ChannelOtherSettings: dto.ChannelOtherSettings{UpstreamTaskCallbackEnabled: true}}}
	task := makeTask(641, 641, 100, 0, BillingSourceWallet, 0)
	task.TaskID = "task_token_only"

	subscription := RecordTaskCallbackSubscription(c, task, nil, info)
	require.NotNil(t, subscription)
	stored, found, err := model.GetTaskCallbackSubscriptionByID(subscription.ID)
	require.NoError(t, err)
	require.True(t, found)
	assert.Empty(t, stored.CallbackURL, "token-only subscriptions carry no customer URL")
	require.NotNil(t, stored.UpstreamTokenHash)
	assert.Equal(t, TaskUpstreamCallbackTokenHash(token), *stored.UpstreamTokenHash)
	assert.NotContains(t, *stored.UpstreamTokenHash, token, "only the hash is stored")
	assert.Greater(t, stored.UpstreamTokenExpiresAt, common.GetTimestamp()+3000)

	assert.Nil(t, RecordTaskCallbackSubscription(c, task, nil, &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{}}),
		"no customer URL and no channel opt-in means no subscription")
}
