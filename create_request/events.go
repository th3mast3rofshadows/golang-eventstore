package create_request

type Event interface {
	isEvent()
	Name() string
}

func (e CreateRequestApproved) isEvent()  {}
func (e CreateRequestRejected) isEvent()  {}
func (e CreateRequestRequested) isEvent() {}

type CreateRequestApproved struct {
	Status string `json:"status"`
}

func (e CreateRequestApproved) Name() string {
	return "CreateRequestApproved"
}

type CreateRequestRejected struct {
	Status string `json:"status"`
}

func (e CreateRequestRejected) Name() string {
	return "CreateRequestRejected"
}

type CreateRequestRequested struct {
	RequestId string `json:"request_id"`
	VendorId  int    `json:"vendor_id"`
	ExtId     int    `json:"ext_id"`
	Platform  string `json:"platform"`
}

func (e CreateRequestRequested) Name() string {
	return "CreateRequestRequested"
}
