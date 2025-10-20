package loex

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrPanic = errors.New("panic")
)

// ExecTasks execute multiple tasks concurrently using Go routines
// maxConcurrentTasks behaves similarly as `pool size`, pass 0 to set no limit.
// In case you want to cancel the execution, use context.WithTimeout() or context.WithCancel().
// nolint: gocognit
func ExecTasks(
	ctx context.Context,
	maxConcurrentTasks uint,
	stopOnError bool,
	tasks ...func(ctx context.Context) error,
) map[int]error {
	taskCount := len(tasks)
	if taskCount == 0 {
		return nil
	}

	type execTaskResult struct {
		Index int
		Error error
	}

	stopped := &atomic.Bool{}
	resultChan := make(chan *execTaskResult, taskCount)
	var limiterChan chan struct{}
	if maxConcurrentTasks != 0 && maxConcurrentTasks < uint(taskCount) {
		limiterChan = make(chan struct{}, maxConcurrentTasks)
	} else {
		maxConcurrentTasks = 0
	}

	for i := 0; i < taskCount; i++ {
		// In case we set pool size, when out of slot, this will wait until one to be available again
		if maxConcurrentTasks != 0 {
			limiterChan <- struct{}{}
		}

		go func(i int, task func(ctx context.Context) error) {
			defer func() {
				// In case we set pool size, release the slot when the task ends
				if maxConcurrentTasks != 0 {
					<-limiterChan
				}

				r := recover()
				if r != nil && !stopped.Load() {
					resultChan <- &execTaskResult{Index: i, Error: fmt.Errorf("%w: %v", ErrPanic, r)}
				}
			}()

			if stopOnError && stopped.Load() {
				return
			}

			if err := ctx.Err(); err != nil {
				resultChan <- &execTaskResult{Index: i, Error: err}
				return
			}

			err := task(ctx)
			if err != nil {
				resultChan <- &execTaskResult{Index: i, Error: err}
			} else {
				resultChan <- nil
			}
		}(i, tasks[i])
	}

	errResult := map[int]error{}
	for i := 0; i < taskCount; i++ {
		res := <-resultChan
		if res == nil {
			continue
		}
		errResult[res.Index] = res.Error
		if stopOnError {
			stopped.Store(true)
			break
		}
	}
	return errResult
}

// ExecTaskFunc executes a function on every target objects
func ExecTaskFunc[T any](
	ctx context.Context,
	maxConcurrentTasks uint,
	stopOnError bool,
	taskFunc func(ctx context.Context, obj T, index int) error,
	targetObjects ...T,
) map[int]error {
	tasks := make([]func(ctx context.Context) error, len(targetObjects))
	for i := range targetObjects {
		obj := targetObjects[i]
		tasks[i] = func(ctx context.Context) error {
			return taskFunc(ctx, obj, i)
		}
	}
	return ExecTasks(ctx, maxConcurrentTasks, stopOnError, tasks...)
}
