package notify_request

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
)

type NotifyRequest struct {
	Url     string            `json:"url"`
	Header  map[string]string `json:"header"`
	Timeout int64             `json:"timeout"`
}

type NotifyBody struct {
	Origin  string        `json:"origin"`
	Id      string        `json:"id"`
	Payload NotifyPayload `json:"payload"`
}

type NotifyPayload struct {
	Status   custom_types.Status        `json:"status"`
	Type     custom_types.ExecutionType `json:"type"`
	Response map[string]any             `json:"response"`
}

func (request NotifyRequest) RequestBody(body NotifyBody) (*bytes.Reader, error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonStr)
	return bodyReader, nil
}
