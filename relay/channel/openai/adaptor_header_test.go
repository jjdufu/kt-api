package openai

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// OpenRouter requests must not carry upstream-project attribution headers;
// operators who want them can still set them through header overrides.
func TestSetupRequestHeaderOpenRouterAddsNoAttributionHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType: constant.ChannelTypeOpenRouter,
			ApiKey:      "sk-test",
		},
	}
	header := http.Header{}

	require.NoError(t, (&Adaptor{}).SetupRequestHeader(ctx, &header, info))
	require.Equal(t, "Bearer sk-test", header.Get("Authorization"))
	require.Empty(t, header.Get("HTTP-Referer"))
	require.Empty(t, header.Get("X-OpenRouter-Title"))
}
