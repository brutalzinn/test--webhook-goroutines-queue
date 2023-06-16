package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Job struct {
	Notify   `json:"notify"`
	Request  `json:"request"`
	Response any
	Priority `json:"priority"`
	Status   `json:"status"`
}

func (job *Job) Execute() {
	job.Status = Created
	request_body, err := job.Request.RequestBody()
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		job.Status = Rejected
		return
	}
	request, err := http.NewRequest(job.Request.Verb, job.Request.Url, request_body)
	job.Status = Pending
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		job.Status = Rejected
	}
	for key, value := range job.Request.Header {
		request.Header.Set(key, value)
	}
	client := http.Client{
		Timeout: time.Duration(job.Request.Timeout) * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n. Reprocess this after", err)
		job.Status = Rejected
		return
	}
	if response.StatusCode >= 200 && response.StatusCode <= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		job.Response = string(body)
		job.Status = Approved
		return
	}
	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		body, _ := ioutil.ReadAll(response.Body)
		job.Response = string(body)
		job.Status = Reprocess
		return
	}
	job.Status = Rejected
}
