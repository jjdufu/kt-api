package service

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskCallbackSignatureVector(t *testing.T) {
	secret := TaskCallbackUserSecret("master-secret", 42)
	assert.Equal(t, "37b66edad4c2eb7fc494b9f2db025fed48c496693a9dd470ac9b2fd00838a105", secret)
	assert.Equal(t, secret, TaskCallbackUserSecret("master-secret", 42))
	assert.NotEqual(t, secret, TaskCallbackUserSecret("master-secret", 43))
	assert.NotEqual(t, secret, TaskCallbackUserSecret("other-master", 42))
	assert.Empty(t, TaskCallbackUserSecret("", 42))
	assert.Empty(t, TaskCallbackUserSecret("master-secret", 0))

	// Fixed vector: v1=hex(HMAC-SHA256("k", "1760000000.evt_1.{}")).
	signature := TaskCallbackSignature("k", 1760000000, "evt_1", []byte("{}"))
	assert.Equal(t, "v1=216affa7e5e4e57a43d51be7126cf056a210fe3e94808c738254c392d19f4d6a", signature)
	assert.True(t, VerifyTaskCallbackSignature("k", 1760000000, "evt_1", []byte("{}"), signature))
	assert.True(t, VerifyTaskCallbackSignature("k", 1760000000, "evt_1", []byte("{}"), "v1=deadbeef, "+signature))
	assert.False(t, VerifyTaskCallbackSignature("k", 1760000001, "evt_1", []byte("{}"), signature))
	assert.False(t, VerifyTaskCallbackSignature("k", 1760000000, "evt_2", []byte("{}"), signature))
	assert.False(t, VerifyTaskCallbackSignature("k", 1760000000, "evt_1", []byte("{ }"), signature))
}

func TestTaskCallbackDeliverySignsWhenEnabled(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusOK))
	subscription := h.subscribe("task_d_signed", "/cb")
	payload := `{"id":"task_d_signed","status":"succeeded"}`
	event := h.enqueue(subscription, "succeeded", payload)

	h.run()
	require.Len(t, h.receiver.requests, 1)
	assert.Empty(t, h.receiver.requests[0].Header.Get(TaskCallbackHeaderSignature), "signing is off by default")
	assert.Empty(t, h.receiver.requests[0].Header.Get(TaskCallbackHeaderTimestamp))

	h.setting.SigningEnabled = true
	h.setting.SigningSecret = "master-secret"
	second := h.enqueue(h.subscribe("task_d_signed_2", "/cb"), "succeeded", payload)
	h.run()
	require.Len(t, h.receiver.requests, 2)
	request := h.receiver.requests[1]
	timestamp := request.Header.Get(TaskCallbackHeaderTimestamp)
	assert.Equal(t, h.clock.Now().Unix(), parseInt64(t, timestamp))
	secret := TaskCallbackUserSecret("master-secret", 31)
	assert.True(t, VerifyTaskCallbackSignature(secret, parseInt64(t, timestamp), second.EventID, []byte(payload), request.Header.Get(TaskCallbackHeaderSignature)))
	assert.False(t, VerifyTaskCallbackSignature(secret, parseInt64(t, timestamp), event.EventID, []byte(payload), request.Header.Get(TaskCallbackHeaderSignature)))
	assert.Equal(t, model.TaskCallbackStateSucceeded, h.reload(second).State)
}

func parseInt64(t *testing.T, value string) int64 {
	t.Helper()
	parsed, err := strconv.ParseInt(value, 10, 64)
	require.NoError(t, err)
	return parsed
}
