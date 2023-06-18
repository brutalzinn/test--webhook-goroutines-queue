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
	Id            string
	Options       *WorkerOptions
	Service       custom_types.ServiceType
	ExecutionType custom_types.ExecutionType
	Exec          func() worker.FeedbackModel
	Notify        *notify.Notify
}

type WorkerOptions struct {
	ExecuteAt time.Time
	Priority  custom_types.Priority
}

type WorkerCompleted struct {
	FeedbackModel worker.FeedbackModel
	ExecutionType custom_types.ExecutionType
	Notify        *notify.Notify
}

type WorkerNotify struct {
	Origin  string
	Payload notify_request.NotifyPayload
}

func New(exec func() worker.FeedbackModel, serviceType custom_types.ServiceType) Worker {
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
	execFeedback := worker.Exec()
	workerLog := WorkerLog{
		Worker:          worker,
		Status:          execFeedback.Status,
		RequestPayload:  execFeedback.Request,
		ResponsePayload: execFeedback.Response,
	}
	workerLog.Insert()
	fmt.Printf("Run worker %s at %s\n", worker.Id, time.Now())
	notifyBody := notify_request.NotifyBody{
		Origin: "WEBHOOK_NORMAL",
		Payload: notify_request.NotifyPayload{
			Type:     custom_types.Normal,
			Status:   execFeedback.Status,
			Response: &execFeedback.Response,
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
	workerNotify := WorkerNotify{
		Origin: "WEBHOOK_SHEDULER",
		Payload: notify_request.NotifyPayload{
			Type:     custom_types.Sheduler,
			Status:   custom_types.Created,
			Response: nil,
		},
	}
	worker.sendNotify(workerNotify)

	time.AfterFunc(worker.Options.ExecuteAt.Sub(time.Now()), func() {
		feedbackModel := worker.Exec()
		workerLog := WorkerLog{
			Worker:          worker,
			Status:          feedbackModel.Status,
			RequestPayload:  feedbackModel.Request,
			ResponsePayload: feedbackModel.Response,
		}
		err = workerLog.Update(id)
		fmt.Printf("### Worker sheduler %s update at %s executeAt:%s \n", worker.Id, time.Now(), worker.Options.ExecuteAt)
		if err != nil {
			fmt.Printf("Wrong at worker sheduler update %s at %s\n", worker.Id, time.Now())
		}
		workerNotify := WorkerNotify{
			Origin: "WEBHOOK_SHEDULER",
			Payload: notify_request.NotifyPayload{
				Type:     custom_types.Sheduler,
				Status:   feedbackModel.Status,
				Response: &feedbackModel.Response,
			},
		}
		worker.sendNotify(workerNotify)
	})
}

func (worker *Worker) sendNotify(workerNotify WorkerNotify) {
	notifyBody := notify_request.NotifyBody{
		Origin:  workerNotify.Origin,
		Payload: workerNotify.Payload,
	}
	worker.Notify.Execute(notifyBody)
	fmt.Printf("Notify sended to work %s", worker.Id)
}
