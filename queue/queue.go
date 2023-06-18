package queue

import (
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/worker"
)

type Queue struct {
	Workers []worker.Worker
}

func (queue *Queue) Enqueue(worker worker.Worker) {
	queue.Workers = append(queue.Workers, worker)
}

func (queue *Queue) Current() worker.Worker {
	item := queue.Workers[0]
	return item
}
func (queue *Queue) Dequeue() worker.Worker {
	item := queue.Workers[0]
	queue.Workers = queue.Workers[1:]
	return item
}

func (queue *Queue) IsEmpty() bool {
	return len(queue.Workers) == 0
}
