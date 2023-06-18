package notify

import (
	"fmt"
	"net/http"
	"time"

	notify_request "github.com/brutalzinn/test-webhook-goroutines-queue.git/notify/models"
)

type Notify struct {
	Request notify_request.NotifyRequest
	Status
}

func (notify *Notify) Execute(body notify_request.NotifyBody) {
	notify.Status = Created
	request_body, err := notify.Request.RequestBody(body)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		notify.Status = Rejected
		return
	}
	req, err := http.NewRequest(http.MethodPost, notify.Request.Url, request_body)
	notify.Status = Pending
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		notify.Status = Rejected
	}
	for key, value := range notify.Request.Header {
		req.Header.Set(key, value)
	}
	client := http.Client{
		Timeout: time.Duration(notify.Request.Timeout) * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http response: %s\n. Reprocess this after", err)
		notify.Status = Rejected
		return
	}
	if response.StatusCode >= 200 && response.StatusCode <= 300 {
		notify.Status = Approved
		return
	}
	notify.Status = Rejected
}
