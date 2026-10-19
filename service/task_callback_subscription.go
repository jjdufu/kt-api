package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"

	"github.com/gin-gonic/gin"
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
	return subscription
}
