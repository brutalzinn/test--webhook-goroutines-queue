package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Notify struct {
	Request NotifyRequest
	Status
}

type NotifyRequest struct {
	Url     string            `json:"url"`
	Header  map[string]string `json:"header"`
	Timeout int64             `json:"timeout"`
	Body    any
}

func (notify *Notify) Execute() {
	notify.Status = Created
	request_body, err := notify.Request.RequestBody()
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

func (request NotifyRequest) RequestBody() (*bytes.Reader, error) {
	jsonStr, err := json.Marshal(request.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonStr)
	return bodyReader, nil
}
