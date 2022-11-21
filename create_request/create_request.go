package create_request

import (
	"errors"
	"github.com/gofrs/uuid"
)

type CreateRequest struct {
	id        uuid.UUID
	requestId string
	vendorId  int
	extId     int
	platform  string
	status    string

	changes []Event
	version int
}

const Approved = "approved"
const Rejected = "rejected"
const New = "new"

func NewFromEvents(events []Event) *CreateRequest {
	cr := &CreateRequest{}

	for _, event := range events {
		cr.On(event, false)
	}

	return cr
}

func (cr *CreateRequest) RequestId() string {
	return cr.requestId
}

func (cr *CreateRequest) VendorId() int {
	return cr.vendorId
}

func (cr *CreateRequest) ExtId() int {
	return cr.extId
}

func (cr *CreateRequest) Platform() string {
	return cr.platform
}

func (cr *CreateRequest) Status() string {
	return cr.status
}

func (cr *CreateRequest) Approved() bool {
	return cr.status == Approved
}

func (cr *CreateRequest) Rejected() bool {
	return cr.status == Rejected
}

func (cr *CreateRequest) InManualValidation() bool {
	return cr.status == New
}

func (cr *CreateRequest) ID() uuid.UUID {
	return cr.id
}

func Request(requestId string, vendorId int, extId int, platform string) *CreateRequest {
	cr := &CreateRequest{}

	cr.raise(CreateRequestRequested{
		RequestId: requestId,
		VendorId:  vendorId,
		ExtId:     extId,
		Platform:  platform,
	})

	return cr
}

func (cr *CreateRequest) Approve() error {
	if !cr.InManualValidation() {
		return errors.New("cannot approve Create Request as it is not in Manual Validation")
	}

	cr.raise(CreateRequestApproved{
		Status: Approved,
	})

	return nil
}

func (cr *CreateRequest) Reject() error {
	if !cr.InManualValidation() {
		return errors.New("cannot reject Create Request as it is not in Manual Validation")
	}

	cr.raise(CreateRequestApproved{
		Status: Rejected,
	})

	return nil
}

func (cr *CreateRequest) On(event Event, new bool) {
	switch e := event.(type) {
	case CreateRequestApproved:
		cr.OnCreateRequestApproved(e)
	case CreateRequestRequested:
		cr.OnCreateRequestRequested(e)
	case CreateRequestRejected:
		cr.OnCreateRequestRejected(e)
	}

	if !new {
		cr.version++
	}
}

func (cr *CreateRequest) OnCreateRequestApproved(event CreateRequestApproved) {
	cr.status = event.Status
}

func (cr *CreateRequest) OnCreateRequestRejected(event CreateRequestRejected) {
	cr.status = event.Status
}

func (cr *CreateRequest) OnCreateRequestRequested(event CreateRequestRequested) {
	cr.id = uuid.Must(uuid.NewV4())
	cr.requestId = event.RequestId
	cr.vendorId = event.VendorId
	cr.extId = event.ExtId
	cr.platform = event.Platform
	cr.status = New
}

// Events returns the uncommitted events from the patient aggregate.
func (cr *CreateRequest) Events() []Event {
	return cr.changes
}

// Version returns the last version of the aggregate before changes.
func (cr *CreateRequest) Version() int {
	return cr.version
}

func (cr *CreateRequest) raise(event Event) {
	cr.changes = append(cr.changes, event)
	cr.On(event, true)
}
