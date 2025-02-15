package workerpool

import (
	"sync"
	"time"
)

type WorkerPool struct {
	mu                  sync.Mutex
	wg                  sync.WaitGroup
	taskChan            chan Task
	stopTasksChan       chan bool
	numWorkers          int
	taskTimeoutDuration time.Duration
	taskFunc            func(interface{}) error
}

// NewWorkerPool creates a new WorkerPool and starts the workers.
func NewWorkerPool(numWorkers int, taskTimeoutDuration time.Duration, taskFunc func(interface{}) error) *WorkerPool {
	newWorkerPool := WorkerPool{
		taskChan:            make(chan Task),
		stopTasksChan:       make(chan bool),
		numWorkers:          numWorkers,
		taskTimeoutDuration: taskTimeoutDuration,
		taskFunc:            taskFunc,
		wg:                  sync.WaitGroup{},
	}

	startWorkers(&newWorkerPool)
	return &newWorkerPool
}

// StopWorkers stops the workers in the WorkerPool. The WorkerPool should not be used after StopWorkers is called.
func (m *WorkerPool) StopWorkers() {
	close(m.stopTasksChan)
	m.wg.Wait()
}

func startWorkers(m *WorkerPool) {
	m.stopTasksChan = make(chan bool)
	for range m.numWorkers {
		m.wg.Add(1)
		go taskWorker(m.taskTimeoutDuration, m.taskChan, m.stopTasksChan, m.taskFunc, &m.wg)
	}
}
