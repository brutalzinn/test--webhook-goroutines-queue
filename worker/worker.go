package worker

import (
	"fmt"
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/notify"
	notify_request "github.com/brutalzinn/test-webhook-goroutines-queue.git/notify/models"
	worker "github.com/brutalzinn/test-webhook-goroutines-queue.git/worker/models"
	"github.com/gofrs/uuid"
)

type Worker struct {
	Id        string
	Options   *WorkerOptions
	Service   custom_types.ServiceType
	Execution custom_types.ExecutionType
	Exec      func() worker.ExecFeedbackModel
	Notify    *notify.Notify
}

type WorkerOptions struct {
	ExecuteAt time.Time
	Priority  custom_types.Priority
}

func New(exec func() worker.ExecFeedbackModel, serviceType custom_types.ServiceType) Worker {
	id, _ := uuid.NewV4()
	return Worker{
		Id:      id.String(),
		Exec:    exec,
		Service: serviceType,
	}
}

func (worker Worker) WithOptions(options *WorkerOptions) Worker {
	worker.Options = options
	return worker
}

func (worker Worker) WithNotify(notify *notify.Notify) Worker {
	worker.Notify = notify
	return worker
}

func (worker *Worker) Execute() {
	execModel := worker.Exec()
	workerLog := WorkerLog{
		Worker:          worker,
		Status:          execModel.Status,
		RequestPayload:  execModel.Request,
		ResponsePayload: execModel.Response,
	}
	workerLog.Insert()
	fmt.Printf("Run worker %s at %s\n", worker.Id, time.Now())
	notifyBody := notify_request.NotifyBody{
		Origin: "WEBHOOK",
		Payload: map[string]any{
			"type":     custom_types.Normal,
			"status":   execModel.Status,
			"response": execModel.Response,
		},
	}
	worker.Notify.Execute(notifyBody)
}

func (worker *Worker) ExecuteShedule() {
	workerLog := WorkerLog{
		Worker:          worker,
		Status:          custom_types.Created,
		RequestPayload:  "{}",
		ResponsePayload: "{}",
	}
	id, err := workerLog.Insert()
	fmt.Printf("### Worker sheduler  %s inserted at %s executeAt:%s \n", worker.Id, time.Now(), worker.Options.ExecuteAt)
	if err != nil {
		fmt.Printf("Wrong at worker sheduler insert %s at %s\n", worker.Id, time.Now())
	}

	time.AfterFunc(worker.Options.ExecuteAt.Sub(time.Now()), func() {
		execModel := worker.Exec()
		workerLog := WorkerLog{
			Worker:          worker,
			Status:          execModel.Status,
			RequestPayload:  execModel.Request,
			ResponsePayload: execModel.Response,
		}
		err = workerLog.Update(id)
		fmt.Printf("### Worker sheduler %s update at %s executeAt:%s \n", worker.Id, time.Now(), worker.Options.ExecuteAt)
		if err != nil {
			fmt.Printf("Wrong at worker sheduler update %s at %s\n", worker.Id, time.Now())
		}
		notifyBody := notify_request.NotifyBody{
			Origin: "WEBHOOK",
			Payload: map[string]any{
				"type":     custom_types.Sheduler,
				"status":   execModel.Status,
				"response": execModel.Response,
			},
		}
		worker.Notify.Execute(notifyBody)
	})
	notifyBody := notify_request.NotifyBody{
		Origin: "WEBHOOK",
		Payload: map[string]any{
			"type":     custom_types.Sheduler,
			"status":   custom_types.Created,
			"response": nil,
		},
	}
	worker.Notify.Execute(notifyBody)
}
