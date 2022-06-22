package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	if workersCount <= 0 {
		return nil
	}

	if maxErrorsCount < 0 {
		maxErrorsCount = 0
	}

	chTask := make(chan Task, len(tasks))

	for _, task := range tasks {
		chTask <- task
	}
	close(chTask)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var errCount int32
	var mainErr error

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func(ctx context.Context, ch <-chan Task) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-ch:
					if !ok {
						cancel()

						return
					}

					if atomic.LoadInt32(&errCount) > int32(maxErrorsCount) {
						cancel()
						mainErr = ErrErrorsLimitExceeded
						return
					}

					err := task()
					if err != nil {
						atomic.AddInt32(&errCount, 1)
					}
				}
			}
		}(ctx, chTask)
	}

	wg.Wait()

	return mainErr
}
