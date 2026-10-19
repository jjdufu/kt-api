package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
)

// taskCallbackSubscriptionView never exposes the callback URL path, query
// or the upstream token: only the host and a hash prefix to compare URLs.
type taskCallbackSubscriptionView struct {
	CallbackHost     string `json:"callback_host"`
	CallbackURLHash  string `json:"callback_url_hash"`
	UpstreamFastPath bool   `json:"upstream_fast_path"`
	LastStatus       string `json:"last_status"`
	CreatedAt        int64  `json:"created_at"`
}

type taskCallbackEventView struct {
	EventID        string `json:"event_id"`
	TaskID         string `json:"task_id"`
	UserID         int    `json:"user_id"`
	Status         string `json:"status"`
	Seq            int    `json:"seq"`
	State          string `json:"state"`
	Attempts       int    `json:"attempts"`
	MaxAttempts    int    `json:"max_attempts"`
	NextAttemptAt  int64  `json:"next_attempt_at"`
	LastHTTPStatus int    `json:"last_http_status"`
	LastError      string `json:"last_error"`
	LastLatencyMs  int    `json:"last_latency_ms"`
	DeliveredAt    int64  `json:"delivered_at"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	PayloadSHA256  string `json:"payload_sha256"`
	CallbackHost   string `json:"callback_host,omitempty"`
	Payload        string `json:"payload,omitempty"`
	Redeliverable  bool   `json:"redeliverable"`
}

func newTaskCallbackSubscriptionView(subscription *model.TaskCallbackSubscription) *taskCallbackSubscriptionView {
	if subscription == nil {
		return nil
	}
	return &taskCallbackSubscriptionView{
		CallbackHost:     subscription.CallbackHost,
		CallbackURLHash:  service.TaskCallbackURLHashPrefix(subscription.CallbackURLHash),
		UpstreamFastPath: subscription.UpstreamTokenHash != nil,
		LastStatus:       subscription.LastStatus,
		CreatedAt:        subscription.CreatedAt,
	}
}

func newTaskCallbackEventView(event *model.TaskCallbackEvent) taskCallbackEventView {
	return taskCallbackEventView{
		EventID:        event.EventID,
		TaskID:         event.TaskID,
		UserID:         event.UserID,
		Status:         event.Status,
		Seq:            event.Seq,
		State:          event.State,
		Attempts:       event.Attempts,
		MaxAttempts:    event.MaxAttempts,
		NextAttemptAt:  event.NextAttemptAt,
		LastHTTPStatus: event.LastHTTPStatus,
		LastError:      event.LastError,
		LastLatencyMs:  event.LastLatencyMs,
		DeliveredAt:    event.DeliveredAt,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
		PayloadSHA256:  event.PayloadSHA256,
		Redeliverable:  event.State == model.TaskCallbackStateFailed || event.State == model.TaskCallbackStateCancelled,
	}
}

// GetTaskCallbacks returns the callback subscription and delivery events of
// one task. Users see their own tasks; administrators on a session see any.
func GetTaskCallbacks(c *gin.Context) {
	task, exists, err := getTaskForArtifactRequest(c, c.Param("task_id"))
	if err != nil {
		common.ApiErrorMsg(c, "failed to query task")
		return
	}
	if !exists || task == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "task not found"})
		return
	}
	subscription, found, err := model.GetTaskCallbackSubscriptionByTaskID(task.TaskID)
	if err != nil {
		common.ApiErrorMsg(c, "failed to query task callbacks")
		return
	}
	events := []taskCallbackEventView{}
	var subscriptionView *taskCallbackSubscriptionView
	if found {
		subscriptionView = newTaskCallbackSubscriptionView(subscription)
		rows, err := model.ListTaskCallbackEventsByTask(task.TaskID)
		if err != nil {
			common.ApiErrorMsg(c, "failed to query task callbacks")
			return
		}
		for _, row := range rows {
			view := newTaskCallbackEventView(row)
			view.CallbackHost = subscription.CallbackHost
			view.Payload = row.Payload
			events = append(events, view)
		}
	}
	c.Header("Cache-Control", "private, no-store")
	common.ApiSuccess(c, gin.H{
		"task_id":      task.TaskID,
		"enabled":      system_setting.GetTaskCallbackSetting().Enabled,
		"subscription": subscriptionView,
		"events":       events,
	})
}

// AdminListTaskCallbackEvents lists callback events across tasks, newest
// first, filtered by task, user, delivery state, task status and time.
func AdminListTaskCallbackEvents(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	userID, _ := strconv.Atoi(c.Query("user_id"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	query := model.TaskCallbackEventQuery{
		TaskID:         strings.TrimSpace(c.Query("task_id")),
		UserID:         userID,
		State:          strings.TrimSpace(c.Query("state")),
		Status:         strings.TrimSpace(c.Query("status")),
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}
	rows, total, err := model.ListTaskCallbackEvents(query, pageInfo.GetStartIdx(), pageInfo.GetPageSize())
	if err != nil {
		common.ApiErrorMsg(c, "failed to query task callbacks")
		return
	}
	subscriptionIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		subscriptionIDs = append(subscriptionIDs, row.SubscriptionID)
	}
	subscriptions, err := model.ListTaskCallbackSubscriptionsByIDs(subscriptionIDs)
	if err != nil {
		common.ApiErrorMsg(c, "failed to query task callbacks")
		return
	}
	items := make([]taskCallbackEventView, 0, len(rows))
	for _, row := range rows {
		view := newTaskCallbackEventView(row)
		if subscription, ok := subscriptions[row.SubscriptionID]; ok {
			view.CallbackHost = subscription.CallbackHost
		}
		items = append(items, view)
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(items)
	common.ApiSuccess(c, pageInfo)
}

// AdminRedeliverTaskCallbackEvent re-arms a failed or cancelled event with a
// fresh attempt budget; the event_id stays the same so receivers can dedupe.
func AdminRedeliverTaskCallbackEvent(c *gin.Context) {
	eventID := strings.TrimSpace(c.Param("event_id"))
	rearmed, err := model.RedeliverTaskCallbackEvent(eventID)
	params := map[string]any{"event_id": eventID, "success": err == nil && rearmed}
	defer func() {
		action := "task_callback.redeliver"
		auditInfo := &model.AuditRequestInfo{
			Method: c.Request.Method, Route: c.FullPath(), Path: c.FullPath(),
			Status: c.Writer.Status(), Success: err == nil && rearmed,
		}
		model.RecordOperationAuditLog(c.GetInt("id"), c.GetInt("role"), auditContentEN(action, params), c.ClientIP(), action, params, auditOperatorInfo(c), auditInfo, c)
		markAuditLogged(c)
	}()
	if err != nil {
		common.ApiErrorMsg(c, "failed to redeliver task callback")
		return
	}
	if !rearmed {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "only failed or cancelled callbacks can be redelivered"})
		return
	}
	service.WakeTaskCallbackDispatcher()
	common.ApiSuccess(c, gin.H{"event_id": eventID})
}
