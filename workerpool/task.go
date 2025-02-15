package workerpool

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type Task struct {
	LoopTimes int
	Delay     time.Duration
	Data      interface{}
}

// StartNewTask creates a new Task and sends it to the worker pool.
func (m *WorkerPool) StartNewTask(loopTimes int, delay time.Duration, taskData interface{}) error {
	if loopTimes <= 0 {
		return errors.New("loopTimes must be greater than 0")
	}
	newTask := Task{
		LoopTimes: loopTimes,
		Delay:     delay,
		Data:      taskData,
	}

	m.taskChan <- newTask

	return nil
}

// StopAllRunningTasks stops any running task after its current task loop and delay are complete and starts new workers.
func (m *WorkerPool) StopAllRunningTasks() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.StopWorkers()
	startWorkers(m)
}

func taskWorker(timeout time.Duration, taskChan <-chan Task, stopTasksChan <-chan bool, taskFunc func(interface{}) error, wg *sync.WaitGroup) {
	defer wg.Done()

mainLoop:
	for {
		select {
		case <-stopTasksChan:
			break mainLoop
		case task := <-taskChan:
			timeoutTime := time.Now().Add(timeout)

		taskLoop:
			for i := range task.LoopTimes {
				select {
				case <-stopTasksChan:
					break mainLoop
				default:
					if time.Now().After(timeoutTime) {
						fmt.Fprintln(os.Stderr, "max timeout reached")
						break taskLoop
					}

					err := taskFunc(task.Data)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						break taskLoop
					}

					if i < task.LoopTimes-1 && task.Delay != 0 {
						time.Sleep(task.Delay)
					}
				}
			}
		}
	}
}
