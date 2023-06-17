package webhook

type Priority uint8
type Status uint8

const (
	Low    Priority = 0
	Medium          = 1
	High            = 2
)

const (
	Created   Status = 0
	Pending          = 1
	Approved         = 2
	Rejected         = 3
	Reprocess        = 4
)
