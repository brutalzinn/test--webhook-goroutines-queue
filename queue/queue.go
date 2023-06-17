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

/* func (queue Queue) insert() (response QueueResponse) {
	database_url := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close(context.Background())
	err = conn.QueryRow(context.Background(), "INSERT INTO queue (name, request_payload, response_payload, priority, status) VALUES ($1, $2, $3, $4, $5) returning id", queue.Name, queue.Request, queue.Response, queue.Priority, queue.Status).Scan(&response.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	retNewV4()
}

func (queue Queue) update(new_queue Queue) (response QueueResponse) {
	database_url := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE queue set payload_request=$1, payload_response=$2, priority=$3, status=$4", queue.Request, queue.Response, queue.Priority, queue.Status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	return
} */
