package webhook

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
	webhook_models "github.com/brutalzinn/test-webhook-goroutines-queue.git/webhook/models"
	worker "github.com/brutalzinn/test-webhook-goroutines-queue.git/worker/models"
)

type Webhook struct {
	Request  webhook_models.WebhookRequest
	Status   custom_types.Status
	Response string
}

func (wh *Webhook) Execute() worker.FeedbackModel {
	wh.Status = custom_types.Created
	request_body, err := wh.Request.RequestBody()
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		wh.Status = custom_types.Rejected
		return wh.createFeedbackModel()
	}
	request, err := http.NewRequest(wh.Request.Verb, wh.Request.Url, request_body)
	wh.Status = custom_types.Pending
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		wh.Status = custom_types.Rejected
		return wh.createFeedbackModel()
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
		wh.Status = custom_types.Rejected
		return wh.createFeedbackModel()
	}
	if response.StatusCode >= 200 && response.StatusCode <= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		wh.Response = string(body)
		wh.Status = custom_types.Approved
		return wh.createFeedbackModel()
	}
	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		body, _ := ioutil.ReadAll(response.Body)
		wh.Response = string(body)
		wh.Status = custom_types.Error
		return wh.createFeedbackModel()
	}
	wh.Status = custom_types.Rejected
	return wh.createFeedbackModel()
}

func (wh *Webhook) createFeedbackModel() worker.FeedbackModel {
	execFeedbackModel := worker.FeedbackModel{
		ExecuteAt: time.Now(),
		Response:  wh.Response,
		Request:   wh.Request.RequestBodyString(),
		Status:    wh.Status,
	}
	return execFeedbackModel
}
