package model

import (
	"errors"

	"github.com/QuantumNous/new-api/common"

	"gorm.io/gorm"
)

// Task callback event delivery states.
const (
	TaskCallbackStatePending    = "pending"
	TaskCallbackStateDelivering = "delivering"
	TaskCallbackStateSucceeded  = "succeeded"
	TaskCallbackStateFailed     = "failed"
	TaskCallbackStateSuperseded = "superseded"
	TaskCallbackStateCancelled  = "cancelled"
)

// TaskCallbackMaxPayloadBytes keeps snapshots inside a MySQL TEXT column.
const TaskCallbackMaxPayloadBytes = 60000

// TaskCallbackSubscription links one task to the customer's callback_url
// and, when the upstream fast path is used, to kt-api's own upstream
// callback token. One row per task.
type TaskCallbackSubscription struct {
	ID     int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskID string `json:"task_id" gorm:"column:task_id;type:varchar(191);not null;uniqueIndex:uk_task_callback_subscriptions_task_id"`
	UserID int    `json:"user_id" gorm:"column:user_id;not null;index:idx_task_callback_subscriptions_user_id"`
	// CallbackURL is the customer's URL (plaintext in this stage). Empty for a
	// subscription that only carries the upstream token.
	CallbackURL     string `json:"-" gorm:"column:callback_url;type:text"`
	CallbackURLHash string `json:"callback_url_hash" gorm:"column:callback_url_hash;type:varchar(64)"`
	CallbackHost    string `json:"callback_host" gorm:"column:callback_host;type:varchar(255)"`
	// UpstreamTokenHash is the SHA-256 of the per-task upstream callback
	// token; NULL when the fast path is not used (NULLs are not unique).
	UpstreamTokenHash      *string `json:"-" gorm:"column:upstream_token_hash;type:varchar(64);uniqueIndex:uk_task_callback_subscriptions_upstream_token"`
	UpstreamTokenExpiresAt int64   `json:"-" gorm:"column:upstream_token_expires_at;bigint"`
	LastStatus             string  `json:"last_status" gorm:"column:last_status;type:varchar(20)"`
	LastSeq                int     `json:"last_seq" gorm:"column:last_seq"`
	CreatedAt              int64   `json:"created_at" gorm:"column:created_at;bigint"`
	UpdatedAt              int64   `json:"updated_at" gorm:"column:updated_at;bigint;index:idx_task_callback_subscriptions_updated_at"`
}

func (TaskCallbackSubscription) TableName() string { return "task_callback_subscriptions" }

// TaskCallbackEvent is one outbox row: one customer callback per visible
// task status change, delivered at least once.
type TaskCallbackEvent struct {
	ID             int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	EventID        string `json:"event_id" gorm:"column:event_id;type:varchar(64);not null;uniqueIndex:uk_task_callback_events_event_id"`
	SubscriptionID int64  `json:"subscription_id" gorm:"column:subscription_id;not null;index:idx_task_callback_events_subscription"`
	TaskID         string `json:"task_id" gorm:"column:task_id;type:varchar(191);not null;uniqueIndex:uk_task_callback_events_task_status,priority:1;index:idx_task_callback_events_task_seq,priority:1"`
	Status         string `json:"status" gorm:"column:status;type:varchar(20);not null;uniqueIndex:uk_task_callback_events_task_status,priority:2"`
	UserID         int    `json:"user_id" gorm:"column:user_id;not null;index:idx_task_callback_events_user_created,priority:1"`
	Seq            int    `json:"seq" gorm:"column:seq;not null;index:idx_task_callback_events_task_seq,priority:2"`
	Payload        string `json:"-" gorm:"column:payload;type:text"`
	PayloadSHA256  string `json:"payload_sha256" gorm:"column:payload_sha256;type:varchar(64)"`
	State          string `json:"state" gorm:"column:state;type:varchar(16);not null;index:idx_task_callback_events_state_next,priority:1;index:idx_task_callback_events_state_lease,priority:1"`
	Attempts       int    `json:"attempts" gorm:"column:attempts;not null"`
	MaxAttempts    int    `json:"max_attempts" gorm:"column:max_attempts;not null"`
	NextAttemptAt  int64  `json:"next_attempt_at" gorm:"column:next_attempt_at;bigint;index:idx_task_callback_events_state_next,priority:2"`
	LeaseOwner     string `json:"-" gorm:"column:lease_owner;type:varchar(64)"`
	LeaseUntil     int64  `json:"-" gorm:"column:lease_until;bigint;index:idx_task_callback_events_state_lease,priority:2"`
	LastHTTPStatus int    `json:"last_http_status" gorm:"column:last_http_status"`
	LastError      string `json:"last_error" gorm:"column:last_error;type:varchar(512)"`
	LastLatencyMs  int    `json:"last_latency_ms" gorm:"column:last_latency_ms"`
	DeliveredAt    int64  `json:"delivered_at" gorm:"column:delivered_at;bigint"`
	CreatedAt      int64  `json:"created_at" gorm:"column:created_at;bigint;index:idx_task_callback_events_user_created,priority:2"`
	UpdatedAt      int64  `json:"updated_at" gorm:"column:updated_at;bigint"`
}

