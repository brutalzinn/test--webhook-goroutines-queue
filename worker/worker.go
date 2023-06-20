package worker

import (
	"fmt"
	"strings"
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/notify"
	notify_request "github.com/brutalzinn/test-webhook-goroutines-queue.git/notify/models"
	"github.com/gofrs/uuid"
)

type Worker struct {
	Id            string
	Options       *WorkerOptions
	ServiceType   custom_types.ServiceType
	ExecutionType custom_types.ExecutionType
	Exec          func() WorkerFeedbackModel
	Notify        *notify.Notify
}

type WorkerOptions struct {
	ExecuteAt time.Time
	Priority  custom_types.Priority
}

type WorkerCompleted struct {
	Id            string
	FeedbackModel WorkerFeedbackModel
	WorkerNotify  WorkerNotify
}

type WorkerNotify struct {
	Origin  string
	Payload notify_request.NotifyPayload
}
type WorkerFeedbackModel struct {
	Response  map[string]any
	Request   map[string]any
	ExecuteAt time.Time
	Status    custom_types.Status
}

func New(exec func() WorkerFeedbackModel, serviceType custom_types.ServiceType) Worker {
	id, _ := uuid.NewV4()
	return Worker{
		Id:            id.String(),
		Exec:          exec,
		ExecutionType: custom_types.Normal,
		ServiceType:   serviceType,
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
	id, err := workerLog.Insert()
	if err != nil {
		fmt.Printf("Wrong at worker insert %s at %s\n", id, time.Now())
	}
	fmt.Printf("Run worker %s at %s\n", id, time.Now())
	onCompletedEvent := WorkerCompleted{
		FeedbackModel: execFeedback,
		Id:            id,
	}
	worker.onComplete(onCompletedEvent)
}

func (worker *Worker) ExecuteShedule() {
	workerLog := WorkerLog{
		Worker:          worker,
		Status:          custom_types.Created,
		RequestPayload:  map[string]any{},
		ResponsePayload: map[string]any{},
	}
	id, err := workerLog.Insert()
	fmt.Printf("### Worker sheduler  %s inserted at %s executeAt:%s \n", id, time.Now(), worker.Options.ExecuteAt)
	if err != nil {
		fmt.Printf("Wrong at worker sheduler insert %s at %s\n", id, time.Now())
	}
	onCompletedEvent := WorkerCompleted{
		FeedbackModel: WorkerFeedbackModel{
			Status: custom_types.Created,
		},
		Id: id,
	}
	worker.onComplete(onCompletedEvent)
	time.AfterFunc(worker.Options.ExecuteAt.Sub(time.Now()), func() {
		feedbackModel := worker.Exec()
		workerLog := WorkerLog{
			Worker:          worker,
			Status:          feedbackModel.Status,
			RequestPayload:  feedbackModel.Request,
			ResponsePayload: feedbackModel.Response,
		}
		err = workerLog.Update(id)
		fmt.Printf("### Worker sheduler %s update at %s executeAt:%s \n", id, time.Now(), worker.Options.ExecuteAt)
		if err != nil {
			fmt.Printf("Wrong at worker sheduler update %s at %s\n", id, time.Now())
		}
		onCompletedEvent := WorkerCompleted{
			FeedbackModel: feedbackModel,
			Id:            id,
		}
		worker.onComplete(onCompletedEvent)
	})
}
func (worker *Worker) onComplete(workerCompleted WorkerCompleted) {
	worker.sendWorkerNotify(workerCompleted)
	fmt.Printf("On complete to work %s", workerCompleted.Id)
}
func (worker *Worker) sendWorkerNotify(workerCompleted WorkerCompleted) {
	origin := fmt.Sprintf("%s_%s", worker.ServiceType.String(), worker.ExecutionType.String())
	notifyBody := notify_request.NotifyBody{
		Id:     workerCompleted.Id,
		Origin: strings.ToUpper(origin),
		Payload: notify_request.NotifyPayload{
			Type:     worker.ExecutionType,
			Status:   workerCompleted.FeedbackModel.Status,
			Response: workerCompleted.FeedbackModel.Response,
		},
	}
	worker.Notify.Execute(notifyBody)
	fmt.Printf("Notify sended to work %s", workerCompleted.Id)
}
