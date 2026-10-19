package system_setting

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/setting/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskCallbackSettingDefaults(t *testing.T) {
	setting := GetTaskCallbackSetting()
	assert.True(t, setting.Enabled)
	assert.True(t, setting.ExtraHeadersEnabled)
	assert.False(t, setting.HTTPSOnly)
	assert.False(t, setting.AllowPrivateIp)
	assert.False(t, setting.ExtendedRetry)
	assert.False(t, setting.SigningEnabled)
	assert.Equal(t, 5*time.Second, setting.Timeout())
	assert.Equal(t, 5*time.Second, setting.RetryInterval())
	assert.Equal(t, 4, setting.TerminalAttempts())
	assert.Equal(t, 2048, setting.URLLengthLimit())
	assert.Empty(t, setting.UpstreamCallbackBase(), "upstream fast path is off by default")
	assert.False(t, setting.InjectUpstreamCallback(true))
}

func TestTaskCallbackSettingClamps(t *testing.T) {
	setting := TaskCallbackSetting{TimeoutSeconds: 999, RetryIntervalSeconds: -1, TerminalMaxAttempts: 50, Workers: 0, PerHostConcurrency: 1000, RetentionDays: 400, MaxURLLength: 10}
	assert.Equal(t, 30*time.Second, setting.Timeout())
	assert.Equal(t, 5*time.Second, setting.RetryInterval())
	assert.Equal(t, 10, setting.TerminalAttempts())
	assert.Equal(t, 16, setting.WorkerCount())
	assert.Equal(t, 64, setting.HostConcurrency())
	assert.Equal(t, 90*24*time.Hour, setting.Retention())
	assert.Equal(t, 64, setting.URLLengthLimit())
}

func TestTaskCallbackSettingUpstreamInjection(t *testing.T) {
	setting := TaskCallbackSetting{UpstreamPublicBase: " https://api.example.com/ ", UpstreamInjectMode: TaskCallbackUpstreamInjectOnlyWhenCustomerCallback}
	assert.Equal(t, "https://api.example.com", setting.UpstreamCallbackBase())
	assert.True(t, setting.InjectUpstreamCallback(true))
	assert.False(t, setting.InjectUpstreamCallback(false))
	setting.UpstreamInjectMode = TaskCallbackUpstreamInjectAlways
	assert.True(t, setting.InjectUpstreamCallback(false))
	setting.UpstreamInjectMode = TaskCallbackUpstreamInjectNever
	assert.False(t, setting.InjectUpstreamCallback(true))
	assert.Empty(t, setting.UpstreamCallbackBase())
}

func TestTaskCallbackSettingLoadsFromOptions(t *testing.T) {
	previous := *GetTaskCallbackSetting()
	t.Cleanup(func() { *GetTaskCallbackSetting() = previous })
	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"task_callback_setting.https_only":     "true",
		"task_callback_setting.domain_list":    `["*.example.com"]`,
		"task_callback_setting.workers":        "8",
		"task_callback_setting.signing_secret": "s3cret",
	}))
	setting := GetTaskCallbackSetting()
	assert.True(t, setting.HTTPSOnly)
	assert.Equal(t, []string{"*.example.com"}, setting.DomainList)
	assert.Equal(t, 8, setting.WorkerCount())
	assert.Equal(t, "s3cret", setting.SigningSecret)
}
