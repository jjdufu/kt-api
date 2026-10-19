package service

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTaskCallbackClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeTaskCallbackClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeTaskCallbackClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

type callbackReceiver struct {
	server   *httptest.Server
	mu       sync.Mutex
	requests []*http.Request
	bodies   []string
	respond  func(n int, w http.ResponseWriter)
	hits     atomic.Int32
}

func newCallbackReceiver(t *testing.T, respond func(n int, w http.ResponseWriter)) *callbackReceiver {
	t.Helper()
	receiver := &callbackReceiver{respond: respond}
	receiver.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(receiver.hits.Add(1))
		body, _ := io.ReadAll(r.Body)
		receiver.mu.Lock()
		receiver.requests = append(receiver.requests, r.Clone(context.Background()))
		receiver.bodies = append(receiver.bodies, string(body))
		receiver.mu.Unlock()
		receiver.respond(n, w)
	}))
	t.Cleanup(receiver.server.Close)
	return receiver
}

func (r *callbackReceiver) port(t *testing.T) string {
	_, port, err := net.SplitHostPort(strings.TrimPrefix(r.server.URL, "http://"))
	require.NoError(t, err)
	return port
}

type dispatcherHarness struct {
	t          *testing.T
	clock      *fakeTaskCallbackClock
	setting    *system_setting.TaskCallbackSetting
	dispatcher *TaskCallbackDispatcher
	receiver   *callbackReceiver
}

func newDispatcherHarness(t *testing.T, respond func(n int, w http.ResponseWriter)) *dispatcherHarness {
	t.Helper()
	truncate(t)
	receiver := newCallbackReceiver(t, respond)
	setting := &system_setting.TaskCallbackSetting{
		Enabled:              true,
		ExtraHeadersEnabled:  true,
		AllowPrivateIp:       true,
		AllowedPorts:         []string{receiver.port(t)},
		TimeoutSeconds:       5,
		RetryIntervalSeconds: 5,
		TerminalMaxAttempts:  4,
		Workers:              4,
		PerHostConcurrency:   4,
	}
	clock := &fakeTaskCallbackClock{now: time.Unix(time.Now().Unix(), 0)}
	getSetting := func() *system_setting.TaskCallbackSetting { return setting }
	dispatcher := newTaskCallbackDispatcher(newTaskCallbackHTTPClient(nil, nil, getSetting), clock.Now, getSetting)
	return &dispatcherHarness{t: t, clock: clock, setting: setting, dispatcher: dispatcher, receiver: receiver}
}

func (h *dispatcherHarness) newDispatcher() *TaskCallbackDispatcher {
	getSetting := func() *system_setting.TaskCallbackSetting { return h.setting }
	return newTaskCallbackDispatcher(newTaskCallbackHTTPClient(nil, nil, getSetting), h.clock.Now, getSetting)
}

func (h *dispatcherHarness) subscribe(taskID, path string) *model.TaskCallbackSubscription {
	h.t.Helper()
	url := h.receiver.server.URL + path
	subscription := &model.TaskCallbackSubscription{TaskID: taskID, UserID: 31, CallbackURL: url, CallbackURLHash: TaskCallbackURLHash(url), CallbackHost: "127.0.0.1"}
	require.NoError(h.t, model.CreateTaskCallbackSubscription(subscription))
	return subscription
}

func (h *dispatcherHarness) enqueue(subscription *model.TaskCallbackSubscription, status, payload string) *model.TaskCallbackEvent {
	h.t.Helper()
	event := &model.TaskCallbackEvent{
		EventID:       newTaskCallbackEventID(),
		Status:        status,
		Payload:       payload,
		MaxAttempts:   taskCallbackMaxAttemptsFor(h.setting, status),
		NextAttemptAt: h.clock.Now().Unix(),
	}
	inserted, err := model.InsertTaskCallbackEventTx(model.DB, subscription, event)
	require.NoError(h.t, err)
	require.True(h.t, inserted)
	return event
}

func taskCallbackMaxAttemptsFor(setting *system_setting.TaskCallbackSetting, status string) int {
	if taskCallbackDocumentedRetries(status) {
		return setting.TerminalAttempts()
	}
	return 1
}

func (h *dispatcherHarness) run() int {
	h.t.Helper()
	started := h.dispatcher.DispatchOnce(context.Background())
	h.dispatcher.Wait()
	return started
}

