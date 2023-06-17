package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/models"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/queue"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/webhook"
	webhook_models "github.com/brutalzinn/test-webhook-goroutines-queue.git/webhook/models"
	"github.com/brutalzinn/test-webhook-goroutines-queue.git/worker"
)

// lets suppose in some place of world i need to implement a queue to execute web hooks to all the requests failed.
// But i need to do this without dependencies because my i dont want much more thant a simple socket opened to receive do the works for me.
// So we use the gouroutines do this with channels. With this we attack my asks for a queue or a stack principle?!
// Sometimes i do some mistakes with overthinking about the design and patterns concepts before focus the main problem. Now we need to simplifies this cases.
// ITS TIME TO CHANGE THIS.
// FOR MY ORGANIZATIONAL THERAPEUTIC: THIS IS NOT A HYPERFOCUS PROBLEM. ITS A LIFE PROBLEM. HOW CAN I SOLVE A PROBLEM WITHOUT PUT IN A QUEUE FIRST and order by priority? IF SOMETHING HAPPENS,
// I STILL NEED REPONDS THE OTHER PROCESSES. THIS IS NOT A PERFORMACE ISSUE. ITS A LOGICAL MIND PORPOUSE.

func main() {

	processQueue, shedulerQueue := queue.Queue{}, queue.Queue{}

	go func() {
		for range time.Tick(time.Second * 1) {
			for !processQueue.IsEmpty() {
				item := processQueue.Dequeue()
				if item.Options.ExecuteAt.After(time.Now()) {
					fmt.Printf("Jumping %s and sheduler at %s\n ", item.Id, item.Options.ExecuteAt)
					shedulerQueue.Enqueue(item)
					continue
				}
				item.Execute()
				fmt.Printf("execute the worker %s\n", item.Id)
			}
			fmt.Printf("Process all processQueue %s lenght %b\n", time.Now(), len(processQueue.Workers))
		}
	}()

	go func() {
		for range time.Tick(time.Minute * 1) {
			for !shedulerQueue.IsEmpty() {
				item := shedulerQueue.Dequeue()
				item.ExecuteShedule()
				fmt.Printf("execute the worker sheduled %s\n", item.Id)
			}
			fmt.Printf("Process all sheduler queue %s lenght %b\n", time.Now(), len(shedulerQueue.Workers))
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
			Priority:  queueRequest.Priority,
		}
		worker := worker.New(wh.Execute).WithOptions(&options)
		processQueue.Enqueue(worker)
		w.WriteHeader(200)
	})

	// http.HandleFunc("/requeue", func(w http.ResponseWriter, r *http.Request) {
	// 	id := r.URL.Query().Get("id")
	// 	var job Job
	// 	err := json.NewDecoder(r.Body).Decode(&job)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	go worker(job_channel)
	// 	job_channel <- job
	// 	fmt.Fprintf(w, "Body: %+v", r.Body)
	// })
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received Body: %+v", r.Body)
		fmt.Fprintf(w, "Body: %+v", r.Body)
	})

	http.ListenAndServe(":9000", nil)

}
