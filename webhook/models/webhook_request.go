package webhook_models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type WebhookRequest struct {
	Url     string            `json:"url"`
	Verb    string            `json:"verb"`
	Timeout int64             `json:"timeout"`
	Header  map[string]string `json:"header"`
	Body    any               `json:"body"`
}

func (request WebhookRequest) RequestBody() (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(request.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonBytes)
	return bodyReader, nil
}

func (request WebhookRequest) RequestBodyString() string {
	jsonBytes, err := json.Marshal(request.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return ""
	}
	bodyString := string(jsonBytes)
	return bodyString
}
