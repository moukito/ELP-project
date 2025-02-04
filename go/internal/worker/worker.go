package worker

/*
Package worker provides a flexible and scalable worker pool architecture to handle concurrent task processing.

---

### Task[T any, R any]:
Represents a single unit of work to be executed by the worker pool.
The struct uses Go generics to support various types for input (T) and output (R).

Fields:
- `Conn net.Conn`: Represents the associated network connection for the task.
- `Input T`: The input data for the task.
- `Output R`: The result of task processing.
- `Err error`: Captures any error that occurs during task processing.
- `ResultChan chan Task[T, R]`: A channel to communicate results after task completion.
- `Function func(T) (R, error)`: A user-defined function to process the task.

---

### StartWorkerPool[T any, R any]:
Starts a pool of workers that process tasks from a channel concurrently.

Parameters:
- `name string`: Name of the worker pool (useful for logging).
- `numWorkers int`: Number of workers in the pool.
- `workerFunc func(Task[T, R])`: Function executed by each worker to process tasks.
- `tasks <-chan Task[T, R]`: Channel from which workers fetch tasks for processing.

Behavior:
- Creates `numWorkers` goroutines, each executing the provided `workerFunc`.
- Logs when workers start and stop.
- Processes tasks continuously until the `tasks` channel is closed.

Example Usage:
```go
tasks := make(chan Task[int, string])
StartWorkerPool("ExamplePool", 3, TreatmentWorker[int, string], tasks)
```

---

### TreatmentWorker[T any, R any]:
Processes a single task by applying its associated function and manages the task's lifecycle.

Parameters:
- `task Task[T, R]`: A task to be processed.

Behavior:
1. Logs the start of task processing.
2. If no `Function` is provided, logs an error, sets the `Err` field, and sends the result back via `ResultChan` (if specified).
3. Executes the `Function` with `Input`, stores the result in `Output`, and captures any errors in `Err`.
4. Sends the processed task back via `ResultChan` for further handling (if specified).
5. Logs the conclusion of task processing.

Example Usage:
```go
task := Task[int, string]{
    Conn: conn, // Provide a net.Conn instance
    Input: 42,
    Function: func(input int) (string, error) {
        return fmt.Sprintf("Processed: %d", input), nil
    },
    ResultChan: resultChan, // Channel to receive results
}
TreatmentWorker(task)
```

---

### Logging:
- Logs worker activity (start/stop) and individual task processing events.
- Transparent error reporting via structured logging, aiding troubleshooting and monitoring.

### Scalability:
- Uses goroutines for concurrency, enabling efficient task handling even at large scales.
- Flexible and reusable workers adapt to varying types and numbers of tasks dynamically.

### Use-Cases:
- Task delegation in server applications, especially for operations associated with network requests (`net.Conn`).
- Parallel processing of computational tasks, I/O-bound activities, or transformations on streams of data.

*/

import (
	"errors"
	"log"
	"net"
)

type Task[T any, R any] struct {
	Conn       net.Conn
	Input      T
	Output     R
	Err        error
	ResultChan chan Task[T, R]
	Function   func(T) (R, error)
}

func StartWorkerPool[T any, R any](name string, numWorkers int, workerFunc func(Task[T, R]), tasks <-chan Task[T, R]) {
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			log.Printf("%s Worker %d started", name, workerID)
			for task := range tasks {
				workerFunc(task)
			}
			log.Printf("%s Worker %d stopped", name, workerID)
		}(i)
	}
}

func TreatmentWorker[T any, R any](task Task[T, R]) {
	log.Printf("Processing task for connection: %v", task.Conn.RemoteAddr())

	if task.Function == nil {
		task.Err = errors.New("no processing function provided")
		if task.ResultChan != nil {
			task.ResultChan <- task
		}
		log.Printf("No function provided for task from: %v", task.Conn.RemoteAddr())
		return
	}

	output, err := task.Function(task.Input)
	task.Output = output
	task.Err = err

	if task.ResultChan != nil {
		task.ResultChan <- task
	}
	log.Printf("Task processing completed for connection: %v", task.Conn.RemoteAddr())
}
