package service

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

const (
	taskCallbackPollInterval     = time.Second
	taskCallbackReclaimInterval  = time.Minute
	taskCallbackLease            = 30 * time.Second
	taskCallbackMaxResponseBytes = 64 << 10
	taskCallbackShutdownWait     = 10 * time.Second
	// taskCallbackExtendedRetryWindow bounds the optional extended retry.
	taskCallbackExtendedRetryWindow = 24 * time.Hour
)

// taskCallbackExtendedBackoff is the optional back-off after the documented
// retries are exhausted (off by default; customers may then receive
// callbacks outside the documented schedule). The last step repeats until
// the 24h window closes.
var taskCallbackExtendedBackoff = []time.Duration{
	30 * time.Second,
	2 * time.Minute,
	10 * time.Minute,
	30 * time.Minute,
	time.Hour,
}

// Delivery failure classes, written to last_error and the attempt log.
const (
	taskCallbackResultSucceeded   = "succeeded"
	taskCallbackResultTimeout     = "timeout"
	taskCallbackResultDNS         = "dns"
	taskCallbackResultSSRFBlocked = "ssrf_blocked"
	taskCallbackResultConnect     = "connect"
	taskCallbackResultTLS         = "tls"
	taskCallbackResultInvalid     = "invalid_request"
)

// TaskCallbackDispatcher delivers pending outbox events. Every instance runs
// one; conditional-UPDATE claims with leases make concurrent dispatchers
// safe, and per-task ordering allows one in-flight delivery per task.
type TaskCallbackDispatcher struct {
	owner      string
	client     *http.Client
	now        func() time.Time
	getSetting func() *system_setting.TaskCallbackSetting
	wake       <-chan struct{}

	mu       sync.Mutex
	inflight map[int64]struct{}
	tasks    map[string]struct{}
	hosts    map[string]int
	wg       sync.WaitGroup
}

// NewTaskCallbackDispatcher builds a dispatcher with the protected delivery
// client.
func NewTaskCallbackDispatcher() *TaskCallbackDispatcher {
	return newTaskCallbackDispatcher(newTaskCallbackHTTPClient(nil, nil, nil), time.Now, system_setting.GetTaskCallbackSetting)
}

// WithClock replaces the dispatcher's clock (used by end-to-end tests to
// step through retry intervals without sleeping).
func (d *TaskCallbackDispatcher) WithClock(now func() time.Time) *TaskCallbackDispatcher {
	if now != nil {
		d.now = now
	}
	return d
}

func newTaskCallbackDispatcher(client *http.Client, now func() time.Time, getSetting func() *system_setting.TaskCallbackSetting) *TaskCallbackDispatcher {
	return &TaskCallbackDispatcher{
		owner:      newTaskCallbackLeaseOwner(),
		client:     client,
		now:        now,
		getSetting: getSetting,
		wake:       taskCallbackWake,
		inflight:   make(map[int64]struct{}),
		tasks:      make(map[string]struct{}),
		hosts:      make(map[string]int),
	}
}

func newTaskCallbackLeaseOwner() string {
	buf := make([]byte, 6)
	_, _ = rand.Read(buf)
	node := common.NodeName
	if node == "" {
		node = "node"
	}
	if len(node) > 48 {
		node = node[:48]
	}
	return node + "-" + hex.EncodeToString(buf)
}

// StartTaskCallbackDispatcher starts this instance's dispatcher. The returned
// stop function cancels it, waits briefly for in-flight deliveries and
// releases this instance's leases.
func StartTaskCallbackDispatcher() (stop func()) {
	ctx, cancel := context.WithCancel(context.Background())
	dispatcher := NewTaskCallbackDispatcher()
	done := make(chan struct{})
	go func() {
		defer close(done)
		dispatcher.Run(ctx)
	}()
	return func() {
		cancel()
		select {
		case <-done:
		case <-time.After(taskCallbackShutdownWait + 5*time.Second):
		}
	}
}

