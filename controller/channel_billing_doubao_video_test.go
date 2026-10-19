package controller

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const doubaoVideoTestKey = "sk-seedance-secret"

func newDoubaoVideoBalanceChannel(t *testing.T, serverURL, settings string) *model.Channel {
	t.Helper()
	baseURL := serverURL + "/"
	channel := &model.Channel{
		Type:          constant.ChannelTypeDoubaoVideo,
		BaseURL:       &baseURL,
		Key:           doubaoVideoTestKey,
		Name:          "seedance-upstream",
		Status:        common.ChannelStatusEnabled,
		OtherSettings: settings,
		Balance:       7,
	}
	require.NoError(t, model.DB.Create(channel).Error)
	return channel
}

func TestUpdateChannelDoubaoVideoBalanceFromNewAPIUpstream(t *testing.T) {
	for _, tc := range []struct {
		name     string
		settings string
		usage    string
		want     float64
	}{
		{name: "default 500000 quota per unit", usage: `{"code":true,"message":"ok","data":{"object":"token_usage","name":"k","total_granted":5000000,"total_used":1000000,"total_available":4000000,"unlimited_quota":false,"model_limits_enabled":false,"expires_at":0}}`, want: 8},
		{name: "configured quota per unit", settings: `{"upstream_quota_per_unit":1000000}`, usage: `{"code":true,"data":{"total_available":4000000,"unlimited_quota":false}}`, want: 4},
		{name: "non-positive quota per unit falls back to the default", settings: `{"upstream_quota_per_unit":-3}`, usage: `{"code":true,"data":{"total_available":250000,"unlimited_quota":false}}`, want: 0.5},
		{name: "unlimited key reports its owner's quota", usage: `{"code":true,"data":{"total_available":0,"unlimited_quota":true,"user_quota":{"total_available":1500000}}}`, want: 3},
		{name: "owner quota below the key quota", usage: `{"code":true,"data":{"total_available":9000000,"unlimited_quota":false,"user_quota":{"total_available":500000}}}`, want: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupSub2APIBalanceTest(t)
			var statusCalls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/usage/token/":
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "Bearer "+doubaoVideoTestKey, r.Header.Get("Authorization"))
					_, _ = w.Write([]byte(tc.usage))
				case "/api/status":
					statusCalls.Add(1)
					http.NotFound(w, r)
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			channel := newDoubaoVideoBalanceChannel(t, server.URL, tc.settings)

			result, err := updateChannelBalance(channel)
			require.NoError(t, err)
			assert.InDelta(t, tc.want, result.Balance, 1e-9)
			assert.Empty(t, result.RawResponse)
			assert.Zero(t, statusCalls.Load(), "the conversion comes from the channel setting, not the upstream status")

			stored, err := model.GetChannelById(channel.Id, true)
			require.NoError(t, err)
			assert.InDelta(t, tc.want, stored.Balance, 1e-9)
			assert.NotZero(t, stored.BalanceUpdatedTime)
		})
	}
}

func TestUpdateChannelDoubaoVideoBalanceErrorsWithoutTouchingTheChannel(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler http.HandlerFunc
		wantErr string
	}{
		{name: "official Ark has no usage API", handler: func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }, wantErr: "status code: 404"},
		{name: "HTML page", handler: func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("<html>console</html>")) }, wantErr: "did not return New API usage JSON"},
		{name: "rejected key", handler: func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"code":false,"message":"invalid token"}`))
		}, wantErr: "invalid token"},
		{name: "missing amount", handler: func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"code":true,"data":{}}`)) }, wantErr: "missing data.total_available"},
		{name: "missing data", handler: func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"code":true}`)) }, wantErr: "missing data.total_available"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setupSub2APIBalanceTest(t)
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			channel := newDoubaoVideoBalanceChannel(t, server.URL, "")

			_, err := updateChannelBalance(channel)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "New API compatible upstream")
			assert.Contains(t, err.Error(), tc.wantErr)
			assert.NotContains(t, err.Error(), doubaoVideoTestKey)
			var message *common.Message
			require.ErrorAs(t, err, &message, "the console translates the error")
			assert.Equal(t, doubaoVideoBalanceUnsupportedMessage, message.Key)

			stored, err := model.GetChannelById(channel.Id, true)
			require.NoError(t, err)
			assert.InDelta(t, 7, stored.Balance, 1e-9, "a failed query keeps the last balance")
			assert.Zero(t, stored.BalanceUpdatedTime)
		})
	}

	t.Run("unlimited key without owner quota", func(t *testing.T) {
		setupSub2APIBalanceTest(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"code":true,"data":{"total_available":0,"unlimited_quota":true}}`))
		}))
		defer server.Close()
		channel := newDoubaoVideoBalanceChannel(t, server.URL, "")

		_, err := updateChannelBalance(channel)
		require.ErrorContains(t, err, "Allow querying user quota")
	})
}

func TestUpdateAllChannelsBalanceKeepsDoubaoVideoChannelOnQueryError(t *testing.T) {
	setupSub2APIBalanceTest(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }))
	defer server.Close()
	failing := newDoubaoVideoBalanceChannel(t, server.URL, "")

	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"code":true,"data":{"total_available":1000000,"unlimited_quota":false}}`))
	}))
	defer okServer.Close()
	working := newDoubaoVideoBalanceChannel(t, okServer.URL, "")

	require.NoError(t, updateAllChannelsBalance())

	stored, err := model.GetChannelById(failing.Id, true)
	require.NoError(t, err)
	assert.Equal(t, common.ChannelStatusEnabled, stored.Status, "a non-New-API upstream is not disabled")
	assert.InDelta(t, 7, stored.Balance, 1e-9)

	stored, err = model.GetChannelById(working.Id, true)
	require.NoError(t, err)
	assert.Equal(t, common.ChannelStatusEnabled, stored.Status)
	assert.InDelta(t, 2, stored.Balance, 1e-9, "the periodic update stores the converted balance")
}
