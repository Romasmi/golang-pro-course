package hw05parallelexecution

import (
	"errors"
	"math"
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

	jobs := make(chan Task, len(tasks))
	errsCounter := &atomic.Int64{}
	go func() {
		for _, v := range tasks {
			jobs <- v
			if errsCounter.Load() >= int64(m) {
				return
			}
		}
		defer close(jobs)
	}()

	wg := sync.WaitGroup{}

	workersCount := int(math.Min(float64(n), float64(len(tasks))))
	for range workersCount {
		wg.Go(func() {
			worker(jobs, errsCounter, m)
		})
	}
	wg.Wait()

	if errsCounter.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(jobs <-chan Task, errCounter *atomic.Int64, maxErrCount int) {
	for job := range jobs {
		if job() != nil {
			errCounter.Add(1)
		}
		if errCounter.Load() >= int64(maxErrCount) {
			return
		}
	}
}
