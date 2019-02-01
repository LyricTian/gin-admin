package queue

// Jober an asynchronous task that can be executed
type Jober interface {
	Job()
}

// SyncJober a synchronization task that can be executed
type SyncJober interface {
	Jober
	Wait() <-chan interface{}
	Error() error
}

type job struct {
	v        interface{}
	callback func(interface{})
}

// NewJob create an asynchronous task
func NewJob(v interface{}, fn func(interface{})) Jober {
	return &job{
		v:        v,
		callback: fn,
	}
}

func (j *job) Job() {
	j.callback(j.v)
}

type syncJob struct {
	err      error
	result   chan interface{}
	v        interface{}
	callback func(interface{}) (interface{}, error)
}

// NewSyncJob create a synchronization task
func NewSyncJob(v interface{}, fn func(interface{}) (interface{}, error)) SyncJober {
	return &syncJob{
		result:   make(chan interface{}, 1),
		v:        v,
		callback: fn,
	}
}

func (j *syncJob) Job() {
	result, err := j.callback(j.v)
	if err != nil {
		j.err = err
		close(j.result)
		return
	}

	j.result <- result

	close(j.result)
}

func (j *syncJob) Wait() <-chan interface{} {
	return j.result
}

func (j *syncJob) Error() error {
	return j.err
}
