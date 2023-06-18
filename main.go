package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/models"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/notify"
	notify_request "github.com/brutalzinn/test-webhook-goroutines-queue.git/notify/models"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/queue"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/webhook"
	webhook_models "github.com/brutalzinn/test-webhook-goroutines-queue.git/webhook/models"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/worker"
)

/*
[BLOG.ROBERTOCPAES.DEV - HYPERFOCUS - PERSONAL DATA COMMENTARY - IGNORE]
	lets suppose in some place of world i need to implement a queue to execute web hooks.
	But i need to do this without dependencies because my i dont want much more thant a simple socket opened to receive and do the works for me.
	Sometimes i do some mistakes with overthinking about the design and patterns concepts before focus the main problem. Be adaptative its not wrong.. but i need be responsible to changes more than adaptative to problems.
	I not talking about software development. I talking about how i execute my normal activities.
	Attack the main problem without create more complexy than the problem really is.
	Now we need to simplifies this case.
	lets suppose this project is a really work and very important.
	lets suppose the ambient that i was writting this comementary has many counter situations to my sensitive perceptions.
	Now the drain of energy intercepts all my body requests and i am trying to not break all the other functional activies.
	The problem  is the constantly drain of energy and if we have much sensitive tasks envolved too.
	Like when we are developing and we really ok to do that, but something happened on the ambience and now we cant reponds everything at same time.
	But.. now we have the consience that computacional abstraction and the life are soo divergent. I cant assume that my body will reponds like a Circuit Break and no one will see that i am not really responsive.
	I am really OK and this is normal. But can dispatch some fractual errors on common dialogues, singularity of personal interprations of abstract language and language figures interpretations are broken.
	ITS TIME TO UNDESTAND HOW TO HANDLE WITH THIS KIND OF SITUATION.

	FOR MY ORGANIZATIONAL THERAPEUTIC: THIS IS NOT A HYPERFOCUS PROBLEM. ITS CAN BE LIFE HANDLE EXCEPTION PROBLEM. HOW CAN I SOLVE A PROBLEM WITHOUT PUT IN A QUEUE FIRST and order by priority?
	IF SOMETHING HAPPENS, I STILL NEED REPONDS THE OTHER PROCESSES. THIS IS NOT A PERFORMACE ISSUE. ITS A LOGICAL MIND PROCEDURE that needs be viewed at therapy ambient.
[BLOG.ROBERTOCPAES.DEV - HYPERFOCUS - PERSONAL DATA COMMENTARY - IGNORE]
*/

func main() {
	normalQueue, shedulerQueue := queue.Queue{}, queue.Queue{}
	queue_channel := make(chan queue.Queue)
	// goroutine to perform normal queue operation
	go func() {
		for range queue_channel {
			for !normalQueue.IsEmpty() {
				item := normalQueue.Dequeue()
				if item.Options.ExecuteAt.After(time.Now()) {
					fmt.Printf("Jumping %s and sheduler to %s\n ", item.Id, item.Options.ExecuteAt)
					shedulerQueue.Enqueue(item)
					continue
				}
				item.Execute()
				fmt.Printf("execute the worker %s\n", item.Id)
			}
			fmt.Printf("Process all normalQueue %s lenght %b\n", time.Now(), len(normalQueue.Workers))
		}
	}()
	// goroutine to perform sheduler queue operation
	go func() {
		for range time.Tick(time.Minute * 1) {
			for !shedulerQueue.IsEmpty() {
				item := shedulerQueue.Dequeue()
				item.ExecutionType = custom_types.Sheduler
				item.ExecuteShedule()
				fmt.Printf("execute the worker sheduled %s\n", item.Id)
			}
			fmt.Printf("Process all shedulerQueue at %s lenght %b\n", time.Now(), len(shedulerQueue.Workers))
		}
	}()

	http.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		var queueRequest models.QueueRequest
		err := json.NewDecoder(r.Body).Decode(&queueRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		wh := webhook.Webhook{
			Request: webhook_models.WebhookRequest{
				Url:     queueRequest.Request.Url,
				Verb:    queueRequest.Request.Verb,
				Timeout: queueRequest.Request.Timeout,
				Header:  queueRequest.Request.Header,
				Body:    queueRequest.Request.Body,
			},
		}
		options := worker.WorkerOptions{
			ExecuteAt: queueRequest.ExecuteAt,
			Priority:  custom_types.Low,
		}
		notify := notify.Notify{
			Request: notify_request.NotifyRequest{
				Url:     queueRequest.Notify.Url,
				Header:  queueRequest.Notify.Header,
				Timeout: queueRequest.Notify.Timeout,
			},
		}
		worker := worker.New(wh.Execute, custom_types.Webhook)
		worker = worker.WithOptions(&options)
		worker = worker.WithNotify(&notify)
		normalQueue.Enqueue(worker)
		queue_channel <- normalQueue
		w.WriteHeader(200)
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received Body: %+v", r.Body)
		fmt.Fprintf(w, "Body: %+v", r.Body)
	})
	http.ListenAndServe(":9000", nil)
}