func (TaskCallbackEvent) TableName() string { return "task_callback_events" }

// TaskCallbackSupersedableStatuses are the public statuses whose pending
// events a newer event of the same task supersedes.
var TaskCallbackSupersedableStatuses = []string{"queued", "running"}

// TaskCallbackOpenStates are the states that still lead to a delivery.
var TaskCallbackOpenStates = []string{TaskCallbackStatePending, TaskCallbackStateDelivering}

// TaskCallbackFinishedStates are terminal delivery states.
var TaskCallbackFinishedStates = []string{
	TaskCallbackStateSucceeded,
	TaskCallbackStateFailed,
	TaskCallbackStateSuperseded,
	TaskCallbackStateCancelled,
}

func CreateTaskCallbackSubscription(subscription *TaskCallbackSubscription) error {
	now := common.GetTimestamp()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	return DB.Create(subscription).Error
}

func GetTaskCallbackSubscriptionByTaskID(taskID string) (*TaskCallbackSubscription, bool, error) {
	return getTaskCallbackSubscription(DB, "task_id = ?", taskID)
}

func GetTaskCallbackSubscriptionByTaskIDTx(tx *gorm.DB, taskID string) (*TaskCallbackSubscription, bool, error) {
	return getTaskCallbackSubscription(tx, "task_id = ?", taskID)
}

// GetTaskCallbackSubscriptionByID loads a subscription by primary key.
func GetTaskCallbackSubscriptionByID(id int64) (*TaskCallbackSubscription, bool, error) {
	return getTaskCallbackSubscription(DB, "id = ?", id)
}

func GetTaskCallbackSubscriptionByUpstreamTokenHash(tokenHash string) (*TaskCallbackSubscription, bool, error) {
	if tokenHash == "" {
		return nil, false, nil
	}
	return getTaskCallbackSubscription(DB, "upstream_token_hash = ?", tokenHash)
}

func getTaskCallbackSubscription(db *gorm.DB, query string, value any) (*TaskCallbackSubscription, bool, error) {
	var subscription TaskCallbackSubscription
	err := db.Where(query, value).Take(&subscription).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &subscription, true, nil
}

// ListOpenTaskCallbackSubscriptions pages customer subscriptions whose last
// enqueued status is not terminal, in id order after afterID.
func ListOpenTaskCallbackSubscriptions(afterID int64, terminalStatuses []string, limit int) ([]*TaskCallbackSubscription, error) {
	var subscriptions []*TaskCallbackSubscription
	query := DB.Where("id > ?", afterID).Where("callback_url <> ?", "")
	if len(terminalStatuses) > 0 {
		query = query.Where("last_status NOT IN ?", terminalStatuses)
	}
	err := query.Order("id").Limit(limit).Find(&subscriptions).Error
	return subscriptions, err
}

// HasOpenTaskCallbackSubscriptions reports whether any customer subscription
// still awaits a terminal status (gates the reconciliation system task).
func HasOpenTaskCallbackSubscriptions(terminalStatuses []string) bool {
	var subscription TaskCallbackSubscription
	query := DB.Select("id").Where("callback_url <> ?", "")
	if len(terminalStatuses) > 0 {
		query = query.Where("last_status NOT IN ?", terminalStatuses)
	}
	err := query.Take(&subscription).Error
	return err == nil
}

// RecordTaskCallbackSubscriptionStatus advances last_status without an event
// (used for statuses that are not delivered, e.g. cancelled).
func RecordTaskCallbackSubscriptionStatus(tx *gorm.DB, subscriptionID int64, status string) error {
	return tx.Model(&TaskCallbackSubscription{}).Where("id = ?", subscriptionID).Updates(map[string]any{
		"last_status": status,
		"updated_at":  common.GetTimestamp(),
	}).Error
}

