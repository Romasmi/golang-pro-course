package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkersCount = errors.New("invalid workers count")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrInvalidWorkersCount
	}

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	workersCount := min(n, len(tasks))
	jobs := make(chan Task, workersCount)
	errsCounter := &atomic.Int64{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer close(jobs)
		for _, v := range tasks {
			select {
			case <-ctx.Done():
				return
			case jobs <- v:
			}
		}
	}()

	wg := sync.WaitGroup{}
	for range workersCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, jobs, errsCounter, m, cancel)
		}()
	}
	wg.Wait()

	if errsCounter.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(
	ctx context.Context, jobs <-chan Task, errCounter *atomic.Int64, maxErrCount int, cancel context.CancelFunc,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			if job() != nil {
				if errCounter.Add(1) >= int64(maxErrCount) {
					cancel()
					return
				}
			}
		}
	}
}