func (h *dispatcherHarness) reload(event *model.TaskCallbackEvent) *model.TaskCallbackEvent {
	h.t.Helper()
	reloaded, err := model.GetTaskCallbackEvent(event.ID)
	require.NoError(h.t, err)
	return reloaded
}

func respondStatus(status int) func(int, http.ResponseWriter) {
	return func(_ int, w http.ResponseWriter) { w.WriteHeader(status) }
}

func TestTaskCallbackDeliverySucceedsWithHeadersAndExactBody(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusNoContent))
	subscription := h.subscribe("task_d_ok", "/hook?token=secret")
	payload := `{"id":"task_d_ok","status":"succeeded","content":{"video_url":"https://cdn.example/v.mp4?a=1\u0026b=2"}}`
	event := h.enqueue(subscription, "succeeded", payload)

	assert.Equal(t, 1, h.run())
	stored := h.reload(event)
	assert.Equal(t, model.TaskCallbackStateSucceeded, stored.State)
	assert.Equal(t, 1, stored.Attempts)
	assert.Equal(t, http.StatusNoContent, stored.LastHTTPStatus)
	assert.Equal(t, h.clock.Now().Unix(), stored.DeliveredAt)
	assert.Empty(t, stored.LastError)

	require.Len(t, h.receiver.requests, 1)
	request := h.receiver.requests[0]
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/hook", request.URL.Path)
	assert.Equal(t, "secret", request.URL.Query().Get("token"))
	assert.Equal(t, payload, h.receiver.bodies[0])
	assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	assert.True(t, strings.HasPrefix(request.Header.Get("User-Agent"), "kt-api-callback/"))
	assert.Equal(t, event.EventID, request.Header.Get("X-KtApi-Event-Id"))
	assert.Equal(t, "1", request.Header.Get("X-KtApi-Delivery-Attempt"))
	assert.Equal(t, "task_d_ok", request.Header.Get("X-KtApi-Task-Id"))

	// Delivered events are not sent again.
	assert.Equal(t, 0, h.run())
	assert.Equal(t, int32(1), h.receiver.hits.Load())
}

func TestTaskCallbackTerminalRetriesThreeTimesFiveSecondsApart(t *testing.T) {
	h := newDispatcherHarness(t, func(n int, w http.ResponseWriter) {
		if n < 4 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	event := h.enqueue(h.subscribe("task_d_retry", "/cb"), "failed", `{"id":"task_d_retry","status":"failed"}`)

	for attempt := 1; attempt <= 3; attempt++ {
		assert.Equal(t, 1, h.run())
		stored := h.reload(event)
		assert.Equal(t, model.TaskCallbackStatePending, stored.State)
		assert.Equal(t, attempt, stored.Attempts)
		assert.Equal(t, h.clock.Now().Unix()+5, stored.NextAttemptAt)
		assert.Equal(t, 500, stored.LastHTTPStatus)
		assert.Equal(t, "http_5xx", stored.LastError)
		// Not due before the 5s interval.
		h.clock.Advance(4 * time.Second)
		assert.Equal(t, 0, h.run())
		h.clock.Advance(time.Second)
	}
	assert.Equal(t, 1, h.run())
	stored := h.reload(event)
	assert.Equal(t, model.TaskCallbackStateSucceeded, stored.State)
	assert.Equal(t, 4, stored.Attempts)
	var attempts []string
	for _, request := range h.receiver.requests {
		attempts = append(attempts, request.Header.Get("X-KtApi-Delivery-Attempt"))
		assert.Equal(t, event.EventID, request.Header.Get("X-KtApi-Event-Id"))
	}
	assert.Equal(t, []string{"1", "2", "3", "4"}, attempts)
}

func TestTaskCallbackGivesUpAfterFourTerminalAttemptsAndOneOtherwise(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusServiceUnavailable))
	terminal := h.enqueue(h.subscribe("task_d_fail", "/cb"), "succeeded", `{"status":"succeeded"}`)
	running := h.enqueue(h.subscribe("task_d_running", "/cb"), "running", `{"status":"running"}`)
	expired := h.enqueue(h.subscribe("task_d_expired", "/cb"), "expired", `{"status":"expired"}`)

	for i := 0; i < 6; i++ {
		h.run()
		h.clock.Advance(5 * time.Second)
	}
	assert.Equal(t, model.TaskCallbackStateFailed, h.reload(terminal).State)
	assert.Equal(t, 4, h.reload(terminal).Attempts)
	assert.Equal(t, model.TaskCallbackStateFailed, h.reload(running).State)
	assert.Equal(t, 1, h.reload(running).Attempts)
	assert.Equal(t, model.TaskCallbackStateFailed, h.reload(expired).State)
	assert.Equal(t, 1, h.reload(expired).Attempts)
	assert.Equal(t, int32(6), h.receiver.hits.Load())
}

