package worker

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Worker struct {
	Id      string
	Options *WorkerOptions
	Exec    func()
}

type WorkerOptions struct {
	ExecuteAt time.Time
	Priority  int
}

func New(exec func()) Worker {
	id, _ := uuid.NewV4()
	return Worker{
		Id:   id.String(),
		Exec: exec,
	}
}

func (worker Worker) WithOptions(options *WorkerOptions) Worker {
	worker.Options = options
	return worker
}

func (worker *Worker) Execute() {
	worker.Exec()
	fmt.Printf("Run worker %s at %s\n", worker.Id, time.Now())
}

func (worker *Worker) ExecuteShedule() {
	time.AfterFunc(time.Now().Sub(worker.Options.ExecuteAt), worker.Exec)
	fmt.Printf("Run worker %s at %s\n", worker.Id, worker.Options.ExecuteAt)
}
