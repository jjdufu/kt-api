package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// recordingObserver writes one subscription row per observed transition so
// the test can check that observer writes commit (or roll back) with the
// task update.
type recordingObserver struct {
	fail     bool
	observed []string
	after    []string
}

func (o *recordingObserver) ObservesTask(task *Task) bool { return task.Action == "observed" }

func (o *recordingObserver) OnTaskTransition(tx *gorm.DB, task *Task) (bool, error) {
	o.observed = append(o.observed, string(task.Status))
	if err := tx.Create(&TaskCallbackSubscription{TaskID: task.TaskID + "-" + string(task.Status), UserID: 1}).Error; err != nil {
		return false, err
	}
	if o.fail {
		return true, errors.New("observer failed")
	}
	return true, nil
}

func (o *recordingObserver) AfterTaskTransition(task *Task) {
	o.after = append(o.after, string(task.Status))
}

func TestTaskTransitionObserverDatabaseMatrix(t *testing.T) {
	for _, dialect := range []string{"sqlite", "mysql", "postgres"} {
		t.Run(dialect, func(t *testing.T) {
			db, _ := openTaskCallbackTestDB(t, dialect)
			// The shared tasks table is not dropped (other packages may use the
			// same test database); only this test's rows are removed.
			require.NoError(t, db.AutoMigrate(&Task{}, &TaskCallbackSubscription{}, &TaskCallbackEvent{}))
			cleanupRows := func() {
				_ = db.Where("task_id IN ?", []string{"task_observed", "task_plain"}).Delete(&Task{}).Error
			}
			cleanupRows()
			t.Cleanup(cleanupRows)
			observer := &recordingObserver{}
			SetTaskTransitionObserver(observer)
			t.Cleanup(func() { SetTaskTransitionObserver(nil) })

			task := &Task{TaskID: "task_observed", Action: "observed", Status: TaskStatusSubmitted}
			require.NoError(t, db.Create(task).Error)
			stale := *task

			task.Status = TaskStatusQueued
			won, err := task.UpdateWithStatus(TaskStatusSubmitted)
			require.NoError(t, err)
			assert.True(t, won)

			// The CAS loser neither observes nor notifies.
			stale.Status = TaskStatusFailure
			won, err = stale.UpdateWithStatus(TaskStatusSubmitted)
			require.NoError(t, err)
			assert.False(t, won)
			assert.Equal(t, []string{"QUEUED"}, observer.observed)
			assert.Equal(t, []string{"QUEUED"}, observer.after)

			// An observer error rolls back only the observer's writes; the
			// task update still commits.
			observer.fail = true
			task.Status = TaskStatusInProgress
			won, err = task.UpdateWithStatus(TaskStatusQueued)
			require.NoError(t, err)
			assert.True(t, won)
			assert.Equal(t, []string{"QUEUED"}, observer.after)
			var persisted Task
			require.NoError(t, db.Where("task_id = ?", task.TaskID).Take(&persisted).Error)
			assert.EqualValues(t, TaskStatusInProgress, persisted.Status)
			var subscriptions []TaskCallbackSubscription
			require.NoError(t, db.Order("id").Find(&subscriptions).Error)
			require.Len(t, subscriptions, 1)
			assert.Equal(t, "task_observed-QUEUED", subscriptions[0].TaskID)

			// Unobserved tasks keep the plain CAS path.
			plain := &Task{TaskID: "task_plain", Action: "other", Status: TaskStatusSubmitted}
			require.NoError(t, db.Create(plain).Error)
			plain.Status = TaskStatusQueued
			won, err = plain.UpdateWithStatus(TaskStatusSubmitted)
			require.NoError(t, err)
			assert.True(t, won)
			assert.Len(t, observer.observed, 2)
		})
	}
}
