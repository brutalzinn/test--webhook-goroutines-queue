package main

import (
	"fmt"
	"net/http"
	"time"
)

type Priority uint8
type Status uint8

const (
	Low    Priority = 0
	Medium          = 1
	High            = 2
)

const (
	Created  Status = 0
	Pending         = 1
	Approved        = 2
	Rejected        = 3
)

type Job struct {
	Request  Request
	Response Response
	Priority Priority
	Status   Status
}

func (job *Job) Execute() {
	job.Status = Created
	request_body, err := job.Request.RequestBody()
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		job.Status = Rejected
		return
	}
	req, err := http.NewRequest(job.Request.Verb, job.Request.Url, request_body)
	job.Status = Pending
	if err != nil {
		fmt.Printf("client: could not create request: %s\n Reprocess this after", err)
		job.Status = Rejected
	}
	for key, value := range job.Request.Header {
		req.Header.Set(key, value)
	}
	client := http.Client{
		Timeout: time.Duration(job.Request.Timeout) * time.Second,
	}
	request, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n. Reprocess this after", err)
		job.Status = Rejected
	}
	if request.StatusCode == job.Response.Code {
		job.Status = Approved
	}
	job.Status = Rejected
}
