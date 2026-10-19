package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SeedanceCallbackTargetFromContext returns the validated customer
// callback_url stored by the submit middleware, if any.
func SeedanceCallbackTargetFromContext(c *gin.Context) *TaskCallbackTarget {
	value, ok := c.Get(ContextKeySeedanceCallbackTarget)
	if !ok {
		return nil
	}
	target, _ := value.(*TaskCallbackTarget)
	return target
}

// RecordTaskCallbackSubscription stores the customer callback_url and/or
// the upstream fast-path token hash for a durable task. A failure is logged
// and does not fail the submission: the task is already persisted and
// billed, and remains queryable via GET.
func RecordTaskCallbackSubscription(c *gin.Context, task *model.Task, target *TaskCallbackTarget, info *relaycommon.RelayInfo) *model.TaskCallbackSubscription {
	if task == nil {
		return nil
	}
	upstreamToken := taskUpstreamCallbackTokenFromContext(c, info)
	if target == nil && upstreamToken == "" {
		return nil
	}
	subscription := &model.TaskCallbackSubscription{
		TaskID: task.TaskID,
		UserID: task.UserId,
	}
	if target != nil {
		subscription.CallbackURL = target.URL
		subscription.CallbackURLHash = target.URLHash
		subscription.CallbackHost = target.Host
	}
	if upstreamToken != "" {
		hash := TaskUpstreamCallbackTokenHash(upstreamToken)
		subscription.UpstreamTokenHash = &hash
		subscription.UpstreamTokenExpiresAt = ArkTaskDeadline(task) + taskUpstreamTokenTTLAfterDeadline
	}
	if err := model.CreateTaskCallbackSubscription(subscription); err != nil {
		logger.LogError(c, common.LogText("task_callback event=subscribe_failed task_id=%s target=%s error=%q", task.TaskID, target.LogRef(), err.Error()))
		return nil
	}
	logger.LogInfo(c, common.LogText("task_callback event=subscribed task_id=%s target=%s upstream_fast_path=%t", task.TaskID, target.LogRef(), upstreamToken != ""))
	if subscription.CallbackURL != "" {
		if _, err := EnqueueInitialTaskCallback(task.TaskID); err != nil {
			// Reconciliation retries the enqueue; the submission still succeeds.
			logger.LogError(c, common.LogText("task_callback event=initial_enqueue_failed task_id=%s error=%q", task.TaskID, err.Error()))
		}
	}
	return subscription
}

// EnqueueInitialTaskCallback delivers the first status a customer sees right
// after a successful submit (none -> queued). It renders the task exactly as
// GET does at that moment, so the body is byte-identical to GET. It does
// nothing once the subscription already has a status: if the poller or the
// upstream fast path observed a later status first, that status was (or will
// be) enqueued instead and queued is never sent after it.
func EnqueueInitialTaskCallback(taskID string) (bool, error) {
	if !system_setting.GetTaskCallbackSetting().Enabled {
		return false, nil
	}
	enqueued := false
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		subscription, found, err := model.GetTaskCallbackSubscriptionByTaskIDTx(tx, taskID)
		if err != nil || !found || subscription.CallbackURL == "" || subscription.LastStatus != "" {
			return err
		}
		var task model.Task
		if err := tx.Where("task_id = ?", taskID).First(&task).Error; err != nil {
			return err
		}
		if !IsArkVideoTask(&task) {
			return nil
		}
		enqueued, err = EnqueueTaskCallbackIfChanged(tx, subscription, &task)
		return err
	})
	if err != nil {
		return false, err
	}
	if enqueued {
		WakeTaskCallbackDispatcher()
	}
	return enqueued, nil
}
