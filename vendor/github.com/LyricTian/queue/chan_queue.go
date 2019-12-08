package queue

import (
	"sync"
	"sync/atomic"
)

// NewQueue create a queue that specifies the number of buffers and the number of worker threads
func NewQueue(maxCapacity, maxThread int) *Queue {
	return &Queue{
		jobQueue:   make(chan Jober, maxCapacity),
		maxWorkers: maxThread,
		workerPool: make(chan chan Jober, maxThread),
		workers:    make([]*worker, maxThread),
		wg:         new(sync.WaitGroup),
	}
}

// Queue a task queue for mitigating server pressure in high concurrency situations
// and improving task processing
type Queue struct {
	maxWorkers int
	jobQueue   chan Jober
	workerPool chan chan Jober
	workers    []*worker
	running    uint32
	wg         *sync.WaitGroup
}

// Run start running queues
func (q *Queue) Run() {
	if atomic.LoadUint32(&q.running) == 1 {
		return
	}

	atomic.StoreUint32(&q.running, 1)
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i] = newWorker(q.workerPool, q.wg)
		q.workers[i].Start()
	}

	go q.dispatcher()
}

func (q *Queue) dispatcher() {
	for job := range q.jobQueue {
		worker := <-q.workerPool
		worker <- job
	}
}

// Terminate terminate the queue to receive the task and release the resource
func (q *Queue) Terminate() {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	atomic.StoreUint32(&q.running, 0)
	q.wg.Wait()

	close(q.jobQueue)
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}

// Push put the executable task into the queue
func (q *Queue) Push(job Jober) {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	q.wg.Add(1)
	q.jobQueue <- job
}
