package custom_types

type ServiceType uint8
type Priority uint8
type ExecutionType uint8
type Status uint8

const (
	Normal   ExecutionType = 1
	Sheduler               = 2
)

const (
	Low    Priority = 1
	Medium          = 2
	High            = 3
)

const (
	Webhook ServiceType = 1
	///i like when i suppose that i can have other integrations types after. This is a overkill for everything i need now.
)

const (
	Created  Status = 1
	Pending         = 2
	Approved        = 3
	Rejected        = 4
	Error           = 5
)