func TestTaskCallbackTimeoutRedirectAndSSRFClassification(t *testing.T) {
	t.Run("slow receiver times out and is retried", func(t *testing.T) {
		h := newDispatcherHarness(t, func(_ int, w http.ResponseWriter) {
			time.Sleep(1500 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		})
		h.setting.TimeoutSeconds = 1
		event := h.enqueue(h.subscribe("task_d_slow", "/cb"), "succeeded", `{}`)
		h.run()
		stored := h.reload(event)
		assert.Equal(t, model.TaskCallbackStatePending, stored.State)
		assert.True(t, strings.HasPrefix(stored.LastError, "timeout"), stored.LastError)
		assert.NotContains(t, stored.LastError, h.receiver.server.URL)
	})
	t.Run("redirect is not followed", func(t *testing.T) {
		h := newDispatcherHarness(t, func(_ int, w http.ResponseWriter) {
			w.Header().Set("Location", "http://169.254.169.254/")
			w.WriteHeader(http.StatusFound)
		})
		event := h.enqueue(h.subscribe("task_d_redirect", "/cb"), "running", `{}`)
		h.run()
		stored := h.reload(event)
		assert.Equal(t, model.TaskCallbackStateFailed, stored.State)
		assert.Equal(t, "http_3xx", stored.LastError)
		assert.Equal(t, int32(1), h.receiver.hits.Load())
	})
	t.Run("private target blocked at delivery is not retried", func(t *testing.T) {
		h := newDispatcherHarness(t, respondStatus(http.StatusOK))
		h.setting.AllowPrivateIp = false
		event := h.enqueue(h.subscribe("task_d_ssrf", "/cb"), "succeeded", `{}`)
		h.run()
		stored := h.reload(event)
		assert.Equal(t, model.TaskCallbackStateFailed, stored.State)
		assert.Equal(t, 1, stored.Attempts)
		assert.True(t, strings.HasPrefix(stored.LastError, "ssrf_blocked"), stored.LastError)
		assert.Equal(t, int32(0), h.receiver.hits.Load())
	})
	t.Run("extra headers can be disabled", func(t *testing.T) {
		h := newDispatcherHarness(t, respondStatus(http.StatusOK))
		h.setting.ExtraHeadersEnabled = false
		h.enqueue(h.subscribe("task_d_noheaders", "/cb"), "queued", `{}`)
		h.run()
		require.Len(t, h.receiver.requests, 1)
		assert.Empty(t, h.receiver.requests[0].Header.Get("X-KtApi-Event-Id"))
		assert.Empty(t, h.receiver.requests[0].Header.Get("X-KtApi-Delivery-Attempt"))
	})
}

func TestTaskCallbackPerTaskOrderingAndSupersede(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusOK))
	subscription := h.subscribe("task_d_order", "/cb")
	queued := h.enqueue(subscription, "queued", `{"status":"queued"}`)
	running := h.enqueue(subscription, "running", `{"status":"running"}`)
	succeeded := h.enqueue(subscription, "succeeded", `{"status":"succeeded"}`)

	// A newer event supersedes undelivered queued/running ones at enqueue;
	// only the newest is sent.
	assert.Equal(t, model.TaskCallbackStateSuperseded, h.reload(queued).State)
	assert.Equal(t, model.TaskCallbackStateSuperseded, h.reload(running).State)
	h.run()
	assert.Equal(t, model.TaskCallbackStateSucceeded, h.reload(succeeded).State)
	assert.Equal(t, []string{`{"status":"succeeded"}`}, h.receiver.bodies)

	// An earlier event in flight elsewhere blocks later ones of the task.
	other := h.subscribe("task_d_wait", "/cb")
	first := h.enqueue(other, "running", `{"status":"running"}`)
	claimed, err := model.ClaimTaskCallbackEvent(first.ID, "other-instance", h.clock.Now().Unix(), h.clock.Now().Unix()+30)
	require.NoError(t, err)
	require.True(t, claimed)
	second := h.enqueue(other, "succeeded", `{"status":"succeeded"}`)
	assert.Equal(t, model.TaskCallbackStateDelivering, h.reload(first).State)
	assert.Equal(t, 0, h.run())
	assert.Equal(t, model.TaskCallbackStatePending, h.reload(second).State)
	_, err = model.FinishTaskCallbackAttempt(first.ID, "other-instance", model.TaskCallbackAttemptResult{State: model.TaskCallbackStateSucceeded})
	require.NoError(t, err)
	assert.Equal(t, 1, h.run())
	assert.Equal(t, model.TaskCallbackStateSucceeded, h.reload(second).State)
}

