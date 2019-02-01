package queue

import (
	"sync"
)

// create a worker thread
func newWorker(pool chan chan Jober, wg *sync.WaitGroup) *worker {
	return &worker{
		pool:    pool,
		wg:      wg,
		jobChan: make(chan Jober),
		quit:    make(chan struct{}),
	}
}

// worker thread
type worker struct {
	pool    chan chan Jober
	wg      *sync.WaitGroup
	jobChan chan Jober
	quit    chan struct{}
}

// start the worker
func (w *worker) Start() {
	w.pool <- w.jobChan
	go w.dispatcher()
}

func (w *worker) dispatcher() {
	for {
		select {
		case j := <-w.jobChan:
			j.Job()
			w.pool <- w.jobChan
			w.wg.Done()
		case <-w.quit:
			<-w.pool
			close(w.jobChan)
			return
		}
	}
}

// stop the worker
func (w *worker) Stop() {
	close(w.quit)
}
