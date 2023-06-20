package webhook_models

import (
	"encoding/json"
	"fmt"
)

type WebhookResponse struct {
	StatusCode int            `json:"status_code"`
	Body       map[string]any `json:"body"`
}

func (response WebhookResponse) ResponseBodyMap() map[string]any {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil
	}
	var keysPair map[string]any
	_ = json.Unmarshal(jsonBytes, &keysPair)
	return keysPair
}
