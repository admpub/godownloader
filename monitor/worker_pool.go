package monitor

import (
	"sync/atomic"
)

type WorkerPool struct {
	workers map[string]*MonitoredWorker
	total   int32
	done    int32
}

func (wp *WorkerPool) AppendWork(iv *MonitoredWorker) {
	if wp.workers == nil {
		wp.workers = make(map[string]*MonitoredWorker)
	}
	iv.ondone = func() {
		atomic.AddInt32(&wp.done, 1)
	}
	wp.workers[iv.GetId()] = iv
	atomic.AddInt32(&wp.total, 1)
}

func (wp *WorkerPool) Completed() bool {
	return atomic.LoadInt32(&wp.total) == atomic.LoadInt32(&wp.done)
}

func (wp *WorkerPool) StartAll() []error {
	var errs []error
	for _, value := range wp.workers {
		if err := value.Start(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (wp *WorkerPool) StopAll() []error {
	var errs []error
	for _, value := range wp.workers {
		if err := value.Stop(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (wp *WorkerPool) GetAllProgress() interface{} {
	var pr []interface{}
	for _, value := range wp.workers {
		pr = append(pr, value.GetProgress())
	}
	return pr
}
