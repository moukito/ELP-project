package worker

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
