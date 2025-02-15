# 🔧 Go Worker Pool
A Go library for running tasks in a worker pool.

## Features
- Run a task multiple times with a delay between loops.
- Limit the number of concurrent tasks to run at once.

## Notes
- If a running task is stopping early (stopped workers or stopped tasks or taskTimeoutDuration is exceeded) the current loop and delay are allowed to complete before the task is stopped.
- "StartNewTask" will block until there is a worker available.
- A worker only works on a single task until all loops and delays for that task are complete.
- A task is stopped if an error is returned from the "taskFunc".

## Getting Started
### Example Usage
```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/john-pettigrew/workerpool/workerpool"
)

func printMessage(taskData interface{}) error {
	num, ok := taskData.(int)
	if !ok {
		return errors.New("invalid data")
	}

	fmt.Println(num)
	return nil
}

func main() {
	pool := workerpool.NewWorkerPool(
		3,                             // number of workers in pool
		time.Duration(10)*time.Second, // task timeout duration including loops and delays
		printMessage,                  // function to run
	)

	for i := range 3 {
		err := pool.StartNewTask(
			10,          // number of times to loop
			time.Second, // delay between loops
			i,           // data to send to taskFunc
		)

		if err != nil {
			fmt.Println("error starting task")
			return
		}
	}

	time.Sleep(time.Millisecond * 2500)

	// Stop tasks early after current task loops and delays are complete
	pool.StopAllRunningTasks()
}
```