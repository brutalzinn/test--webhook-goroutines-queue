package notify_request

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type NotifyRequest struct {
	Url     string            `json:"url"`
	Header  map[string]string `json:"header"`
	Timeout int64             `json:"timeout"`
}

type NotifyBody struct {
	Origin  string
	Payload any
}

func (request *NotifyRequest) RequestBody(body NotifyBody) (*bytes.Reader, error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonStr)
	return bodyReader, nil
}
