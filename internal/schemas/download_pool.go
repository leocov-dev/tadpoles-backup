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

func (dpl *DownloadPool) Add(proc *AttachmentProc) {
	dpl.queue <- proc
}

func (dpl *DownloadPool) init() {
	for worker := 0; worker < dpl.maxWorkers; worker++ {
		dpl.wg.Add(1)

		go func(worker int) {
			defer dpl.wg.Done()

			for proc := range dpl.queue {
				proc.Execute()
			}
		}(worker)
	}
}

func (dpl *DownloadPool) Process() {
	close(dpl.queue)
	dpl.wg.Wait()
}
