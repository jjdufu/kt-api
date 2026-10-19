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
