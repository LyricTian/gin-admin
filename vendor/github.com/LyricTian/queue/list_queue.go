package queue

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// NewListQueue create a list queue that specifies the number of worker threads
func NewListQueue(maxThread int) *ListQueue {
	return NewListQueueWithMaxLen(maxThread, 0)
}

// NewListQueueWithMaxLen create a list queue that specifies the number of worker threads
// and the maximum number of elements
func NewListQueueWithMaxLen(maxThread, maxLen int) *ListQueue {
	return &ListQueue{
		maxLen:     maxLen,
		maxWorker:  maxThread,
		workers:    make([]*worker, maxThread),
		workerPool: make(chan chan Jober, maxThread),
		list:       list.New(),
		lock:       new(sync.RWMutex),
		wg:         new(sync.WaitGroup),
	}
}

// ListQueue a list task queue for mitigating server pressure in high concurrency situations
// and improving task processing
type ListQueue struct {
	maxLen     int
	maxWorker  int
	workers    []*worker
	workerPool chan chan Jober
	list       *list.List
	lock       *sync.RWMutex
	wg         *sync.WaitGroup
	running    uint32
}

// Run start running queues
func (q *ListQueue) Run() {
	if atomic.LoadUint32(&q.running) == 1 {
		return
	}
	atomic.StoreUint32(&q.running, 1)

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i] = newWorker(q.workerPool, q.wg)
		q.workers[i].Start()
	}

	go q.dispatcher()
}

func (q *ListQueue) dispatcher() {
	for {
		q.lock.RLock()
		if atomic.LoadUint32(&q.running) != 1 && q.list.Len() == 0 {
			q.lock.RUnlock()
			break
		}
		ele := q.list.Front()
		q.lock.RUnlock()

		if ele == nil {
			time.Sleep(time.Millisecond * 10)
			continue
		}

		worker := <-q.workerPool
		worker <- ele.Value.(Jober)

		q.lock.Lock()
		q.list.Remove(ele)
		q.lock.Unlock()
	}
}

// Push put the executable task into the queue
func (q *ListQueue) Push(job Jober) {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	if q.maxLen > 0 {
		q.lock.RLock()
		if q.list.Len() > q.maxLen {
			q.lock.RUnlock()
			return
		}
		q.lock.RUnlock()
	}

	q.wg.Add(1)
	q.lock.Lock()
	q.list.PushBack(job)
	q.lock.Unlock()
}

// Terminate terminate the queue to receive the task and release the resource
func (q *ListQueue) Terminate() {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	atomic.StoreUint32(&q.running, 0)
	q.wg.Wait()

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}
