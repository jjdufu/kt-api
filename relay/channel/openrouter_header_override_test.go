package channel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// The OpenRouter adaptor no longer injects attribution headers, so an
// administrator header override is the only source and must pass through.
func TestProcessHeaderOverride_OpenRouterAttributionFromAdminOverride(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType: constant.ChannelTypeOpenRouter,
			HeadersOverride: map[string]any{
				"HTTP-Referer":       "https://relay.example.com",
				"X-OpenRouter-Title": "Example Relay",
			},
		},
	}

	headers, err := processHeaderOverride(info, ctx)
	require.NoError(t, err)
	require.Equal(t, "https://relay.example.com", headers["http-referer"])
	require.Equal(t, "Example Relay", headers["x-openrouter-title"])
}
