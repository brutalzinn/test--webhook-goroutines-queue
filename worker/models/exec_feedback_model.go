package worker

import (
	"time"

	"github.com/brutalzinn/test-webhook-goroutines-queue.git/custom_types"
)

type ExecFeedbackModel struct {
	Response  string
	Request   string
	ExecuteAt time.Time
	Status    custom_types.Status
}
