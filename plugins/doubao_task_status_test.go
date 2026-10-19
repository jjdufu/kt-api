package plugins_test

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoubaoParseTaskResultSupplierFailedEnvelope(t *testing.T) {
	_, plugin := newDoubaoPlugin(t)
	queryContext := map[string]any{"upstreamModel": "doubao-seedance-2-0-260128", "model": "doubao-seedance-2-0-260128", "action": "text_to_video"}
	parse := func(t *testing.T, body map[string]any, status int) map[string]any {
		t.Helper()
		value, err := plugin.Engine.Call(t.Context(), "parseTaskResult", queryContext, body, map[string]any{"status": status, "headers": map[string]any{}})
		require.NoError(t, err)
		encoded, err := common.Marshal(value)
		require.NoError(t, err)
		decoded := map[string]any{}
		require.NoError(t, common.Unmarshal(encoded, &decoded))
		return decoded
	}
	envelope := map[string]any{"error": map[string]any{"code": "task_not_exist", "message": "task_not_exist", "type": "yisu_api_error", "request_id": "req"}}

	t.Run("HTTP 400 bare error envelope is a failure", func(t *testing.T) {
		result := parse(t, envelope, 400)
		assert.Equal(t, "FAILURE", result["status"])
		assert.Equal(t, "100%", result["progress"])
		assert.Equal(t, "task_not_exist", result["reason"])
	})
	t.Run("envelope with only a code uses the code", func(t *testing.T) {
		result := parse(t, map[string]any{"error": map[string]any{"code": "OutputVideoSensitiveContentDetected"}}, 400)
		assert.Equal(t, "FAILURE", result["status"])
		assert.Equal(t, "OutputVideoSensitiveContentDetected", result["reason"])
	})
	t.Run("the same envelope on another status stays unrecognized", func(t *testing.T) {
		result := parse(t, envelope, 500)
		assert.Equal(t, "UNKNOWN", result["status"])
	})
	t.Run("an Ark failed body keeps its own status mapping", func(t *testing.T) {
		result := parse(t, map[string]any{"id": "cgt", "status": "failed", "error": map[string]any{"code": "x", "message": "boom"}}, 200)
		assert.Equal(t, "FAILURE", result["status"])
		assert.Equal(t, "boom", result["reason"])
	})
	t.Run("expired and cancelled stay failures", func(t *testing.T) {
		assert.Equal(t, "FAILURE", parse(t, map[string]any{"status": "expired"}, 200)["status"])
		assert.Equal(t, "FAILURE", parse(t, map[string]any{"status": "cancelled"}, 200)["status"])
	})
}

func TestDoubaoVideoSubmitCallbackURL(t *testing.T) {
	_, plugin := newDoubaoPlugin(t)
	model := "doubao-seedance-2-0-260128"
	submit := func(t *testing.T, extra map[string]any) map[string]any {
		t.Helper()
		driver := map[string]any{
			"upstreamModel": model, "model": model, "action": "text_to_video", "baseUrl": "https://ark.example", "apiKey": "k",
			"requestBody": map[string]any{
				"model":  model,
				"prompt": "a cat",
				"metadata": map[string]any{
					"content":      []any{map[string]any{"type": "text", "text": "a cat"}},
					"callback_url": "https://customer.example/hook?secret=1",
				},
			},
		}
		for key, value := range extra {
			driver[key] = value
		}
		value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", driver)
		require.NoError(t, err)
		return alibabaObject(t, value)["body"].(map[string]any)
	}

	t.Run("customer callback_url never reaches upstream", func(t *testing.T) {
		body := submit(t, nil)
		assert.NotContains(t, body, "callback_url")
	})
	t.Run("fast path URL is sent when the host provides one", func(t *testing.T) {
		body := submit(t, map[string]any{"upstreamCallbackUrl": "https://kt.example/api/task-callback/upstream/tok"})
		assert.Equal(t, "https://kt.example/api/task-callback/upstream/tok", body["callback_url"])
	})
}
