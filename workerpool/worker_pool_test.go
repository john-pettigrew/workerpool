package workerpool

import (
	"sync"
	"testing"
	"time"
)

func TestStopWorkers(t *testing.T) {
	const expectedFunctionCount int = 2

	mu := sync.Mutex{}
	functionCount := 0
	testFunc := (func(taskData interface{}) error {
		mu.Lock()
		defer mu.Unlock()

		functionCount += 1
		return nil
	})

	pool := NewWorkerPool(
		3,
		time.Duration(10)*time.Second,
		testFunc,
	)

	err := pool.StartNewTask(
		10,
		time.Millisecond*50,
		nil,
	)

	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Millisecond * 75)

	if functionCount != expectedFunctionCount {
		t.Errorf("function called incorrect number of times")
		return
	}

	pool.StopWorkers()
	if _, ok := <-pool.stopTasksChan; ok {
		t.Errorf("channel should be closed")
		return
	}

	time.Sleep(time.Millisecond * 200)

	if functionCount != expectedFunctionCount {
		t.Errorf("function continued to be called")
		return
	}
}
