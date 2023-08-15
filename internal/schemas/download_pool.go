// https://medium.com/@zufolo/a-pattern-for-limiting-the-number-of-goroutines-in-execution-56e13b226e72

package schemas

import "sync"

type DownloadPool struct {
	maxWorkers int
	queue      chan *AttachmentProc
	wg         sync.WaitGroup
}

func NewDownloadPool(concurrencyLimit int) *DownloadPool {
	dpl := &DownloadPool{
		maxWorkers: concurrencyLimit,
		queue:      make(chan *AttachmentProc),
		wg:         sync.WaitGroup{},
	}
	dpl.init()

	return dpl
}

// Add
// put a new attachment download request into the download queue
func (dpl *DownloadPool) Add(proc *AttachmentProc) {
	dpl.queue <- proc
}

// init
// build out the download workers and set them listening to the queue
func (dpl *DownloadPool) init() {
	// make a number of worker goroutines
	for worker := 0; worker < dpl.maxWorkers; worker++ {
		// each worker is a wait group
		dpl.wg.Add(1)

		// start a new goroutine
		go func() {
			// if the worker has no more to do (function has exited)
			// then mark its wait group as done
			defer dpl.wg.Done()

			// "ranging" a channel is a blocking call that ticks when new
			// data is pushed on the channel.
			// This allows each worker to listen to the queue channel, pull
			// a download request from it, and execute the request.
			for proc := range dpl.queue {
				proc.Execute()
			}
		}()
	}
}

func (dpl *DownloadPool) Process() {
	// since our channel does not specify a length we need to close it
	// this will make it so that the workers don't wait for more data
	// once the current data in the queue is exhausted
	close(dpl.queue)

	// wait for all the wait groups to be done (all worker goroutines have exited)
	dpl.wg.Wait()
}