// InsertTaskCallbackEventTx inserts the event unless one already exists for
// the same (task_id, status) and advances the subscription. It reports
// whether a new event was inserted.
func InsertTaskCallbackEventTx(tx *gorm.DB, subscription *TaskCallbackSubscription, event *TaskCallbackEvent) (bool, error) {
	now := common.GetTimestamp()
	event.SubscriptionID = subscription.ID
	event.TaskID = subscription.TaskID
	event.UserID = subscription.UserID
	event.Seq = subscription.LastSeq + 1
	if event.State == "" {
		event.State = TaskCallbackStatePending
	}
	if event.NextAttemptAt == 0 {
		event.NextAttemptAt = now
	}
	event.CreatedAt = now
	event.UpdatedAt = now
	exists, err := taskCallbackEventExists(tx, event.TaskID, event.Status)
	if err != nil || exists {
		return false, err
	}
	// The insert runs in a savepoint: a concurrent writer that won the
	// (task_id, status) unique key only fails this statement, and the
	// re-check turns that race into "already enqueued" on every database
	// (PostgreSQL would otherwise abort the whole transaction).
	if insertErr := tx.Transaction(func(savepoint *gorm.DB) error {
		return savepoint.Create(event).Error
	}); insertErr != nil {
		if exists, checkErr := taskCallbackEventExists(tx, event.TaskID, event.Status); checkErr == nil && exists {
			return false, nil
		}
		return false, insertErr
	}
	update := tx.Model(&TaskCallbackSubscription{}).Where("id = ?", subscription.ID).Updates(map[string]any{
		"last_status": event.Status,
		"last_seq":    event.Seq,
		"updated_at":  now,
	})
	if update.Error != nil {
		return false, update.Error
	}
	// A newer status makes undelivered queued/running events stale: they
	// are superseded instead of being sent late. Terminal statuses are
	// never superseded.
	if _, err := supersedeTaskCallbackEvents(tx, event.TaskID, event.Seq, TaskCallbackSupersedableStatuses); err != nil {
		return false, err
	}
	subscription.LastStatus = event.Status
	subscription.LastSeq = event.Seq
	return true, nil
}

func taskCallbackEventExists(tx *gorm.DB, taskID, status string) (bool, error) {
	var count int64
	err := tx.Model(&TaskCallbackEvent{}).Where("task_id = ? AND status = ?", taskID, status).Count(&count).Error
	return count > 0, err
}

// SupersedeTaskCallbackEvents marks still-pending events of the task with a
// smaller seq and one of the given (non-terminal) statuses as superseded.
func SupersedeTaskCallbackEvents(taskID string, beforeSeq int, statuses []string) (int64, error) {
	return supersedeTaskCallbackEvents(DB, taskID, beforeSeq, statuses)
}

func supersedeTaskCallbackEvents(db *gorm.DB, taskID string, beforeSeq int, statuses []string) (int64, error) {
	result := db.Model(&TaskCallbackEvent{}).
		Where("task_id = ? AND seq < ? AND state = ? AND status IN ?", taskID, beforeSeq, TaskCallbackStatePending, statuses).
		Updates(map[string]any{"state": TaskCallbackStateSuperseded, "updated_at": common.GetTimestamp()})
	return result.RowsAffected, result.Error
}

// HasEarlierOpenTaskCallbackEvent reports whether the task has an event with
// a smaller seq that is still pending or being delivered.
func HasEarlierOpenTaskCallbackEvent(taskID string, seq int) (bool, error) {
	var count int64
	err := DB.Model(&TaskCallbackEvent{}).
		Where("task_id = ? AND seq < ? AND state IN ?", taskID, seq, TaskCallbackOpenStates).
		Count(&count).Error
	return count > 0, err
}

// ListDueTaskCallbackEvents returns pending events whose next attempt is due.
func ListDueTaskCallbackEvents(now int64, limit int) ([]*TaskCallbackEvent, error) {
	var events []*TaskCallbackEvent
	err := DB.Where("state = ? AND next_attempt_at <= ?", TaskCallbackStatePending, now).
		Order("next_attempt_at").Order("id").
		Limit(limit).
		Find(&events).Error
	return events, err
}

