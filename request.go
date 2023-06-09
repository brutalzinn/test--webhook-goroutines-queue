package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Request struct {
	Url     string
	Verb    string
	Timeout int64
	Header  map[string]string
	Body    any
}

func (request Request) RequestBody() (*bytes.Reader, error) {
	jsonStr, err := json.Marshal(request.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(jsonStr)
	return bodyReader, nil
}
