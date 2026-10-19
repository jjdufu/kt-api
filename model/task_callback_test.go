package model

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// sqlCaptureLogger records every statement GORM executes.
type sqlCaptureLogger struct {
	mu         sync.Mutex
	statements []string
}

func (l *sqlCaptureLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface { return l }
func (l *sqlCaptureLogger) Info(context.Context, string, ...any)             {}
func (l *sqlCaptureLogger) Warn(context.Context, string, ...any)             {}
func (l *sqlCaptureLogger) Error(context.Context, string, ...any)            {}
func (l *sqlCaptureLogger) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	sql, _ := fc()
	l.mu.Lock()
	l.statements = append(l.statements, sql)
	l.mu.Unlock()
}

func (l *sqlCaptureLogger) reset() {
	l.mu.Lock()
	l.statements = nil
	l.mu.Unlock()
}

func (l *sqlCaptureLogger) schemaChanges() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	var changes []string
	for _, statement := range l.statements {
		upper := strings.ToUpper(strings.TrimSpace(statement))
		if strings.HasPrefix(upper, "ALTER") || strings.HasPrefix(upper, "CREATE") || strings.HasPrefix(upper, "DROP") {
			changes = append(changes, statement)
		}
	}
	return changes
}

func openTaskCallbackTestDB(t *testing.T, dialect string) (*gorm.DB, *sqlCaptureLogger) {
	t.Helper()
	var driver gorm.Dialector
	switch dialect {
	case "sqlite":
		driver = sqlite.Open("file:task_callback_" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared")
	case "mysql":
		dsn := os.Getenv("TEST_MYSQL_DSN")
		if dsn == "" {
			t.Skip("TEST_MYSQL_DSN not configured")
		}
		driver = mysql.Open(dsn)
	case "postgres":
		dsn := os.Getenv("TEST_POSTGRES_DSN")
		if dsn == "" {
			t.Skip("TEST_POSTGRES_DSN not configured")
		}
		driver = postgres.Open(dsn)
	}
	capture := &sqlCaptureLogger{}
	db, err := gorm.Open(driver, &gorm.Config{Logger: capture})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(4)
	require.NoError(t, db.Migrator().DropTable(&TaskCallbackEvent{}, &TaskCallbackSubscription{}))
	previousDB, previousType := DB, common.MainDatabaseType()
	DB = db
	common.SetMainDatabaseType(common.DatabaseType(dialect))
	t.Cleanup(func() {
		require.NoError(t, db.Migrator().DropTable(&TaskCallbackEvent{}, &TaskCallbackSubscription{}))
		DB = previousDB
		common.SetMainDatabaseType(previousType)
		require.NoError(t, sqlDB.Close())
	})
	return db, capture
}

func TestTaskCallbackModelDatabaseMatrix(t *testing.T) {
	for _, dialect := range []string{"sqlite", "mysql", "postgres"} {
		t.Run(dialect, func(t *testing.T) {
			db, capture := openTaskCallbackTestDB(t, dialect)
			require.NoError(t, db.AutoMigrate(&TaskCallbackSubscription{}, &TaskCallbackEvent{}))
			var version string
			query := "SELECT version()"
			if dialect == "sqlite" {
				query = "SELECT sqlite_version()"
			}
			require.NoError(t, db.Raw(query).Scan(&version).Error)
			t.Logf("database version: %s", version)

			// A second migration (a restart) must not change the schema.
			capture.reset()
			require.NoError(t, db.AutoMigrate(&TaskCallbackSubscription{}, &TaskCallbackEvent{}))
			assert.Empty(t, capture.schemaChanges(), "repeated AutoMigrate must not alter the schema")
			for _, index := range []string{"uk_task_callback_events_task_status", "uk_task_callback_events_event_id", "idx_task_callback_events_state_next", "idx_task_callback_events_state_lease", "uk_task_callback_subscriptions_task_id", "uk_task_callback_subscriptions_upstream_token"} {
				table := "task_callback_events"
				if strings.HasPrefix(index, "uk_task_callback_subscriptions") {
					table = "task_callback_subscriptions"
				}
				assert.True(t, db.Migrator().HasIndex(table, index), "index %s", index)
			}

			// Subscriptions: unique task id, NULL upstream tokens do not collide.
			first := &TaskCallbackSubscription{TaskID: "task_a", UserID: 7, CallbackURL: "https://example.com/a"}
			second := &TaskCallbackSubscription{TaskID: "task_b", UserID: 7, CallbackURL: "https://example.com/b"}
			require.NoError(t, CreateTaskCallbackSubscription(first))
			require.NoError(t, CreateTaskCallbackSubscription(second))
			assert.Error(t, CreateTaskCallbackSubscription(&TaskCallbackSubscription{TaskID: "task_a", UserID: 7}))
			tokenHash := strings.Repeat("ab", 32)
			withToken := &TaskCallbackSubscription{TaskID: "task_c", UserID: 8, UpstreamTokenHash: &tokenHash}
			require.NoError(t, CreateTaskCallbackSubscription(withToken))
			found, ok, err := GetTaskCallbackSubscriptionByUpstreamTokenHash(tokenHash)
			require.NoError(t, err)
			require.True(t, ok)
			assert.Equal(t, "task_c", found.TaskID)
			_, ok, err = GetTaskCallbackSubscriptionByUpstreamTokenHash(strings.Repeat("cd", 32))
			require.NoError(t, err)
			assert.False(t, ok)

			// Events: one per (task, status); the duplicate is silently ignored.
			insert := func(subscription *TaskCallbackSubscription, status, eventID string, maxAttempts int) bool {
				var inserted bool
				require.NoError(t, DB.Transaction(func(tx *gorm.DB) error {
					var err error
					inserted, err = InsertTaskCallbackEventTx(tx, subscription, &TaskCallbackEvent{EventID: eventID, Status: status, Payload: `{"status":"` + status + `"}`, MaxAttempts: maxAttempts})
					return err
				}))
				return inserted
			}
			assert.True(t, insert(first, "queued", "evt_1", 1))
			assert.True(t, insert(first, "running", "evt_2", 1))
			assert.False(t, insert(first, "running", "evt_dup", 1), "same status twice is deduplicated")
			assert.True(t, insert(first, "succeeded", "evt_3", 4))
			reloaded, ok, err := GetTaskCallbackSubscriptionByTaskID("task_a")
			require.NoError(t, err)
			require.True(t, ok)
			assert.Equal(t, "succeeded", reloaded.LastStatus)
			assert.Equal(t, 3, reloaded.LastSeq)

			open, err := ListOpenTaskCallbackSubscriptions(0, []string{"succeeded", "failed", "expired"}, 10)
			require.NoError(t, err)
			require.Len(t, open, 1)
			assert.Equal(t, "task_b", open[0].TaskID, "token-only and terminal subscriptions are not reconciled")

			hasEarlier, err := HasEarlierOpenTaskCallbackEvent("task_a", 3)
			require.NoError(t, err)
			assert.True(t, hasEarlier)
			superseded, err := SupersedeTaskCallbackEvents("task_a", 3, []string{"queued", "running"})
			require.NoError(t, err)
			assert.EqualValues(t, 2, superseded)
			hasEarlier, err = HasEarlierOpenTaskCallbackEvent("task_a", 3)
			require.NoError(t, err)
			assert.False(t, hasEarlier)

			// Claim is exclusive across concurrent callers.
			now := common.GetTimestamp()
			due, err := ListDueTaskCallbackEvents(now, 10)
			require.NoError(t, err)
			require.Len(t, due, 1)
			event := due[0]
			assert.Equal(t, "succeeded", event.Status)
			var wins int
			var winsMu sync.Mutex
			var wg sync.WaitGroup
			for i := 0; i < 4; i++ {
				wg.Add(1)
				go func(owner string) {
					defer wg.Done()
					won, claimErr := ClaimTaskCallbackEvent(event.ID, owner, now, now+30)
					assert.NoError(t, claimErr)
					if won {
						winsMu.Lock()
						wins++
						winsMu.Unlock()
					}
				}("owner-" + string(rune('a'+i)))
			}
			wg.Wait()
			assert.Equal(t, 1, wins, "exactly one instance claims an event")
			claimed, err := GetTaskCallbackEvent(event.ID)
			require.NoError(t, err)
			assert.Equal(t, TaskCallbackStateDelivering, claimed.State)
			assert.Equal(t, 1, claimed.Attempts)

			// A stale owner cannot write back.
			applied, err := FinishTaskCallbackAttempt(event.ID, "someone-else", TaskCallbackAttemptResult{State: TaskCallbackStateSucceeded})
			require.NoError(t, err)
			assert.False(t, applied)
			applied, err = FinishTaskCallbackAttempt(event.ID, claimed.LeaseOwner, TaskCallbackAttemptResult{State: TaskCallbackStatePending, NextAttemptAt: now + 5, HTTPStatus: 502, Error: strings.Repeat("x", 600), LatencyMs: 12})
			require.NoError(t, err)
			assert.True(t, applied)
			retried, err := GetTaskCallbackEvent(event.ID)
			require.NoError(t, err)
			assert.Equal(t, TaskCallbackStatePending, retried.State)
			assert.Equal(t, 502, retried.LastHTTPStatus)
			assert.Len(t, retried.LastError, 512)
			assert.Empty(t, retried.LeaseOwner)

			// Expired leases are reclaimed (restart recovery) or failed when spent.
			won, err := ClaimTaskCallbackEvent(event.ID, "dead-instance", now+5, now+6)
			require.NoError(t, err)
			require.True(t, won)
			reclaimed, failed, err := ReclaimExpiredTaskCallbackLeases(now + 10)
			require.NoError(t, err)
			assert.EqualValues(t, 1, reclaimed)
			assert.EqualValues(t, 0, failed)
			require.NoError(t, DB.Model(&TaskCallbackEvent{}).Where("id = ?", event.ID).Update("attempts", 3).Error)
			won, err = ClaimTaskCallbackEvent(event.ID, "dead-instance", now+10, now+11)
			require.NoError(t, err)
			require.True(t, won)
			reclaimed, failed, err = ReclaimExpiredTaskCallbackLeases(now + 20)
			require.NoError(t, err)
			assert.EqualValues(t, 0, reclaimed)
			assert.EqualValues(t, 1, failed)

			// Manual redelivery re-arms the same event id.
			redelivered, err := RedeliverTaskCallbackEvent("evt_3")
			require.NoError(t, err)
			assert.True(t, redelivered)
			rearmed, err := GetTaskCallbackEvent(event.ID)
			require.NoError(t, err)
			assert.Equal(t, TaskCallbackStatePending, rearmed.State)
			assert.Zero(t, rearmed.Attempts)

			// Shutdown release returns in-flight events to pending.
			won, err = ClaimTaskCallbackEvent(event.ID, "stopping", common.GetTimestamp(), common.GetTimestamp()+30)
			require.NoError(t, err)
			require.True(t, won)
			require.NoError(t, ReleaseTaskCallbackLeases("stopping"))
			released, err := GetTaskCallbackEvent(event.ID)
			require.NoError(t, err)
			assert.Equal(t, TaskCallbackStatePending, released.State)

			// Listing.
			events, err := ListTaskCallbackEventsByTask("task_a")
			require.NoError(t, err)
			require.Len(t, events, 3)
			assert.Equal(t, []int{1, 2, 3}, []int{events[0].Seq, events[1].Seq, events[2].Seq})
			page, total, err := ListTaskCallbackEvents(TaskCallbackEventQuery{UserID: 7, State: TaskCallbackStateSuperseded}, 0, 10)
			require.NoError(t, err)
			assert.EqualValues(t, 2, total)
			assert.Len(t, page, 2)

			// Cleanup removes finished events and then event-less subscriptions.
			future := common.GetTimestamp() + 3600
			deleted, err := DeleteFinishedTaskCallbackEvents(future, 1)
			require.NoError(t, err)
			assert.EqualValues(t, 1, deleted, "batches are bounded")
			deleted, err = DeleteFinishedTaskCallbackEvents(future, 100)
			require.NoError(t, err)
			assert.EqualValues(t, 1, deleted, "the pending event is kept")
			deletedSubscriptions, err := DeleteStaleTaskCallbackSubscriptions(future, 100)
			require.NoError(t, err)
			assert.EqualValues(t, 2, deletedSubscriptions, "task_b and task_c have no events; task_a still has one")
		})
	}
}
