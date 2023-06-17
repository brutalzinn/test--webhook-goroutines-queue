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
	Body    any
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
