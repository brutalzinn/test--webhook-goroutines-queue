package models

import "time"

type NotifyRequest struct {
	Url     string            `json:"url"`
	Header  map[string]string `json:"header"`
	Timeout int64             `json:"timeout"`
	Body    any
}

type QueueWebhookRequest struct {
	Url     string            `json:"url"`
	Verb    string            `json:"verb"`
	Timeout int64             `json:"timeout"`
	Header  map[string]string `json:"header"`
	Body    any               `json:"body"`
}

type QueueRequest struct {
	ExecuteAt time.Time           `json:"execute_at"`
	Priority  int                 `json:"priority"`
	Request   QueueWebhookRequest `json:"request"`
	Response  any                 `json:"response"`
	Notify    NotifyRequest       `json:"notify"`
	Status    int                 `json:"status"`
}