// Run scans for due events every second (or when woken), reclaims expired
// leases at start and every minute, and on shutdown waits briefly for
// in-flight deliveries and releases this instance's leases.
func (d *TaskCallbackDispatcher) Run(ctx context.Context) {
	common.SysLog(common.LogText("task_callback event=dispatcher_started owner=%s", d.owner))
	d.reclaim()
	ticker := time.NewTicker(taskCallbackPollInterval)
	defer ticker.Stop()
	lastReclaim := d.now()
	for {
		select {
		case <-ctx.Done():
			d.shutdown()
			return
		case <-ticker.C:
		case <-d.wake:
		}
		if d.now().Sub(lastReclaim) >= taskCallbackReclaimInterval {
			d.reclaim()
			lastReclaim = d.now()
		}
		d.DispatchOnce(ctx)
	}
}

func (d *TaskCallbackDispatcher) shutdown() {
	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(taskCallbackShutdownWait):
	}
	if err := model.ReleaseTaskCallbackLeases(d.owner); err != nil {
		common.SysError(common.LogText("task_callback event=release_leases_failed owner=%s error=%q", d.owner, err.Error()))
	}
	common.SysLog(common.LogText("task_callback event=dispatcher_stopped owner=%s", d.owner))
}

func (d *TaskCallbackDispatcher) reclaim() {
	reclaimed, failed, err := model.ReclaimExpiredTaskCallbackLeases(d.now().Unix())
	if err != nil {
		common.SysError(common.LogText("task_callback event=reclaim_failed error=%q", err.Error()))
		return
	}
	if reclaimed > 0 || failed > 0 {
		common.SysLog(common.LogText("task_callback event=leases_reclaimed reclaimed=%d failed=%d", reclaimed, failed))
		if reclaimed > 0 {
			WakeTaskCallbackDispatcher()
		}
	}
}

// Wait blocks until in-flight deliveries started by DispatchOnce finish.
func (d *TaskCallbackDispatcher) Wait() {
	d.wg.Wait()
}

// DispatchOnce claims due events up to the worker limit and delivers each in
// its own goroutine. It reports how many deliveries it started.
func (d *TaskCallbackDispatcher) DispatchOnce(ctx context.Context) int {
	setting := d.getSetting()
	if !setting.Enabled {
		return 0
	}
	free := setting.WorkerCount() - d.inflightCount()
	if free <= 0 {
		return 0
	}
	now := d.now().Unix()
	events, err := model.ListDueTaskCallbackEvents(now, free*4)
	if err != nil {
		common.SysError(common.LogText("task_callback event=scan_failed error=%q", err.Error()))
		return 0
	}
	started := 0
	for _, event := range events {
		if started >= free || ctx.Err() != nil {
			break
		}
		if d.tryStart(ctx, event, setting) {
			started++
		}
	}
	return started
}

func (d *TaskCallbackDispatcher) inflightCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.inflight)
}

func (d *TaskCallbackDispatcher) tryStart(ctx context.Context, event *model.TaskCallbackEvent, setting *system_setting.TaskCallbackSetting) bool {
	d.mu.Lock()
	_, busyEvent := d.inflight[event.ID]
	_, busyTask := d.tasks[event.TaskID]
	d.mu.Unlock()
	if busyEvent || busyTask {
		return false
	}

	// Per-task ordering: stale queued/running events are superseded by this
	// newer one (normally already done at enqueue); anything else still open,
	// e.g. in flight on another instance, must finish first.
	if _, err := model.SupersedeTaskCallbackEvents(event.TaskID, event.Seq, model.TaskCallbackSupersedableStatuses); err != nil {
		common.SysError(common.LogText("task_callback event=supersede_failed task_id=%s error=%q", event.TaskID, err.Error()))
		return false
	}
	earlier, err := model.HasEarlierOpenTaskCallbackEvent(event.TaskID, event.Seq)
	if err != nil || earlier {
		return false
	}

	subscription, found, err := model.GetTaskCallbackSubscriptionByID(event.SubscriptionID)
	if err != nil {
		return false
	}
	if !found || subscription.CallbackURL == "" {
		d.cancelOrphan(event)
		return false
	}
	host := subscription.CallbackHost
	d.mu.Lock()
	if d.hosts[host] >= setting.HostConcurrency() {
		d.mu.Unlock()
		return false
	}
	d.inflight[event.ID] = struct{}{}
	d.tasks[event.TaskID] = struct{}{}
	d.hosts[host]++
	d.mu.Unlock()

	now := d.now()
	claimed, err := model.ClaimTaskCallbackEvent(event.ID, d.owner, now.Unix(), now.Add(taskCallbackLease).Unix())
	if err != nil || !claimed {
		d.release(event, host)
		return false
	}
	event.Attempts++
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		defer d.release(event, host)
		d.deliver(ctx, event, subscription)
	}()
	return true
}

