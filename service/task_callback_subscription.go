package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"

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

// RecordTaskCallbackSubscription stores the customer callback_url for a
// durable task. A failure is logged and does not fail the submission: the
// task is already persisted and billed, and remains queryable via GET.
func RecordTaskCallbackSubscription(c *gin.Context, task *model.Task, target *TaskCallbackTarget) *model.TaskCallbackSubscription {
	if task == nil || target == nil {
		return nil
	}
	subscription := &model.TaskCallbackSubscription{
		TaskID:          task.TaskID,
		UserID:          task.UserId,
		CallbackURL:     target.URL,
		CallbackURLHash: target.URLHash,
		CallbackHost:    target.Host,
	}
	if err := model.CreateTaskCallbackSubscription(subscription); err != nil {
		logger.LogError(c, common.LogText("task_callback event=subscribe_failed task_id=%s target=%s error=%q", task.TaskID, target.LogRef(), err.Error()))
		return nil
	}
	logger.LogInfo(c, common.LogText("task_callback event=subscribed task_id=%s target=%s", task.TaskID, target.LogRef()))
	return subscription
}
