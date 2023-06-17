package webhook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	webhook_models "github.com/brutalzinn/test-webhook-goroutines-queue.git/webhook/models"
)

type Webhook struct {
	Request  webhook_models.WebhookRequest
	Response any
	Status
}

func (wh *Webhook) Execute() {
	wh.Status = Created
	request_body, err := wh.Request.RequestBody()
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		wh.Status = Rejected
		return
	}
	request, err := http.NewRequest(wh.Request.Verb, wh.Request.Url, request_body)
	wh.Status = Pending
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		wh.Status = Rejected
	}
	for key, value := range wh.Request.Header {
		request.Header.Set(key, value)
	}
	client := http.Client{
		Timeout: time.Duration(wh.Request.Timeout) * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n. Reprocess this after", err)
		wh.Status = Rejected
		return
	}
	if response.StatusCode >= 200 && response.StatusCode <= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		wh.Response = string(body)
		wh.Status = Approved
		return
	}
	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		body, _ := ioutil.ReadAll(response.Body)
		wh.Response = string(body)
		wh.Status = Reprocess
		return
	}
	wh.Status = Rejected
}