func (d *TaskCallbackDispatcher) release(event *model.TaskCallbackEvent, host string) {
	d.mu.Lock()
	delete(d.inflight, event.ID)
	delete(d.tasks, event.TaskID)
	d.hosts[host]--
	if d.hosts[host] <= 0 {
		delete(d.hosts, host)
	}
	d.mu.Unlock()
}

func (d *TaskCallbackDispatcher) cancelOrphan(event *model.TaskCallbackEvent) {
	now := d.now()
	if claimed, err := model.ClaimTaskCallbackEvent(event.ID, d.owner, now.Unix(), now.Add(taskCallbackLease).Unix()); err != nil || !claimed {
		return
	}
	_, _ = model.FinishTaskCallbackAttempt(event.ID, d.owner, model.TaskCallbackAttemptResult{
		State: model.TaskCallbackStateCancelled,
		Error: "subscription not found",
	})
}

// taskCallbackAttempt is the outcome of one HTTP delivery attempt.
type taskCallbackAttempt struct {
	result     string
	httpStatus int
	err        error
	latency    time.Duration
}

func (a taskCallbackAttempt) succeeded() bool {
	return a.result == taskCallbackResultSucceeded
}

func (d *TaskCallbackDispatcher) deliver(parent context.Context, event *model.TaskCallbackEvent, subscription *model.TaskCallbackSubscription) {
	setting := d.getSetting()
	attempt := d.send(parent, event, subscription, setting)
	now := d.now()
	result := model.TaskCallbackAttemptResult{
		HTTPStatus: attempt.httpStatus,
		LatencyMs:  int(attempt.latency.Milliseconds()),
	}
	nextIn := time.Duration(0)
	switch {
	case attempt.succeeded():
		result.State = model.TaskCallbackStateSucceeded
		result.DeliveredAt = now.Unix()
	case attempt.result == taskCallbackResultSSRFBlocked || attempt.result == taskCallbackResultInvalid:
		// The outcome cannot change on retry.
		result.State = model.TaskCallbackStateFailed
	case event.Attempts < event.MaxAttempts:
		nextIn = setting.RetryInterval()
		result.State = model.TaskCallbackStatePending
	default:
		if delay, ok := taskCallbackExtendedDelay(setting, event, now); ok {
			nextIn = delay
			result.State = model.TaskCallbackStatePending
			result.MaxAttempts = event.Attempts + 1
		} else {
			result.State = model.TaskCallbackStateFailed
		}
	}
	if nextIn > 0 {
		result.NextAttemptAt = now.Add(nextIn).Unix()
	}
	if !attempt.succeeded() {
		result.Error = attempt.result
		if attempt.err != nil {
			result.Error = attempt.result + ": " + sanitizeTaskCallbackError(attempt.err, subscription.CallbackURL)
		}
	}
	applied, err := model.FinishTaskCallbackAttempt(event.ID, d.owner, result)
	if err != nil {
		common.SysError(common.LogText("task_callback event=finish_failed event_id=%s error=%q", event.EventID, err.Error()))
	} else if !applied {
		common.SysLog(common.LogText("task_callback event=lease_lost event_id=%s task_id=%s", event.EventID, event.TaskID))
	}
	common.SysLog(common.LogText(
		"task_callback event=delivery_attempt event_id=%s task_id=%s status=%s attempt=%d max_attempts=%d result=%s http_status=%d latency_ms=%d state=%s next_attempt_in_s=%d host=%s url_hash=%s",
		event.EventID, event.TaskID, event.Status, event.Attempts, event.MaxAttempts, attempt.result, attempt.httpStatus,
		attempt.latency.Milliseconds(), result.State, int64(nextIn/time.Second), subscription.CallbackHost,
		TaskCallbackURLHashPrefix(subscription.CallbackURLHash),
	))
	if result.State == model.TaskCallbackStatePending {
		WakeTaskCallbackDispatcher()
	}
}

