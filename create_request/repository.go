package create_request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
	"io"
)

type CreateRequestRepository struct {
	client *esdb.Client
}

func NewCreateRequestRepository(client *esdb.Client) *CreateRequestRepository {
	crr := CreateRequestRepository{client: client}

	return &crr
}

func (crr *CreateRequestRepository) Load(ctx context.Context, id uuid.UUID) *CreateRequest {
	stream, err := crr.client.ReadStream(ctx, id.String(), esdb.ReadStreamOptions{}, 10)

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	var events []Event
	for {
		event, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		// Doing something productive with the event
		switch event.OriginalEvent().EventType {
		case CreateRequestApproved{}.Name():
			domainEvent := CreateRequestApproved{}
			err = json.Unmarshal(event.OriginalEvent().Data, &domainEvent)
			events = append(events, domainEvent)
		case CreateRequestRejected{}.Name():
			domainEvent := CreateRequestRejected{}
			err = json.Unmarshal(event.OriginalEvent().Data, &domainEvent)
			events = append(events, domainEvent)
		case CreateRequestRequested{}.Name():
			domainEvent := CreateRequestRequested{}
			err = json.Unmarshal(event.OriginalEvent().Data, &domainEvent)
			events = append(events, domainEvent)
		}

		if err != nil {
			panic(err)
		}
	}

	return NewFromEvents(events)
}

func (crr *CreateRequestRepository) Save(ctx context.Context, cr *CreateRequest) error {
	records := make([]esdb.EventData, len(cr.Events()))
	var err error
	for i, ev := range cr.Events() {
		data, err := json.Marshal(ev)
		if err != nil {
			panic(err)
		}

		records[i] = esdb.EventData{
			ContentType: esdb.JsonContentType,
			EventType:   ev.Name(),
			EventID:     uuid.Must(uuid.NewV4()),
			Data:        data,
		}
	}

	if err != nil {
		return err
	}

	_, err = crr.client.AppendToStream(ctx, cr.ID().String(), esdb.AppendToStreamOptions{}, records...)
	if err != nil {
		return errors.New(fmt.Sprintf("A crapat call-ul Reason: %s", err.Error()))
	}

	return nil
}
