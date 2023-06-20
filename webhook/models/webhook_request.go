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

func (request WebhookRequest) RequestBodyMap() map[string]any {
	jsonBytes, err := json.Marshal(request.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil
	}
	var keysPair map[string]any
	_ = json.Unmarshal(jsonBytes, &keysPair)
	return keysPair
}
