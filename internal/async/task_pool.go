// https://medium.com/@zufolo/a-pattern-for-limiting-the-number-of-goroutines-in-execution-56e13b226e72

package async

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"sync"
)

type Task interface {
	Run() error
}

type TaskPool struct {
	maxWorkers    int
	queueChan     chan Task
	errorChan     chan error
	cancelContext context.Context
	waitGroup     sync.WaitGroup
	isClosed      bool
	err           *PoolError
	cleanupCb     func()
}

type PoolError struct {
	taskErrors []error
}

func (e *PoolError) Error() string {
	sb := new(strings.Builder)

	for _, e := range e.taskErrors {
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}

func NewTaskPool(
	cancelContext context.Context,
	cleanupCb func(),
) *TaskPool {
	pool := &TaskPool{
		maxWorkers:    runtime.NumCPU() * 2,
		queueChan:     make(chan Task),
		errorChan:     make(chan error),
		cancelContext: cancelContext,
		waitGroup:     sync.WaitGroup{},
		isClosed:      false,
		err:           &PoolError{},
		cleanupCb:     cleanupCb,
	}
	pool.init()

	return pool
}

// AddTask
// put a new attachment download request into the download queue
func (pool *TaskPool) AddTask(task Task) error {
	if pool.isClosed {
		return errors.New("TaskPool has been closed")
	}
	pool.queueChan <- task
	return nil
}

// init
// build out the download workers and set them listening to the queue
func (pool *TaskPool) init() {
	go func() {
		for taskError := range pool.errorChan {
			select {
			case <-pool.cancelContext.Done():
				return
			default:
				log.Debug("handling task error")
				pool.err.taskErrors = append(pool.err.taskErrors, taskError)
			}
		}
	}()

	// make a number of worker goroutines
	for worker := 0; worker < pool.maxWorkers; worker++ {
		// each worker is a wait group
		pool.waitGroup.Add(1)

		// start a new goroutine as an event loop for this worker
		go func() {
			// if the worker has no more to do (function has exited)
			// then mark its wait group as done
			defer pool.waitGroup.Done()

			// "ranging" a channel is a blocking call that ticks when new
			// data is pushed on the channel.
			// This allows each worker to listen to the queue channel, pull
			// a download request from it, and execute the request.
			for nextTask := range pool.queueChan {
				select {
				case <-pool.cancelContext.Done():
					return
				default:
					err := nextTask.Run()
					if err != nil {
						pool.errorChan <- err
					}
				}
			}
		}()
	}
}

func (pool *TaskPool) Wait() {
	// since our channel does not specify a length we need to close it
	// this will make it so that the workers don't wait for more data
	// once the current data in the queue is exhausted
	pool.isClosed = true
	close(pool.queueChan)
	// wait for all the wait groups to be done (all worker goroutines have exited)
	pool.waitGroup.Wait()
	close(pool.errorChan)

	if pool.cleanupCb != nil {
		pool.cleanupCb()
	}
}

func (pool *TaskPool) Errors() error {
	if len(pool.err.taskErrors) > 0 {
		return pool.err
	}

	return nil
}