// taskCallbackExtendedDelay returns the next back-off step when the optional
// extended retry applies (terminal succeeded/failed events only).
func taskCallbackExtendedDelay(setting *system_setting.TaskCallbackSetting, event *model.TaskCallbackEvent, now time.Time) (time.Duration, bool) {
	if !setting.ExtendedRetry || !taskCallbackDocumentedRetries(event.Status) {
		return 0, false
	}
	step := event.Attempts - setting.TerminalAttempts()
	if step < 0 {
		step = 0
	}
	if step >= len(taskCallbackExtendedBackoff) {
		step = len(taskCallbackExtendedBackoff) - 1
	}
	delay := taskCallbackExtendedBackoff[step]
	if now.Add(delay).Sub(time.Unix(event.CreatedAt, 0)) > taskCallbackExtendedRetryWindow {
		return 0, false
	}
	return delay, true
}

func (d *TaskCallbackDispatcher) send(parent context.Context, event *model.TaskCallbackEvent, subscription *model.TaskCallbackSubscription, setting *system_setting.TaskCallbackSetting) taskCallbackAttempt {
	ctx, cancel := context.WithTimeout(parent, setting.Timeout())
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, subscription.CallbackURL, strings.NewReader(event.Payload))
	if err != nil {
		return taskCallbackAttempt{result: taskCallbackResultInvalid, err: err}
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "kt-api-callback/"+common.Version)
	if setting.ExtraHeadersEnabled {
		request.Header.Set("X-KtApi-Event-Id", event.EventID)
		request.Header.Set("X-KtApi-Delivery-Attempt", strconv.Itoa(event.Attempts))
		request.Header.Set("X-KtApi-Task-Id", event.TaskID)
	}
	started := d.now()
	signTaskCallbackRequest(request, setting, event.UserID, event.EventID, []byte(event.Payload), started.Unix())
	response, err := d.client.Do(request)
	if err != nil {
		return taskCallbackAttempt{result: classifyTaskCallbackError(err), err: err, latency: d.now().Sub(started)}
	}
	// The response body is ignored; at most 64 KiB is drained so the
	// connection can be reused.
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, taskCallbackMaxResponseBytes))
	_ = response.Body.Close()
	latency := d.now().Sub(started)
	status := response.StatusCode
	if status >= 200 && status < 300 {
		// Any 2xx within the timeout acknowledges the callback.
		return taskCallbackAttempt{result: taskCallbackResultSucceeded, httpStatus: status, latency: latency}
	}
	return taskCallbackAttempt{result: fmt.Sprintf("http_%dxx", status/100), httpStatus: status, latency: latency}
}

func classifyTaskCallbackError(err error) string {
	if IsTaskCallbackBlocked(err) {
		return taskCallbackResultSSRFBlocked
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return taskCallbackResultDNS
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return taskCallbackResultTimeout
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return taskCallbackResultTimeout
	}
	var certErr *tls.CertificateVerificationError
	var unknownAuthority x509.UnknownAuthorityError
	var hostnameErr x509.HostnameError
	var recordErr tls.RecordHeaderError
	if errors.As(err, &certErr) || errors.As(err, &unknownAuthority) || errors.As(err, &hostnameErr) || errors.As(err, &recordErr) || strings.Contains(err.Error(), "tls:") {
		return taskCallbackResultTLS
	}
	return taskCallbackResultConnect
}

// sanitizeTaskCallbackError removes the callback URL (whose path and query
// may carry customer secrets) from an error message.
func sanitizeTaskCallbackError(err error, callbackURL string) string {
	message := err.Error()
	if callbackURL != "" {
		message = strings.ReplaceAll(message, callbackURL, "<callback_url>")
		message = strings.ReplaceAll(message, strconv.Quote(callbackURL), "<callback_url>")
	}
	if len(message) > 400 {
		message = message[:400]
	}
	return message
}
