package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"gorm.io/gorm"
)

// taskCallbackRenderTimeout bounds rendering a snapshot inside a write
// transaction.
const taskCallbackRenderTimeout = 5 * time.Second

// Reconciliation re-checks subscriptions whose task changed without an
// enqueue (bulk writes, observer errors, historical rows).
const (
	taskCallbackReconcileMinAge    = 30 // seconds since the task's last update
	taskCallbackReconcileBatchSize = 200
	taskCallbackReconcileMaxBatch  = 25
)

// taskCallbackTerminalStatuses close a subscription: no later status can be
// delivered after them.
var taskCallbackTerminalStatuses = []string{
	ArkTaskStatusSucceeded,
	ArkTaskStatusFailed,
	ArkTaskStatusExpired,
	ArkTaskStatusCancelled,
}

// taskCallbackDelivered reports whether a public status triggers a customer
// callback (supplier doc: queued, running, succeeded, failed, expired).
// cancelled is recorded but not delivered.
func taskCallbackDelivered(status string) bool {
	switch status {
	case ArkTaskStatusQueued, ArkTaskStatusRunning, ArkTaskStatusSucceeded, ArkTaskStatusFailed, ArkTaskStatusExpired:
		return true
	}
	return false
}

// taskCallbackDocumentedRetries reports whether the status gets the
// documented retries ("成功或失败回调在发送未获确认时最多重试三次").
func taskCallbackDocumentedRetries(status string) bool {
	return status == ArkTaskStatusSucceeded || status == ArkTaskStatusFailed
}

func taskCallbackMaxAttempts(status string) int {
	if taskCallbackDocumentedRetries(status) {
		return system_setting.GetTaskCallbackSetting().TerminalAttempts()
	}
	return 1
}

func newTaskCallbackEventID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "evt_" + common.GetUUID()
	}
	return "evt_" + hex.EncodeToString(buf)
}

// publicTaskStatus extracts the top-level status of a rendered body.
func publicTaskStatus(body []byte) string {
	var envelope struct {
		Status string `json:"status"`
	}
	if err := common.Unmarshal(body, &envelope); err != nil {
		return ""
	}
	return envelope.Status
}

// EnqueueTaskCallbackIfChanged renders the task exactly as GET would, and
// when its public status differs from the subscription's last status and is
// deliverable, inserts one outbox event (unique per task and status) in tx.
// It reports whether an event was inserted.
func EnqueueTaskCallbackIfChanged(tx *gorm.DB, subscription *model.TaskCallbackSubscription, task *model.Task) (bool, error) {
	if subscription == nil || task == nil || subscription.CallbackURL == "" {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), taskCallbackRenderTimeout)
	defer cancel()
	body, err := RenderSeedanceTaskStatus(ctx, task)
	if err != nil {
		return false, err
	}
	status := publicTaskStatus(body)
	if status == "" || status == subscription.LastStatus {
		return false, nil
	}
	if !taskCallbackDelivered(status) {
		if status == ArkTaskStatusCancelled {
			if err := model.RecordTaskCallbackSubscriptionStatus(tx, subscription.ID, status); err != nil {
				return false, err
			}
			subscription.LastStatus = status
			common.SysLog(common.LogText("task_callback event=not_delivered task_id=%s status=%s reason=status_not_documented", task.TaskID, status))
		}
		return false, nil
	}
	if len(body) > model.TaskCallbackMaxPayloadBytes {
		common.SysError(common.LogText("task_callback event=enqueue_skipped task_id=%s status=%s reason=payload_too_large bytes=%d", task.TaskID, status, len(body)))
		return false, nil
	}
	sum := sha256.Sum256(body)
	event := &model.TaskCallbackEvent{
		EventID:       newTaskCallbackEventID(),
		Status:        status,
		Payload:       string(body),
		PayloadSHA256: hex.EncodeToString(sum[:]),
		MaxAttempts:   taskCallbackMaxAttempts(status),
	}
	inserted, err := model.InsertTaskCallbackEventTx(tx, subscription, event)
	if err != nil {
		return false, err
	}
	if inserted {
		common.SysLog(common.LogText("task_callback event=enqueued event_id=%s task_id=%s status=%s seq=%d max_attempts=%d", event.EventID, task.TaskID, status, event.Seq, event.MaxAttempts))
	}
	return inserted, nil
}

