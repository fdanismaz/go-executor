package executor

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
	"sync"
)

type Parallel struct {
	concurrencyLimit int64
	tasks            []Task
	finished         bool
	continueOnError  bool
}

func NewParallel(concurrencyLimit int64, continueOnError bool) Executor {
	return &Parallel{
		concurrencyLimit: concurrencyLimit,
		continueOnError: continueOnError,
	}
}

func (e *Parallel) Submit(t Task) {
	if t != nil {
		e.tasks = append(e.tasks, t)
	}
}

func (e *Parallel) Execute() error {
	var waitGroup sync.WaitGroup
	var executorError error
	waitGroup.Add(len(e.tasks))

	s := semaphore.NewWeighted(e.concurrencyLimit)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldStop := false

	// A decorator function that releases semaphore resource after the function execution is completed.
	// It also notifies the WaitGroup so that it waits for one less job to be completed
	executeTask := func(f func() error, taskName string, errChan chan error) func() {
		return func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Panic of the task %s recovered\n", taskName)
				}
			}()
			defer waitGroup.Done()
			defer s.Release(1)

			if shouldStop {
				fmt.Printf("The executor is stopped. Cancelling task %s\n", taskName)
				return
			}

			defer fmt.Printf("Task completed: %s\n", taskName)
			fmt.Printf("Starting task: %s\n", taskName)
			err := f()
			if err != nil {
				errChan <- errors.Wrap(err, fmt.Sprintf("An error occurred during executing the task %s\n", taskName))
			}
		}
	}

	// More than one go routine can write to errChan at the same time. Tasks should never be blocked...
	errChan := make(chan error, e.concurrencyLimit)

	// Error channel and Context Done listener
	go func() {
		for {
			select {
			case err := <-errChan:
				executorError = err
				if !e.continueOnError {
					fmt.Println("Shutting down the parallel executor")
					cancel() // Cancel the context
				}
			case <-ctx.Done():
				shouldStop = true
				return
			}
		}
	}()

	for _, task := range e.tasks {
		err := s.Acquire(ctx, 1)
		if err != nil {
			fmt.Printf("Cannot acquire resource for task %s. Skipping it. Reason: %v\n", task.GetName(), err)
			// Since the task is skipped, we should call waitGroup.Done for it.
			waitGroup.Done()
			continue
		}
		run := executeTask(task.Run, task.GetName(), errChan)
		go run()
	}
	waitGroup.Wait()
	e.finished = true
	return executorError
}

func (e *Parallel) IsFinished() bool {
	return e.finished
}