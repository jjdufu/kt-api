package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

// Upstream callback fast path (off by default): kt-api hands upstream its
// own per-task callback URL; an upstream callback only triggers an
// immediate re-query through the polling path and its body is never
// trusted.
const (
	// TaskCallbackUpstreamPath is the public route prefix of the endpoint.
	TaskCallbackUpstreamPath = "/api/task-callback/upstream/"
	// ContextKeyTaskUpstreamCallbackToken carries the raw per-task token
	// from the submit middleware to the submit handler.
	ContextKeyTaskUpstreamCallbackToken = "task_upstream_callback_token"
	// taskUpstreamTokenTTLAfterDeadline keeps a token valid for an hour past
	// the task's execution deadline.
	taskUpstreamTokenTTLAfterDeadline = 3600
	// taskUpstreamRefreshMinInterval coalesces repeated upstream callbacks.
	taskUpstreamRefreshMinInterval = 2 * time.Second
	taskUpstreamRefreshTimeout     = 30 * time.Second
)

// NewTaskUpstreamCallbackToken returns a random base64url token (32 bytes).
func NewTaskUpstreamCallbackToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// TaskUpstreamCallbackTokenHash is the stored form of a token.
func TaskUpstreamCallbackTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// PrepareTaskUpstreamCallback decides at submit time whether this task may
// carry kt-api's upstream callback URL and, if so, stores the token and URL
// in the request context. Whether the URL is actually sent is decided per
// channel (upstream_task_callback_enabled) when the plugin builds the
// upstream request.
func PrepareTaskUpstreamCallback(c *gin.Context, customerCallback bool) error {
	setting := system_setting.GetTaskCallbackSetting()
	if !setting.Enabled || !setting.InjectUpstreamCallback(customerCallback) {
		return nil
	}
	token, err := NewTaskUpstreamCallbackToken()
	if err != nil {
		return err
	}
	c.Set(ContextKeyTaskUpstreamCallbackToken, token)
	common.SetContextKey(c, constant.ContextKeyTaskUpstreamCallbackURL, setting.UpstreamCallbackBase()+TaskCallbackUpstreamPath+token)
	return nil
}

// TaskUpstreamCallbackURLForChannel returns the URL a plugin should send
// upstream, or "" when the request has none or the channel does not enable
// the fast path.
func TaskUpstreamCallbackURLForChannel(c *gin.Context, info *relaycommon.RelayInfo) string {
	if c == nil || info == nil || info.ChannelMeta == nil || !info.ChannelOtherSettings.UpstreamTaskCallbackEnabled {
		return ""
	}
	return common.GetContextKeyString(c, constant.ContextKeyTaskUpstreamCallbackURL)
}

func taskUpstreamCallbackTokenFromContext(c *gin.Context, info *relaycommon.RelayInfo) string {
	if TaskUpstreamCallbackURLForChannel(c, info) == "" {
		return ""
	}
	return c.GetString(ContextKeyTaskUpstreamCallbackToken)
}

// HandleTaskUpstreamCallback resolves a token and schedules an immediate
// refresh of its task. It reports whether the token is known and valid.
func HandleTaskUpstreamCallback(token string, hintStatus string) bool {
	if token == "" || len(token) > 128 {
		return false
	}
	subscription, found, err := model.GetTaskCallbackSubscriptionByUpstreamTokenHash(TaskUpstreamCallbackTokenHash(token))
	if err != nil || !found {
		return false
	}
	if subscription.UpstreamTokenExpiresAt > 0 && subscription.UpstreamTokenExpiresAt < common.GetTimestamp() {
		return false
	}
	accepted := TriggerTaskRefresh(subscription.TaskID)
	common.SysLog(common.LogText("task_callback event=upstream_hint task_id=%s hint_status=%s accepted=%t", subscription.TaskID, hintStatus, accepted))
	return true
}

var taskRefreshState = struct {
	sync.Mutex
	running map[string]bool
	last    map[string]time.Time
}{running: map[string]bool{}, last: map[string]time.Time{}}

// TriggerTaskRefresh re-queries one task now through the polling path
// (FetchTask -> ParseTaskResult -> UpdateWithStatus CAS -> settlement), at
// most once per task every 2 seconds and never concurrently for the same
// task on this instance. It reports whether a refresh was started.
func TriggerTaskRefresh(taskID string) bool {
	now := time.Now()
	taskRefreshState.Lock()
	if taskRefreshState.running[taskID] || now.Sub(taskRefreshState.last[taskID]) < taskUpstreamRefreshMinInterval {
		taskRefreshState.Unlock()
		return false
	}
	taskRefreshState.running[taskID] = true
	taskRefreshState.last[taskID] = now
	for id, at := range taskRefreshState.last {
		if now.Sub(at) > time.Minute && !taskRefreshState.running[id] {
			delete(taskRefreshState.last, id)
		}
	}
	taskRefreshState.Unlock()
	go func() {
		defer func() {
			taskRefreshState.Lock()
			delete(taskRefreshState.running, taskID)
			taskRefreshState.Unlock()
		}()
		ctx, cancel := context.WithTimeout(context.Background(), taskUpstreamRefreshTimeout)
		defer cancel()
		if err := RefreshTaskNow(ctx, taskID); err != nil {
			common.SysError(common.LogText("task_callback event=refresh_failed task_id=%s error=%q", taskID, err.Error()))
		}
	}()
	return true
}

// RefreshTaskNow runs one polling update for an unfinished task.
func RefreshTaskNow(ctx context.Context, taskID string) error {
	task, exists, err := model.GetByOnlyTaskId(taskID)
	if err != nil || !exists || task == nil {
		return err
	}
	if task.Status == model.TaskStatusSuccess || task.Status == model.TaskStatusFailure {
		return nil
	}
	upstreamID := task.GetUpstreamTaskID()
	if upstreamID == "" {
		return nil
	}
	channel, err := model.CacheGetChannel(task.ChannelId)
	if err != nil {
		return err
	}
	adaptor := GetTaskAdaptorFunc(task.Platform)
	if adaptor == nil {
		return nil
	}
	info := &relaycommon.RelayInfo{}
	info.ChannelMeta = &relaycommon.ChannelMeta{
		ChannelType:    channel.Type,
		ChannelId:      channel.Id,
		ChannelBaseUrl: channel.GetBaseURL(),
		ChannelSetting: channel.GetSetting(),
	}
	info.ApiKey = channel.Key
	adaptor.Init(info)
	return updateVideoSingleTask(ctx, adaptor, channel, upstreamID, map[string]*model.Task{upstreamID: task})
}
