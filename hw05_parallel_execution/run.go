package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
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
	var mutex sync.RWMutex
	var errCount int
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

					mutex.RLock()
					if errCount > maxErrorsCount {
						cancel()
						mainErr = ErrErrorsLimitExceeded
						mutex.RUnlock()
						return
					}
					mutex.RUnlock()

					err := task()
					if err != nil {
						mutex.Lock()
						errCount++
						mutex.Unlock()
					}
				}
			}
		}(ctx, chTask)
	}

	wg.Wait()

	return mainErr
}