// taskCallbackObserver enqueues callback events in the same transaction as
// the task's CAS update.
type taskCallbackObserver struct{}

func (taskCallbackObserver) ObservesTask(task *model.Task) bool {
	return system_setting.GetTaskCallbackSetting().Enabled && IsArkVideoTask(task)
}

func (taskCallbackObserver) OnTaskTransition(tx *gorm.DB, task *model.Task) (bool, error) {
	subscription, found, err := model.GetTaskCallbackSubscriptionByTaskIDTx(tx, task.TaskID)
	if err != nil || !found {
		return false, err
	}
	return EnqueueTaskCallbackIfChanged(tx, subscription, task)
}

func (taskCallbackObserver) AfterTaskTransition(*model.Task) {
	WakeTaskCallbackDispatcher()
}

func init() {
	model.SetTaskTransitionObserver(taskCallbackObserver{})
}

// taskCallbackWake nudges the local dispatcher after a commit; dispatchers
// also poll, so a dropped nudge only delays delivery until the next tick.
var taskCallbackWake = make(chan struct{}, 1)

// WakeTaskCallbackDispatcher asks the local dispatcher to scan now.
func WakeTaskCallbackDispatcher() {
	select {
	case taskCallbackWake <- struct{}{}:
	default:
	}
}

// HasOpenTaskCallbackSubscriptions reports whether reconciliation has work.
func HasOpenTaskCallbackSubscriptions() bool {
	return model.HasOpenTaskCallbackSubscriptions(taskCallbackTerminalStatuses)
}

// TaskCallbackReconcileSummary is the result of one reconciliation pass.
type TaskCallbackReconcileSummary struct {
	Scanned  int `json:"scanned"`
	Enqueued int `json:"enqueued"`
	Errors   int `json:"errors"`
}

// ReconcileTaskCallbacks enqueues events for open subscriptions whose task
// status moved without an enqueue (bulk writes, observer failures, rows
// written before the subscription existed). Tasks updated in the last 30
// seconds are left to the transactional path.
func ReconcileTaskCallbacks(ctx context.Context) (TaskCallbackReconcileSummary, error) {
	summary := TaskCallbackReconcileSummary{}
	if !system_setting.GetTaskCallbackSetting().Enabled {
		return summary, nil
	}
	cutoff := common.GetTimestamp() - taskCallbackReconcileMinAge
	afterID := int64(0)
	for batch := 0; batch < taskCallbackReconcileMaxBatch; batch++ {
		if err := ctx.Err(); err != nil {
			return summary, err
		}
		subscriptions, err := model.ListOpenTaskCallbackSubscriptions(afterID, taskCallbackTerminalStatuses, taskCallbackReconcileBatchSize)
		if err != nil {
			return summary, err
		}
		for _, subscription := range subscriptions {
			afterID = subscription.ID
			summary.Scanned++
			enqueued, err := reconcileTaskCallbackSubscription(subscription, cutoff)
			if err != nil {
				summary.Errors++
				common.SysError(common.LogText("task_callback event=reconcile_failed task_id=%s error=%q", subscription.TaskID, err.Error()))
				continue
			}
			if enqueued {
				summary.Enqueued++
			}
		}
		if len(subscriptions) < taskCallbackReconcileBatchSize {
			break
		}
	}
	if summary.Enqueued > 0 {
		WakeTaskCallbackDispatcher()
	}
	return summary, nil
}

func reconcileTaskCallbackSubscription(subscription *model.TaskCallbackSubscription, cutoff int64) (bool, error) {
	task, exists, err := model.GetByOnlyTaskId(subscription.TaskID)
	if err != nil {
		return false, err
	}
	if !exists || task == nil {
		return false, nil
	}
	if task.UpdatedAt > cutoff {
		return false, nil
	}
	enqueued := false
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		current, found, err := model.GetTaskCallbackSubscriptionByTaskIDTx(tx, subscription.TaskID)
		if err != nil || !found {
			return err
		}
		enqueued, err = EnqueueTaskCallbackIfChanged(tx, current, task)
		return err
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return enqueued, err
}
