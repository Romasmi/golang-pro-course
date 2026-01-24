package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")

		// Check specific errors
		require.Equal(t, Run(tasks, workersCount, 0), ErrErrorsLimitExceeded)
		require.Equal(t, Run(tasks, 0, maxErrorsCount), ErrInvalidWorkersCount)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var concurrentTasksCount, maxConcurrentTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				current := atomic.AddInt32(&concurrentTasksCount, 1)
				defer atomic.AddInt32(&concurrentTasksCount, -1)
				for {
					peak := atomic.LoadInt32(&maxConcurrentTasksCount)
					if current <= peak {
						break
					}
					if atomic.CompareAndSwapInt32(&maxConcurrentTasksCount, peak, current) {
						break
					}
				}

				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")

		require.GreaterOrEqual(t, atomic.LoadInt32(&maxConcurrentTasksCount), int32(workersCount),
			"should have at least %d tasks running concurrently", workersCount)
	})
}
