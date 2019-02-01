package queue

var (
	internalQueue Queuer
)

// Queuer a task queue for mitigating server pressure in high concurrency situations
// and improving task processing
type Queuer interface {
	Run()
	Push(job Jober)
	Terminate()
}

// Run start running queues,
// specify the number of buffers, and the number of worker threads
func Run(maxCapacity, maxThread int) {
	if internalQueue == nil {
		internalQueue = NewQueue(maxCapacity, maxThread)
	}
	internalQueue.Run()
}

// RunListQueue start running list queues
// ,specify the number of worker threads
func RunListQueue(maxThread int) {
	if internalQueue == nil {
		internalQueue = NewListQueue(maxThread)
	}
	internalQueue.Run()
}

// Push put the executable task into the queue
func Push(job Jober) {
	if internalQueue == nil {
		return
	}
	internalQueue.Push(job)
}

// Terminate terminate the queue to receive the task and release the resource
func Terminate() {
	if internalQueue == nil {
		return
	}
	internalQueue.Terminate()
}
