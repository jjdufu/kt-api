package controller

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/relay"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	callbackE2EUserID = 31
	callbackE2EModel  = "doubao-seedance-2-0-260128"
)

// callbackE2ESupplier is a programmable Ark-compatible supplier. Each
// upstream task id returns the scripted status/body on GET.
type callbackE2ESupplier struct {
	t      *testing.T
	server *httptest.Server

	mu          sync.Mutex
	seq         int
	submits     []map[string]any
	queryStatus map[string]int
	queryBody   map[string]string
}

func newCallbackE2ESupplier(t *testing.T) *callbackE2ESupplier {
	supplier := &callbackE2ESupplier{t: t, queryStatus: map[string]int{}, queryBody: map[string]string{}}
	supplier.server = httptest.NewServer(http.HandlerFunc(supplier.serve))
	t.Cleanup(supplier.server.Close)
	return supplier
}

func (s *callbackE2ESupplier) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	const prefix = "/api/v3/contents/generations/tasks"
	s.mu.Lock()
	defer s.mu.Unlock()
	switch {
	case r.Method == http.MethodPost && r.URL.Path == prefix:
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = common.Unmarshal(raw, &body)
		s.submits = append(s.submits, body)
		s.seq++
		id := "cgt-e2e-" + string(rune('0'+s.seq))
		s.queryStatus[id] = http.StatusOK
		s.queryBody[id] = `{"id":"` + id + `","model":"` + callbackE2EModel + `","status":"queued","created_at":1786015617,"updated_at":1786015617}`
		_, _ = io.WriteString(w, `{"id":"`+id+`"}`)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, prefix+"/"):
		id := strings.TrimPrefix(r.URL.Path, prefix+"/")
		body, ok := s.queryBody[id]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(s.queryStatus[id])
		_, _ = io.WriteString(w, body)
	default:
		http.NotFound(w, r)
	}
}

func (s *callbackE2ESupplier) script(id string, status int, body string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queryStatus[id] = status
	s.queryBody[id] = body
}

func (s *callbackE2ESupplier) lastSubmit() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	require.NotEmpty(s.t, s.submits)
	return s.submits[len(s.submits)-1]
}

// callbackE2EReceiver records customer callbacks and replies with scripted
// status codes (200 once the script runs out).
type callbackE2EReceiver struct {
	server *httptest.Server

	mu       sync.Mutex
	bodies   []string
	headers  []http.Header
	replies  []int
	received atomic.Int32
}

func newCallbackE2EReceiver(t *testing.T) *callbackE2EReceiver {
	receiver := &callbackE2EReceiver{}
	receiver.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		receiver.mu.Lock()
		receiver.bodies = append(receiver.bodies, string(raw))
		receiver.headers = append(receiver.headers, r.Header.Clone())
		code := http.StatusOK
		if len(receiver.replies) > 0 {
			code = receiver.replies[0]
			receiver.replies = receiver.replies[1:]
		}
		receiver.mu.Unlock()
		receiver.received.Add(1)
		if code == http.StatusFound {
			w.Header().Set("Location", "http://169.254.169.254/latest")
		}
		w.WriteHeader(code)
	}))
	t.Cleanup(receiver.server.Close)
	return receiver
}

func (r *callbackE2EReceiver) snapshot() ([]string, []http.Header) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.bodies...), append([]http.Header(nil), r.headers...)
}

func (r *callbackE2EReceiver) reset(replies ...int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bodies, r.headers, r.replies = nil, nil, replies
}

type callbackE2E struct {
	t        *testing.T
	database *gorm.DB
	supplier *callbackE2ESupplier
	receiver *callbackE2EReceiver
	channel  *model.Channel
	engine   *gin.Engine
	clock    atomic.Int64 // offset seconds added to time.Now
	billing  []*nativeRouteBilling
}

