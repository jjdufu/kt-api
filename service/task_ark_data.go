package service

import (
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

// Ark-compatible (Volcano Ark / Seedance) video task statuses as presented on
// GET /api/v3/contents/generations/tasks/{id} and in task callbacks.
const (
	ArkTaskStatusQueued    = "queued"
	ArkTaskStatusRunning   = "running"
	ArkTaskStatusSucceeded = "succeeded"
	ArkTaskStatusFailed    = "failed"
	ArkTaskStatusExpired   = "expired"
	ArkTaskStatusCancelled = "cancelled"
)

// Error codes kt-api writes into the Ark-shaped task body when it decides a
// terminal state itself instead of observing it upstream.
const (
	arkErrorCodeTaskExpired         = "task_expired"
	arkErrorCodeUpstreamNotFound    = "upstream_task_not_found"
	arkErrorCodeUpstreamQueryFailed = "upstream_query_failed"
	arkErrorCodeTaskFailed          = "task_failed"
)

// arkVideoTaskPluginKey is the task plugin that serves Ark-compatible Seedance
// video tasks (plugins/tasks/doubao).
const arkVideoTaskPluginKey = "doubao"

// IsArkVideoTask reports whether the task is a Seedance video task whose
// persisted Data is an Ark-shaped task object. Image tasks served by the same
// plugin are synchronous and excluded.
func IsArkVideoTask(task *model.Task) bool {
	if task == nil || task.PrivateData.Execution == nil || task.PrivateData.Execution.TaskPlugin == nil {
		return false
	}
	if task.PrivateData.Execution.TaskPlugin.Key != arkVideoTaskPluginKey {
		return false
	}
	return task.Action == "text_to_video" || task.Action == "image_to_video"
}

// arkTaskDataObject decodes the persisted task snapshot as a JSON object. Any
// other shape (absent, array, scalar, invalid) yields an empty object.
func arkTaskDataObject(task *model.Task) map[string]any {
	object := map[string]any{}
	if task == nil || len(task.Data) == 0 {
		return object
	}
	var decoded any
	if err := common.Unmarshal(task.Data, &decoded); err != nil {
		return object
	}
	if typed, ok := decoded.(map[string]any); ok {
		return typed
	}
	return object
}

// arkTerminalTaskData builds the Ark-shaped terminal body kt-api persists when
// it decides (or must synthesize) a terminal state: the last upstream object
// is kept, and status, error and updated_at are replaced. The private
// upstream id is kept in id; public views rewrite it to the public task id.
func arkTerminalTaskData(task *model.Task, status, code, message string, now int64) []byte {
	object := arkTaskDataObject(task)
	// Any previous error object is replaced; the caller supplies code/message.
	delete(object, "error")
	if id, _ := object["id"].(string); id == "" {
		if upstreamID := task.GetUpstreamTaskID(); upstreamID != "" {
			object["id"] = upstreamID
		}
	}
	if properties := task.Properties; properties.UpstreamModelName != "" {
		if model, _ := object["model"].(string); model == "" {
			object["model"] = properties.UpstreamModelName
		}
	}
	if _, ok := object["created_at"]; !ok && task.SubmitTime > 0 {
		object["created_at"] = task.SubmitTime
	}
	if now <= 0 {
		now = time.Now().Unix()
	}
	object["status"] = status
	object["error"] = map[string]any{"code": code, "message": message}
	object["updated_at"] = now
	encoded, err := common.Marshal(object)
	if err != nil {
		return task.Data
	}
	return encoded
}

// arkLocalFailureCode maps the fixed fail reasons used by the polling loop to
// the error code written into the Ark-shaped body.
func arkLocalFailureCode(reason string) string {
	switch reason {
	case "Upstream task not found":
		return arkErrorCodeUpstreamNotFound
	case "Failed to query the upstream task status":
		return arkErrorCodeUpstreamQueryFailed
	default:
		return arkErrorCodeTaskFailed
	}
}

// arkStatuslessFailureCode inspects a failed poll response: when the body is a
// JSON object without a status field it returns the vendor error code (or a
// generic one) so the caller can persist an Ark-shaped failed task body.
func arkStatuslessFailureCode(body []byte) (string, bool) {
	var decoded any
	if err := common.Unmarshal(body, &decoded); err != nil {
		return arkErrorCodeTaskFailed, true
	}
	object, ok := decoded.(map[string]any)
	if !ok {
		return arkErrorCodeTaskFailed, true
	}
	if _, hasStatus := object["status"]; hasStatus {
		return "", false
	}
	if vendorError, ok := object["error"].(map[string]any); ok {
		if code, _ := vendorError["code"].(string); code != "" {
			return code, true
		}
	}
	return arkErrorCodeTaskFailed, true
}

// Seedance execution_expires_after bounds (seconds), from the supplier and Ark
// documents: default 172800 (48h), range [3600, 259200].
const (
	ArkDefaultExecutionExpiresAfter int64 = 172800
	ArkMinExecutionExpiresAfter     int64 = 3600
	ArkMaxExecutionExpiresAfter     int64 = 259200
	// ArkExpiryGraceSeconds lets the upstream report expired itself before
	// kt-api decides locally.
	ArkExpiryGraceSeconds int64 = 1800
)

// NormalizeArkExecutionExpiresAfter returns the requested value when it is in
// the documented range, otherwise zero (the vendor default applies).
func NormalizeArkExecutionExpiresAfter(value int64) int64 {
	if value < ArkMinExecutionExpiresAfter || value > ArkMaxExecutionExpiresAfter {
		return 0
	}
	return value
}

// arkTaskExecutionExpiresAfter resolves the task's execution_expires_after:
// the value recorded at submit, else the one echoed by the upstream snapshot,
// else the documented default.
func arkTaskExecutionExpiresAfter(task *model.Task) int64 {
	if value := NormalizeArkExecutionExpiresAfter(task.PrivateData.ExecutionExpiresAfter); value > 0 {
		return value
	}
	if number, ok := arkTaskDataObject(task)["execution_expires_after"].(float64); ok {
		if value := NormalizeArkExecutionExpiresAfter(int64(number)); value > 0 {
			return value
		}
	}
	return ArkDefaultExecutionExpiresAfter
}

// ArkTaskDeadline is the Unix second after which kt-api may time the task out
// locally: submit time + execution_expires_after + grace.
func ArkTaskDeadline(task *model.Task) int64 {
	start := task.SubmitTime
	if start <= 0 {
		start = task.CreatedAt
	}
	return start + arkTaskExecutionExpiresAfter(task) + ArkExpiryGraceSeconds
}

// arkMaxDeadlineMinutes is the smallest TASK_TIMEOUT_MINUTES that never
// defers a Seedance task (maximum execution_expires_after + grace).
func arkMaxDeadlineMinutes() int64 {
	return (ArkMaxExecutionExpiresAfter + ArkExpiryGraceSeconds + 59) / 60
}

// ContextKeySeedanceExecutionExpiresAfter carries the validated
// execution_expires_after of a Seedance create request to task persistence.
const ContextKeySeedanceExecutionExpiresAfter = "seedance_execution_expires_after"

// rewriteArkBulkFailureData gives Seedance tasks failed by a bulk update the
// Ark failed body, so GET and callbacks report status "failed" with an error
// instead of the last upstream snapshot. Billing is untouched; callback
// reconciliation enqueues the failed event.
func rewriteArkBulkFailureData(tasks map[string]*model.Task, upstreamIDs []string, reason string) {
	now := common.GetTimestamp()
	for _, upstreamID := range upstreamIDs {
		task, ok := tasks[upstreamID]
		if !ok || task == nil || !IsArkVideoTask(task) {
			continue
		}
		data := arkTerminalTaskData(task, ArkTaskStatusFailed, arkLocalFailureCode(reason), reason, now)
		if len(data) == 0 {
			continue
		}
		if err := model.TaskBulkUpdateByID([]int64{task.ID}, map[string]any{"data": common.RawMessage(data)}); err != nil {
			common.SysError(common.LogText("rewrite failed task data for %s: %v", task.TaskID, err))
		}
	}
}
