package workerpool

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestStartNewTask(t *testing.T) {
	type TaskTest struct {
		NumberWorkers         int
		MaxTaskLoopDuration   time.Duration
		LoopCount             int
		LoopDelay             time.Duration
		TestWaitDuration      time.Duration
		NumberTasks           int
		ExpectedFunctionCount int
	}
	tests := []TaskTest{
		// single worker
		{
			NumberWorkers:         1,
			MaxTaskLoopDuration:   time.Duration(10) * time.Second,
			LoopCount:             1,
			LoopDelay:             time.Millisecond * 5,
			NumberTasks:           1,
			ExpectedFunctionCount: 1,
			TestWaitDuration:      (time.Millisecond * 5) + (time.Millisecond * 100),
		},
		// test delay
		{
			NumberWorkers:         1,
			MaxTaskLoopDuration:   time.Duration(10) * time.Second,
			LoopCount:             10,
			LoopDelay:             time.Millisecond * 10,
			NumberTasks:           1,
			ExpectedFunctionCount: 3,
			TestWaitDuration:      time.Millisecond * 25,
		},
		// multiple loops
		{
			NumberWorkers:         1,
			MaxTaskLoopDuration:   time.Duration(10) * time.Second,
			LoopCount:             10,
			LoopDelay:             time.Millisecond * 5,
			NumberTasks:           1,
			ExpectedFunctionCount: 10,
			TestWaitDuration:      (time.Millisecond * 5 * 10) + (time.Millisecond * 100),
		},
		// multiple workers
		{
			NumberWorkers:         3,
			MaxTaskLoopDuration:   time.Duration(10) * time.Second,
			LoopCount:             10,
			LoopDelay:             time.Millisecond * 5,
			NumberTasks:           1,
			ExpectedFunctionCount: 10,
			TestWaitDuration:      (time.Millisecond * 5 * 10) + (time.Millisecond * 100),
		},
		// multiple tasks
		{
			NumberWorkers:         3,
			MaxTaskLoopDuration:   time.Duration(10) * time.Second,
			LoopCount:             10,
			LoopDelay:             time.Millisecond * 5,
			NumberTasks:           3,
			ExpectedFunctionCount: 30,
			TestWaitDuration:      (time.Millisecond * 5 * 10 * 3) + (time.Millisecond * 100),
		},
		// more tasks than workers
		{
			NumberWorkers:         1,
			MaxTaskLoopDuration:   time.Duration(10) * time.Second,
			LoopCount:             2,
			LoopDelay:             time.Millisecond * 200,
			NumberTasks:           2,
			ExpectedFunctionCount: 3,
			TestWaitDuration:      0,
		},
		// MaxTaskLoopDuration reached
		{
			NumberWorkers:         1,
			MaxTaskLoopDuration:   time.Millisecond * 5,
			LoopCount:             10,
			LoopDelay:             time.Millisecond * 10,
			NumberTasks:           1,
			ExpectedFunctionCount: 1,
			TestWaitDuration:      time.Millisecond * 100,
		},
	}

	testTaskData := 42
	for _, currentTest := range tests {
		mu := sync.Mutex{}
		functionCount := 0

		testFunc := (func(taskData interface{}) error {
			mu.Lock()
			defer mu.Unlock()

			readValue, ok := taskData.(int)
			if !ok || readValue != testTaskData {
				const err string = "incorrect value passed to function"
				t.Errorf(err)
				return errors.New(err)
			}

			functionCount += 1
			return nil
		})

		pool := NewWorkerPool(
			currentTest.NumberWorkers,
			currentTest.MaxTaskLoopDuration,
			testFunc,
		)

		for range currentTest.NumberTasks {
			err := pool.StartNewTask(
				currentTest.LoopCount,
				currentTest.LoopDelay,
				testTaskData,
			)

			if err != nil {
				t.Error(err)
				return
			}
		}

		time.Sleep(currentTest.TestWaitDuration)

		if functionCount != currentTest.ExpectedFunctionCount {
			t.Errorf("function called incorrect number of times")
		}

		pool.StopWorkers()
	}
}

func TestStopAllRunningTasks(t *testing.T) {
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

	pool.StopAllRunningTasks()

	time.Sleep(time.Millisecond * 200)

	if functionCount != expectedFunctionCount {
		t.Errorf("function continued to be called after tasks stopped")
		return
	}

	pool.StopWorkers()
}