func setupCallbackE2E(t *testing.T) *callbackE2E {
	t.Helper()
	gin.SetMode(gin.TestMode)
	service.InitHttpClient()

	previousDB, previousLogDB := model.DB, model.LOG_DB
	previousMemoryCache, previousBatchUpdate := common.MemoryCacheEnabled, common.BatchUpdateEnabled
	previousLogConsume, previousRedis := common.LogConsumeEnabled, common.RedisEnabled
	dsn := filepath.Join(t.TempDir(), "e2e.db") + "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)"
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.User{}, &model.Channel{}, &model.Task{}, &model.Log{},
		&model.TaskCallbackSubscription{}, &model.TaskCallbackEvent{}))
	model.DB, model.LOG_DB = database, database
	common.MemoryCacheEnabled, common.BatchUpdateEnabled = false, false
	common.LogConsumeEnabled, common.RedisEnabled = true, false

	previousRatios := ratio_setting.ModelRatio2JSONString()
	require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{"`+callbackE2EModel+`":1}`))
	setting := system_setting.GetTaskCallbackSetting()
	previousSetting := *setting
	previousFactory := service.GetTaskAdaptorFunc
	service.GetTaskAdaptorFunc = func(platform constant.TaskPlatform) service.TaskPollingAdaptor {
		return relay.GetTaskAdaptor(platform)
	}
	t.Cleanup(func() {
		service.GetTaskAdaptorFunc = previousFactory
		*setting = previousSetting
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(previousRatios))
		model.DB, model.LOG_DB = previousDB, previousLogDB
		common.MemoryCacheEnabled, common.BatchUpdateEnabled = previousMemoryCache, previousBatchUpdate
		common.LogConsumeEnabled, common.RedisEnabled = previousLogConsume, previousRedis
		if sqlDB, err := database.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	env := &callbackE2E{t: t, database: database, supplier: newCallbackE2ESupplier(t), receiver: newCallbackE2EReceiver(t)}
	_, receiverPort, err := net.SplitHostPort(strings.TrimPrefix(env.receiver.server.URL, "http://"))
	require.NoError(t, err)
	setting.Enabled = true
	setting.ExtraHeadersEnabled = true
	setting.AllowPrivateIp = true // the mock receiver listens on loopback
	setting.AllowedPorts = []string{receiverPort}
	setting.RetryIntervalSeconds = 5
	setting.TerminalMaxAttempts = 4
	setting.TimeoutSeconds = 2
	setting.ExtendedRetry = false
	setting.UpstreamPublicBase = ""

	require.NoError(t, database.Create(&model.User{Id: callbackE2EUserID, Username: "callback-e2e", Group: "default", Quota: 1_000_000}).Error)
	baseURL := env.supplier.server.URL
	env.channel = &model.Channel{
		Type: constant.ChannelTypeDoubaoVideo, Name: "ark-e2e", Key: "sk-supplier", BaseURL: &baseURL,
		Status: common.ChannelStatusEnabled, Models: callbackE2EModel, Group: "default",
	}
	require.NoError(t, database.Create(env.channel).Error)

	env.engine = gin.New()
	authenticate := func(c *gin.Context) {
		common.SetContextKey(c, constant.ContextKeyUserId, callbackE2EUserID)
		common.SetContextKey(c, constant.ContextKeyUserGroup, "default")
		common.SetContextKey(c, constant.ContextKeyUsingGroup, "default")
		common.SetContextKey(c, constant.ContextKeyTokenGroup, "default")
		common.SetContextKey(c, constant.ContextKeyUserQuota, 1_000_000)
	}
	distribute := func(c *gin.Context) {
		if err := middleware.SetupContextForSelectedChannel(c, env.channel, callbackE2EModel); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
	env.engine.POST("/api/v3/contents/generations/tasks", authenticate, middleware.PinSeedanceVideoAlias(),
		middleware.PrepareTaskPluginRoute(), middleware.SeedanceVideoTaskOptions(), distribute, env.submit)
	env.engine.GET("/api/v3/contents/generations/tasks/:task_id", authenticate, middleware.PinSeedanceVideoAlias(),
		middleware.PrepareTaskPluginRoute())
	env.engine.POST("/api/task-callback/upstream/:token", HandleTaskUpstreamCallback)
	return env
}

func (e *callbackE2E) submit(c *gin.Context) {
	billing := &nativeRouteBilling{userID: callbackE2EUserID}
	e.billing = append(e.billing, billing)
	channel := e.channel
	info := &relaycommon.RelayInfo{
		UserId: callbackE2EUserID, UserGroup: "default", UsingGroup: "default", UserQuota: 1_000_000, TokenGroup: "default",
		OriginModelName: callbackE2EModel, Billing: billing,
		TaskRelayInfo: &relaycommon.TaskRelayInfo{
			Action:        c.GetString("task_action"),
			PublicTaskID:  "task_e2e_" + string(rune('a'+len(e.billing)-1)),
			LockedChannel: channel,
		},
	}
	outcome, taskErr := executeTaskSubmissionWith(c, info, relay.RelayTaskSubmit)
	if taskErr != nil {
		c.AbortWithStatusJSON(taskErr.StatusCode, gin.H{"error": taskErr.Message})
		return
	}
	presentTaskSubmission(c, outcome)
}

func (e *callbackE2E) do(method, path, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	var reader io.Reader
	if body != "" {
		reader = bytes.NewBufferString(body)
	}
	request := httptest.NewRequest(method, path, reader)
	request.Header.Set("Content-Type", "application/json")
	e.engine.ServeHTTP(recorder, request)
	return recorder
}

func (e *callbackE2E) submitTask(callbackURL string) (string, *httptest.ResponseRecorder) {
	body := `{"model":"` + callbackE2EModel + `","content":[{"type":"text","text":"a lighthouse"}]`
	if callbackURL != "" {
		body += `,"callback_url":"` + callbackURL + `"`
	}
	body += `}`
	recorder := e.do(http.MethodPost, "/api/v3/contents/generations/tasks", body)
	var created map[string]any
	_ = common.Unmarshal(recorder.Body.Bytes(), &created)
	id, _ := created["id"].(string)
	return id, recorder
}

func (e *callbackE2E) task(taskID string) *model.Task {
	var task model.Task
	require.NoError(e.t, e.database.Where("task_id = ?", taskID).First(&task).Error)
	return &task
}

func (e *callbackE2E) poll(taskID string) {
	task := e.task(taskID)
	upstreamID := task.PrivateData.UpstreamTaskID
	service.DispatchPlatformUpdate(context.Background(), task.Platform,
		map[int][]string{task.ChannelId: {upstreamID}}, map[string]*model.Task{upstreamID: task})
}

func (e *callbackE2E) get(taskID string) string {
	recorder := e.do(http.MethodGet, "/api/v3/contents/generations/tasks/"+taskID, "")
	require.Equal(e.t, http.StatusOK, recorder.Code, recorder.Body.String())
	return recorder.Body.String()
}

func (e *callbackE2E) dispatcher() *service.TaskCallbackDispatcher {
	return service.NewTaskCallbackDispatcher().WithClock(func() time.Time {
		return time.Now().Add(time.Duration(e.clock.Load()) * time.Second)
	})
}

// deliver runs one dispatch pass and waits for its deliveries.
func (e *callbackE2E) deliver() int {
	dispatcher := e.dispatcher()
	started := dispatcher.DispatchOnce(context.Background())
	dispatcher.Wait()
	return started
}

func (e *callbackE2E) events(taskID string) []model.TaskCallbackEvent {
	var events []model.TaskCallbackEvent
	require.NoError(e.t, e.database.Where("task_id = ?", taskID).Order("seq").Find(&events).Error)
	return events
}

// liveEvents drops the queued event a later status superseded before it was
// delivered (every submit enqueues queued; the tests below often move the task
// on before running the dispatcher).
func (e *callbackE2E) liveEvents(taskID string) []model.TaskCallbackEvent {
	var live []model.TaskCallbackEvent
	for _, event := range e.events(taskID) {
		if event.State == model.TaskCallbackStateSuperseded {
			assert.Equal(e.t, "queued", event.Status, "only the submit-time queued event is superseded here")
			continue
		}
		live = append(live, event)
	}
	return live
}

func (e *callbackE2E) quotaAndLogs() (int, int64) {
	var user model.User
	require.NoError(e.t, e.database.First(&user, callbackE2EUserID).Error)
	var logs int64
	require.NoError(e.t, e.database.Model(&model.Log{}).Count(&logs).Error)
	return user.Quota, logs
}

func arkBody(id, status, extra string) string {
	return `{"id":"` + id + `","model":"` + callbackE2EModel + `","status":"` + status + `","created_at":1786015617,"updated_at":1786015700` + extra + `}`
}

func TestTaskCallbackEndToEnd(t *testing.T) {
	env := setupCallbackE2E(t)
	callbackURL := env.receiver.server.URL + "/hooks/seedance?order=42"

	t.Run("private callback_url is rejected before billing", func(t *testing.T) {
		system_setting.GetTaskCallbackSetting().AllowPrivateIp = false
		defer func() { system_setting.GetTaskCallbackSetting().AllowPrivateIp = true }()
		quota, logs := env.quotaAndLogs()
		_, recorder := env.submitTask(callbackURL)
		assert.Equal(t, http.StatusBadRequest, recorder.Code, recorder.Body.String())
		assert.Contains(t, recorder.Body.String(), "callback_url")
		var tasks int64
		require.NoError(t, env.database.Model(&model.Task{}).Count(&tasks).Error)
		assert.Zero(t, tasks)
		assert.Empty(t, env.supplier.submits)
		gotQuota, gotLogs := env.quotaAndLogs()
		assert.Equal(t, quota, gotQuota)
		assert.Equal(t, logs, gotLogs)
	})

	t.Run("every transition is delivered once with the GET bytes", func(t *testing.T) {
		env.receiver.reset()
		taskID, recorder := env.submitTask(callbackURL)
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		require.NotEmpty(t, taskID)
		assert.NotContains(t, recorder.Body.String(), "cgt-e2e")

		submitted := env.supplier.lastSubmit()
		assert.NotContains(t, submitted, "callback_url", "the customer URL never reaches the supplier")
		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID

		subscription, found, err := model.GetTaskCallbackSubscriptionByTaskID(taskID)
		require.NoError(t, err)
		require.True(t, found)
		assert.Equal(t, callbackURL, subscription.CallbackURL)
		assert.Nil(t, subscription.UpstreamTokenHash, "fast path is off by default")

		// Before the first poll GET already answers a documented queued body,
		// and the submit itself enqueued the queued callback with those bytes.
		task := env.task(taskID)
		queuedGET := env.get(taskID)
		assert.JSONEq(t, `{"id":"`+taskID+`","model":"`+callbackE2EModel+`","status":"queued","created_at":`+
			strconv.FormatInt(task.CreatedAt, 10)+`,"updated_at":`+strconv.FormatInt(task.UpdatedAt, 10)+`}`, queuedGET)
		getBodies := []string{queuedGET}
		require.Equal(t, 1, env.deliver(), "the submit enqueues queued")

		step := func(status string, extra string, deliveries int) {
			env.supplier.script(upstreamID, http.StatusOK, arkBody(upstreamID, status, extra))
			env.poll(taskID)
			if deliveries > 0 {
				getBodies = append(getBodies, env.get(taskID))
			}
			require.Equal(t, deliveries, env.deliver(), "deliveries for %s", status)
			// Repeating the same status does not deliver again.
			env.poll(taskID)
			assert.Zero(t, env.deliver())
		}
		// The supplier's own queued body does not resend queued.
		step("queued", "", 0)
		step("running", "", 1)
		step("succeeded", `,"content":{"video_url":"https://cdn.example/v.mp4?a=1&b=2"},"usage":{"completion_tokens":1000,"total_tokens":1000}`, 1)

		bodies, headers := env.receiver.snapshot()
		require.Len(t, bodies, 3)
		assert.Equal(t, getBodies, bodies, "callback body is byte-identical to GET at that moment")
		for index, body := range bodies {
			assert.Contains(t, body, `"id":"`+taskID+`"`)
			assert.NotContains(t, body, upstreamID)
			assert.NotContains(t, body, env.supplier.server.URL)
			assert.Equal(t, "application/json", headers[index].Get("Content-Type"))
			assert.Regexp(t, `^kt-api-callback/`, headers[index].Get("User-Agent"))
			assert.Regexp(t, `^evt_[0-9a-f]{32}$`, headers[index].Get("X-KtApi-Event-Id"))
			assert.Equal(t, "1", headers[index].Get("X-KtApi-Delivery-Attempt"))
			assert.Equal(t, taskID, headers[index].Get("X-KtApi-Task-Id"))
		}
		assert.Contains(t, bodies[0], `"status":"queued"`)
		assert.Contains(t, bodies[1], `"status":"running"`)
		assert.Contains(t, bodies[2], `"status":"succeeded"`)
		events := env.events(taskID)
		require.Len(t, events, 3)
		for index, event := range events {
			assert.Equal(t, []string{"queued", "running", "succeeded"}[index], event.Status)
			assert.Equal(t, index+1, event.Seq)
			assert.Equal(t, model.TaskCallbackStateSucceeded, event.State)
		}
		enqueued, err := service.EnqueueInitialTaskCallback(taskID)
		require.NoError(t, err)
		assert.False(t, enqueued, "the initial enqueue never runs twice")

		// Billing: callbacks (and a redelivery) change neither quota nor logs.
		quota, logs := env.quotaAndLogs()
		require.NoError(t, env.database.Model(&model.TaskCallbackEvent{}).Where("id = ?", events[len(events)-1].ID).
			Updates(map[string]any{"state": model.TaskCallbackStatePending, "attempts": 0}).Error)
		require.Equal(t, 1, env.deliver())
		gotQuota, gotLogs := env.quotaAndLogs()
		assert.Equal(t, quota, gotQuota)
		assert.Equal(t, logs, gotLogs)
	})

	t.Run("no queued after a later status was observed first", func(t *testing.T) {
		env.receiver.reset()
		taskID, recorder := env.submitTask(callbackURL)
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		// Rewind to the moment before the submit-time enqueue, then let the
		// poller observe running first.
		require.NoError(t, env.database.Where("task_id = ?", taskID).Delete(&model.TaskCallbackEvent{}).Error)
		require.NoError(t, env.database.Model(&model.TaskCallbackSubscription{}).Where("task_id = ?", taskID).
			Updates(map[string]any{"last_status": "", "last_seq": 0}).Error)
		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID
		env.supplier.script(upstreamID, http.StatusOK, arkBody(upstreamID, "running", ""))
		env.poll(taskID)
		enqueued, err := service.EnqueueInitialTaskCallback(taskID)
		require.NoError(t, err)
		assert.False(t, enqueued)
		require.Equal(t, 1, env.deliver())
		bodies, _ := env.receiver.snapshot()
		require.Len(t, bodies, 1)
		assert.Contains(t, bodies[0], `"status":"running"`)
		assert.Equal(t, env.get(taskID), bodies[0])
		events := env.events(taskID)
		require.Len(t, events, 1)
		assert.Equal(t, "running", events[0].Status)
	})

	t.Run("terminal retries: three 500s then 200 is four attempts", func(t *testing.T) {
		env.receiver.reset(http.StatusInternalServerError, http.StatusInternalServerError, http.StatusInternalServerError)
		taskID, recorder := env.submitTask(callbackURL)
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID
		// Supplier failure: HTTP 400 with a bare error object.
		env.supplier.script(upstreamID, http.StatusBadRequest, `{"error":{"code":"InputTextSensitiveContentDetected","message":"sensitive","type":"BadRequest"}}`)
		env.poll(taskID)
		failedGET := env.get(taskID)
		assert.Contains(t, failedGET, `"status":"failed"`)
		assert.Contains(t, failedGET, `"error":{`)

		for attempt := 1; attempt <= 4; attempt++ {
			require.Equal(t, 1, env.deliver(), "attempt %d", attempt)
			assert.Zero(t, env.deliver(), "nothing is due before the retry interval")
			env.clock.Add(5)
		}
		bodies, headers := env.receiver.snapshot()
		require.Len(t, bodies, 4)
		for index, body := range bodies {
			assert.Equal(t, failedGET, body)
			assert.Equal(t, string(rune('1'+index)), headers[index].Get("X-KtApi-Delivery-Attempt"))
		}
		assert.Equal(t, headers[0].Get("X-KtApi-Event-Id"), headers[3].Get("X-KtApi-Event-Id"), "retries keep the event id")
		events := env.liveEvents(taskID)
		require.Len(t, events, 1)
		assert.Equal(t, model.TaskCallbackStateSucceeded, events[0].State)
		assert.Equal(t, 4, events[0].Attempts)
	})

	t.Run("redirects are not followed and count as failures", func(t *testing.T) {
		env.receiver.reset(http.StatusFound, http.StatusFound, http.StatusFound, http.StatusFound)
		taskID, _ := env.submitTask(callbackURL)
		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID
		env.supplier.script(upstreamID, http.StatusOK, arkBody(upstreamID, "expired", ""))
		env.poll(taskID)
		for attempt := 0; attempt < 4; attempt++ {
			env.deliver()
			env.clock.Add(5)
		}
		events := env.liveEvents(taskID)
		require.Len(t, events, 1)
		assert.Equal(t, "expired", events[0].Status)
		assert.Equal(t, 1, events[0].MaxAttempts, "non-terminal (expired) events are sent once")
		assert.Equal(t, model.TaskCallbackStateFailed, events[0].State)
		assert.EqualValues(t, 1, len(func() []string { b, _ := env.receiver.snapshot(); return b }()))
	})

	t.Run("dial-time private target is ssrf_blocked without retries", func(t *testing.T) {
		env.receiver.reset()
		taskID, recorder := env.submitTask(callbackURL)
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID
		env.supplier.script(upstreamID, http.StatusOK, arkBody(upstreamID, "succeeded", `,"content":{"video_url":"https://cdn.example/v.mp4"}`))
		env.poll(taskID)
		system_setting.GetTaskCallbackSetting().AllowPrivateIp = false
		defer func() { system_setting.GetTaskCallbackSetting().AllowPrivateIp = true }()
		require.Equal(t, 1, env.deliver())
		env.clock.Add(5)
		assert.Zero(t, env.deliver())
		events := env.liveEvents(taskID)
		require.Len(t, events, 1)
		assert.Equal(t, model.TaskCallbackStateFailed, events[0].State)
		assert.True(t, strings.HasPrefix(events[0].LastError, "ssrf_blocked:"), events[0].LastError)
		assert.NotContains(t, events[0].LastError, "/hooks/seedance", "the URL is redacted")
		assert.Equal(t, int32(0), int32(len(func() []string { b, _ := env.receiver.snapshot(); return b }())))
	})

	t.Run("crash recovery and concurrent dispatchers deliver once", func(t *testing.T) {
		env.receiver.reset()
		taskID, recorder := env.submitTask(callbackURL)
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID
		env.supplier.script(upstreamID, http.StatusOK, arkBody(upstreamID, "succeeded", `,"content":{"video_url":"https://cdn.example/v.mp4"}`))
		env.poll(taskID)
		events := env.liveEvents(taskID)
		require.Len(t, events, 1)
		// A dispatcher "crashed" mid-delivery and left an expired lease.
		require.NoError(t, env.database.Model(&model.TaskCallbackEvent{}).Where("id = ?", events[0].ID).Updates(map[string]any{
			"state": model.TaskCallbackStateDelivering, "attempts": 1, "lease_owner": "dead-node", "lease_until": time.Now().Unix() - 1,
		}).Error)
		assert.Zero(t, env.deliver(), "a leased event is not claimed")
		_, _, err := model.ReclaimExpiredTaskCallbackLeases(time.Now().Unix())
		require.NoError(t, err)
		env.clock.Add(5)

		first, second := env.dispatcher(), env.dispatcher()
		var wg sync.WaitGroup
		var started atomic.Int32
		for _, dispatcher := range []*service.TaskCallbackDispatcher{first, second} {
			wg.Add(1)
			go func(d *service.TaskCallbackDispatcher) {
				defer wg.Done()
				started.Add(int32(d.DispatchOnce(context.Background())))
				d.Wait()
			}(dispatcher)
		}
		wg.Wait()
		assert.Equal(t, int32(1), started.Load(), "only one dispatcher claims the event")
		bodies, headers := env.receiver.snapshot()
		require.Len(t, bodies, 1)
		assert.Equal(t, events[0].EventID, headers[0].Get("X-KtApi-Event-Id"), "event id is stable across the crash")
		assert.Equal(t, "2", headers[0].Get("X-KtApi-Delivery-Attempt"))
		assert.Equal(t, model.TaskCallbackStateSucceeded, env.liveEvents(taskID)[0].State)
	})

	t.Run("fast path re-queries the supplier and ignores the hint body", func(t *testing.T) {
		setting := system_setting.GetTaskCallbackSetting()
		setting.UpstreamPublicBase = "https://kt.example"
		defer func() { setting.UpstreamPublicBase = "" }()
		require.NoError(t, env.database.Model(env.channel).Update("settings", `{"upstream_task_callback_enabled":true}`).Error)
		env.channel.OtherSettings = `{"upstream_task_callback_enabled":true}`
		defer func() {
			env.channel.OtherSettings = ""
			_ = env.database.Model(env.channel).Update("settings", "").Error
		}()
		env.receiver.reset()

		taskID, recorder := env.submitTask(callbackURL)
		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		upstreamCallback, _ := env.supplier.lastSubmit()["callback_url"].(string)
		require.True(t, strings.HasPrefix(upstreamCallback, "https://kt.example/api/task-callback/upstream/"), upstreamCallback)
		assert.NotEqual(t, callbackURL, upstreamCallback)
		parsed, err := url.Parse(upstreamCallback)
		require.NoError(t, err)

		upstreamID := env.task(taskID).PrivateData.UpstreamTaskID
		env.supplier.script(upstreamID, http.StatusOK, arkBody(upstreamID, "succeeded", `,"content":{"video_url":"https://cdn.example/fast.mp4"}`))
		// The hint claims failure; kt-api trusts only its own re-query.
		hint := env.do(http.MethodPost, parsed.Path, arkBody(upstreamID, "failed", ""))
		require.Equal(t, http.StatusOK, hint.Code, hint.Body.String())
		require.Eventually(t, func() bool {
			return env.task(taskID).Status == model.TaskStatusSuccess
		}, 5*time.Second, 10*time.Millisecond)
		require.Equal(t, 1, env.deliver())
		bodies, _ := env.receiver.snapshot()
		require.Len(t, bodies, 1)
		assert.Contains(t, bodies[0], `"status":"succeeded"`)
		assert.Equal(t, env.get(taskID), bodies[0])

		// A later poll of the same terminal state settles nothing twice.
		quota, logs := env.quotaAndLogs()
		env.poll(taskID)
		assert.Zero(t, env.deliver())
		gotQuota, gotLogs := env.quotaAndLogs()
		assert.Equal(t, quota, gotQuota)
		assert.Equal(t, logs, gotLogs)
		assert.Equal(t, http.StatusNotFound, env.do(http.MethodPost, "/api/task-callback/upstream/forged", `{}`).Code)
	})
}
