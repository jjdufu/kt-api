package jsplugin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitContextUpstreamCallbackURLRequiresChannelOptIn(t *testing.T) {
	setting := system_setting.GetTaskCallbackSetting()
	previous := *setting
	t.Cleanup(func() { *setting = previous })
	setting.Enabled = true
	setting.UpstreamPublicBase = "https://kt.example"
	setting.UpstreamInjectMode = system_setting.TaskCallbackUpstreamInjectOnlyWhenCustomerCallback

	plugin, err := pluginruntime.NewRegistry().Register(mockPlugin, pluginruntime.Options{})
	require.NoError(t, err)
	adaptor := New(plugin)
	for _, tc := range []struct {
		name    string
		enabled bool
	}{{"channel opted in", true}, {"channel not opted in", false}} {
		t.Run(tc.name, func(t *testing.T) {
			info := &relaycommon.RelayInfo{
				ChannelMeta: &relaycommon.ChannelMeta{
					ChannelBaseUrl:       "https://provider.example",
					ApiKey:               "secret",
					ChannelOtherSettings: dto.ChannelOtherSettings{UpstreamTaskCallbackEnabled: tc.enabled},
				},
				TaskRelayInfo: &relaycommon.TaskRelayInfo{},
			}
			adaptor.Init(info)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v3/contents/generations/tasks", nil)
			require.NoError(t, service.PrepareTaskUpstreamCallback(c, true))

			ctx := adaptor.submitContext(c, info)
			value, ok := ctx["upstreamCallbackUrl"]
			if !tc.enabled {
				assert.False(t, ok)
				return
			}
			token := c.GetString(service.ContextKeyTaskUpstreamCallbackToken)
			require.NotEmpty(t, token)
			assert.Equal(t, "https://kt.example/api/task-callback/upstream/"+token, value)
		})
	}
}