// ClaimTaskCallbackEvent leases a due pending event for one delivery
// attempt. Only one caller across all instances can win (conditional UPDATE).
func ClaimTaskCallbackEvent(id int64, owner string, now, leaseUntil int64) (bool, error) {
	result := DB.Model(&TaskCallbackEvent{}).
		Where("id = ? AND state = ? AND next_attempt_at <= ?", id, TaskCallbackStatePending, now).
		Updates(map[string]any{
			"state":       TaskCallbackStateDelivering,
			"lease_owner": owner,
			"lease_until": leaseUntil,
			"attempts":    gorm.Expr("attempts + 1"),
			"updated_at":  now,
		})
	return result.RowsAffected == 1, result.Error
}

// GetTaskCallbackEvent loads one event by primary key.
func GetTaskCallbackEvent(id int64) (*TaskCallbackEvent, error) {
	var event TaskCallbackEvent
	if err := DB.Take(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// TaskCallbackAttemptResult is written back after one delivery attempt.
type TaskCallbackAttemptResult struct {
	State         string
	NextAttemptAt int64
	HTTPStatus    int
	Error         string
	LatencyMs     int
	DeliveredAt   int64
	// MaxAttempts, when positive, raises the attempt budget (extended retry).
	MaxAttempts int
}

// FinishTaskCallbackAttempt records an attempt result if the caller still
// holds the lease. It reports whether the write applied.
func FinishTaskCallbackAttempt(id int64, owner string, result TaskCallbackAttemptResult) (bool, error) {
	if len(result.Error) > 512 {
		result.Error = result.Error[:512]
	}
	updates := map[string]any{
		"state":            result.State,
		"lease_owner":      "",
		"lease_until":      0,
		"last_http_status": result.HTTPStatus,
		"last_error":       result.Error,
		"last_latency_ms":  result.LatencyMs,
		"updated_at":       common.GetTimestamp(),
	}
	if result.NextAttemptAt > 0 {
		updates["next_attempt_at"] = result.NextAttemptAt
	}
	if result.DeliveredAt > 0 {
		updates["delivered_at"] = result.DeliveredAt
	}
	if result.MaxAttempts > 0 {
		updates["max_attempts"] = result.MaxAttempts
	}
	write := DB.Model(&TaskCallbackEvent{}).
		Where("id = ? AND lease_owner = ? AND state = ?", id, owner, TaskCallbackStateDelivering).
		Updates(updates)
	return write.RowsAffected == 1, write.Error
}

// ReclaimExpiredTaskCallbackLeases returns events whose delivery lease
// expired (e.g. the instance died mid-attempt) to pending, or fails them when
// their attempt budget is spent. It returns (reclaimed, failed).
func ReclaimExpiredTaskCallbackLeases(now int64) (int64, int64, error) {
	failed := DB.Model(&TaskCallbackEvent{}).
		Where("state = ? AND lease_until < ? AND attempts >= max_attempts", TaskCallbackStateDelivering, now).
		Updates(map[string]any{
			"state":       TaskCallbackStateFailed,
			"lease_owner": "",
			"lease_until": 0,
			"last_error":  "delivery lease expired",
			"updated_at":  now,
		})
	if failed.Error != nil {
		return 0, 0, failed.Error
	}
	reclaimed := DB.Model(&TaskCallbackEvent{}).
		Where("state = ? AND lease_until < ?", TaskCallbackStateDelivering, now).
		Updates(map[string]any{
			"state":           TaskCallbackStatePending,
			"lease_owner":     "",
			"lease_until":     0,
			"next_attempt_at": now,
			"updated_at":      now,
		})
	if reclaimed.Error != nil {
		return 0, failed.RowsAffected, reclaimed.Error
	}
	return reclaimed.RowsAffected, failed.RowsAffected, nil
}

// ReleaseTaskCallbackLeases returns this owner's in-flight events to pending
// (best effort on shutdown). Attempts already counted are kept.
func ReleaseTaskCallbackLeases(owner string) error {
	now := common.GetTimestamp()
	return DB.Model(&TaskCallbackEvent{}).
		Where("state = ? AND lease_owner = ?", TaskCallbackStateDelivering, owner).
		Updates(map[string]any{
			"state":           TaskCallbackStatePending,
			"lease_owner":     "",
			"lease_until":     0,
			"next_attempt_at": now,
			"updated_at":      now,
		}).Error
}

// RedeliverTaskCallbackEvent re-arms a failed or cancelled event with a
// fresh attempt budget and the same event_id.
func RedeliverTaskCallbackEvent(eventID string) (bool, error) {
	now := common.GetTimestamp()
	result := DB.Model(&TaskCallbackEvent{}).
		Where("event_id = ? AND state IN ?", eventID, []string{TaskCallbackStateFailed, TaskCallbackStateCancelled}).
		Updates(map[string]any{
			"state":           TaskCallbackStatePending,
			"attempts":        0,
			"next_attempt_at": now,
			"last_error":      "",
			"updated_at":      now,
		})
	return result.RowsAffected == 1, result.Error
}

// CancelPendingTaskCallbackEvents cancels every pending event (e.g. when
// delivery is switched off for good by an administrator).
func CancelPendingTaskCallbackEvents() (int64, error) {
	result := DB.Model(&TaskCallbackEvent{}).
		Where("state = ?", TaskCallbackStatePending).
		Updates(map[string]any{"state": TaskCallbackStateCancelled, "updated_at": common.GetTimestamp()})
	return result.RowsAffected, result.Error
}

// DeleteFinishedTaskCallbackEvents deletes up to limit finished events last
// updated before cutoff.
func DeleteFinishedTaskCallbackEvents(cutoff int64, limit int) (int64, error) {
	var ids []int64
	if err := DB.Model(&TaskCallbackEvent{}).
		Where("state IN ? AND updated_at < ?", TaskCallbackFinishedStates, cutoff).
		Order("id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := DB.Where("id IN ?", ids).Delete(&TaskCallbackEvent{})
	return result.RowsAffected, result.Error
}

// DeleteStaleTaskCallbackSubscriptions deletes up to limit subscriptions last
// updated before cutoff that have no remaining events.
func DeleteStaleTaskCallbackSubscriptions(cutoff int64, limit int) (int64, error) {
	var ids []int64
	if err := DB.Model(&TaskCallbackSubscription{}).
		Where("updated_at < ?", cutoff).
		Where("NOT EXISTS (SELECT 1 FROM task_callback_events e WHERE e.subscription_id = task_callback_subscriptions.id)").
		Order("id").Limit(limit).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}
	result := DB.Where("id IN ?", ids).Delete(&TaskCallbackSubscription{})
	return result.RowsAffected, result.Error
}

// ListTaskCallbackEventsByTask returns a task's events in seq order.
func ListTaskCallbackEventsByTask(taskID string) ([]*TaskCallbackEvent, error) {
	var events []*TaskCallbackEvent
	err := DB.Where("task_id = ?", taskID).Order("seq").Order("id").Find(&events).Error
	return events, err
}

// TaskCallbackEventQuery filters the admin event list.
type TaskCallbackEventQuery struct {
	TaskID         string
	UserID         int
	State          string
	Status         string
	StartTimestamp int64
	EndTimestamp   int64
}

func (q TaskCallbackEventQuery) apply(db *gorm.DB) *gorm.DB {
	if q.TaskID != "" {
		db = db.Where("task_id = ?", q.TaskID)
	}
	if q.UserID > 0 {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.State != "" {
		db = db.Where("state = ?", q.State)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.StartTimestamp > 0 {
		db = db.Where("created_at >= ?", q.StartTimestamp)
	}
	if q.EndTimestamp > 0 {
		db = db.Where("created_at <= ?", q.EndTimestamp)
	}
	return db
}

// ListTaskCallbackEvents returns one page of events (newest first) and the
// total count for the filter.
func ListTaskCallbackEvents(query TaskCallbackEventQuery, offset, limit int) ([]*TaskCallbackEvent, int64, error) {
	var total int64
	if err := query.apply(DB.Model(&TaskCallbackEvent{})).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var events []*TaskCallbackEvent
	err := query.apply(DB.Model(&TaskCallbackEvent{})).Order("id DESC").Offset(offset).Limit(limit).Find(&events).Error
	return events, total, err
}

// ListTaskCallbackSubscriptionsByIDs loads subscriptions keyed by id.
func ListTaskCallbackSubscriptionsByIDs(ids []int64) (map[int64]*TaskCallbackSubscription, error) {
	result := make(map[int64]*TaskCallbackSubscription, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	var subscriptions []*TaskCallbackSubscription
	if err := DB.Where("id IN ?", ids).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	for _, subscription := range subscriptions {
		result[subscription.ID] = subscription
	}
	return result, nil
}