func TestTaskCallbackConcurrentDispatchersDeliverOnce(t *testing.T) {
	h := newDispatcherHarness(t, func(_ int, w http.ResponseWriter) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	var events []*model.TaskCallbackEvent
	for i := 0; i < 12; i++ {
		taskID := "task_d_multi_" + string(rune('a'+i))
		events = append(events, h.enqueue(h.subscribe(taskID, "/cb"), "succeeded", `{"id":"`+taskID+`"}`))
	}
	dispatchers := []*TaskCallbackDispatcher{h.dispatcher, h.newDispatcher(), h.newDispatcher()}
	for round := 0; round < 4; round++ {
		var wg sync.WaitGroup
		for _, dispatcher := range dispatchers {
			wg.Add(1)
			go func(d *TaskCallbackDispatcher) {
				defer wg.Done()
				d.DispatchOnce(context.Background())
				d.Wait()
			}(dispatcher)
		}
		wg.Wait()
	}
	assert.Equal(t, int32(12), h.receiver.hits.Load())
	seen := map[string]int{}
	for _, body := range h.receiver.bodies {
		seen[body]++
	}
	assert.Len(t, seen, 12)
	for _, event := range events {
		stored := h.reload(event)
		assert.Equal(t, model.TaskCallbackStateSucceeded, stored.State)
		assert.Equal(t, 1, stored.Attempts)
	}
}

func TestTaskCallbackReclaimAfterCrashKeepsEventID(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusOK))
	event := h.enqueue(h.subscribe("task_d_crash", "/cb"), "succeeded", `{"status":"succeeded"}`)
	// A previous instance claimed the event and died mid-attempt.
	claimed, err := model.ClaimTaskCallbackEvent(event.ID, "dead-instance", h.clock.Now().Unix(), h.clock.Now().Unix()+30)
	require.NoError(t, err)
	require.True(t, claimed)
	assert.Equal(t, 0, h.run())

	h.clock.Advance(31 * time.Second)
	h.dispatcher.reclaim()
	assert.Equal(t, 1, h.run())
	stored := h.reload(event)
	assert.Equal(t, model.TaskCallbackStateSucceeded, stored.State)
	assert.Equal(t, 2, stored.Attempts)
	require.Len(t, h.receiver.requests, 1)
	assert.Equal(t, event.EventID, h.receiver.requests[0].Header.Get("X-KtApi-Event-Id"))
	assert.Equal(t, "2", h.receiver.requests[0].Header.Get("X-KtApi-Delivery-Attempt"))
}

func TestTaskCallbackExtendedRetryIsOptIn(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusBadGateway))
	h.setting.ExtendedRetry = true
	event := h.enqueue(h.subscribe("task_d_extended", "/cb"), "succeeded", `{}`)
	for i := 0; i < 4; i++ {
		h.run()
		h.clock.Advance(5 * time.Second)
	}
	stored := h.reload(event)
	assert.Equal(t, model.TaskCallbackStatePending, stored.State)
	assert.Equal(t, 4, stored.Attempts)
	assert.Equal(t, 5, stored.MaxAttempts)
	assert.Equal(t, h.clock.Now().Unix()-5+30, stored.NextAttemptAt)

	// Non-terminal statuses never get the extended retry.
	running := h.enqueue(h.subscribe("task_d_extended_running", "/cb"), "running", `{}`)
	h.run()
	assert.Equal(t, model.TaskCallbackStateFailed, h.reload(running).State)
}

func TestTaskCallbackDisabledDispatcherIsIdle(t *testing.T) {
	h := newDispatcherHarness(t, respondStatus(http.StatusOK))
	h.setting.Enabled = false
	event := h.enqueue(h.subscribe("task_d_disabled", "/cb"), "succeeded", `{}`)
	assert.Equal(t, 0, h.run())
	assert.Equal(t, model.TaskCallbackStatePending, h.reload(event).State)
}
