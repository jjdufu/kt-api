package service

import (
	"context"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

const (
	taskCallbackCleanupBatch      = 1000
	taskCallbackCleanupMaxBatches = 50
)

// TaskCallbackCleanupSummary is recorded on the task_callback_cleanup system
// task row.
type TaskCallbackCleanupSummary struct {
	EventsDeleted        int64 `json:"events_deleted"`
	SubscriptionsDeleted int64 `json:"subscriptions_deleted"`
}

// CleanupTaskCallbacks deletes finished events and event-less subscriptions
// older than the retention period, in bounded batches.
func CleanupTaskCallbacks(ctx context.Context) (TaskCallbackCleanupSummary, error) {
	summary := TaskCallbackCleanupSummary{}
	cutoff := time.Now().Add(-system_setting.GetTaskCallbackSetting().Retention()).Unix()
	for batch := 0; batch < taskCallbackCleanupMaxBatches; batch++ {
		if ctx.Err() != nil {
			return summary, ctx.Err()
		}
		deleted, err := model.DeleteFinishedTaskCallbackEvents(cutoff, taskCallbackCleanupBatch)
		if err != nil {
			return summary, err
		}
		summary.EventsDeleted += deleted
		if deleted < taskCallbackCleanupBatch {
			break
		}
	}
	for batch := 0; batch < taskCallbackCleanupMaxBatches; batch++ {
		if ctx.Err() != nil {
			return summary, ctx.Err()
		}
		deleted, err := model.DeleteStaleTaskCallbackSubscriptions(cutoff, taskCallbackCleanupBatch)
		if err != nil {
			return summary, err
		}
		summary.SubscriptionsDeleted += deleted
		if deleted < taskCallbackCleanupBatch {
			break
		}
	}
	if summary.EventsDeleted > 0 || summary.SubscriptionsDeleted > 0 {
		logger.LogInfo(ctx, common.LogText("task_callback event=cleanup events_deleted=%d subscriptions_deleted=%d", summary.EventsDeleted, summary.SubscriptionsDeleted))
	}
	return summary, nil
}
