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
Lets suppose i dont need to contains me anymore. Goodbye data tracking of my github profile.
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
					shedulerQueue.Enqueue(item)
					continue
				}
				item.Execute()
			}
		}
	}()
	// goroutine to perform sheduler queue operation
	go func() {
		for range time.Tick(time.Minute * 1) {
			for !shedulerQueue.IsEmpty() {
				item := shedulerQueue.Dequeue()
				item.ExecutionType = custom_types.Sheduler
				item.ExecuteShedule()
			}
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
