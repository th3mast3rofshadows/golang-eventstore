package main

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
	"golang-eventstore/create_request"
	"os"
)

func main() {
	db := setupEventStore()

	defer func(eventDb *esdb.Client) {
		err := eventDb.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	crr := create_request.NewCreateRequestRepository(db)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "publish":
			fmt.Println("publishing")
			publish(crr)
		case "read":
			fmt.Println("reading")
			read(crr, os.Args[2])
		default:
		}
	}
}

func setupEventStore() *esdb.Client {
	settings, err := esdb.ParseConnectionString("esdb://localhost:2113?tls=false")

	if err != nil {
		panic(err)
	}

	db, err := esdb.NewClient(settings)

	if err != nil {
		panic(err)
	}

	return db
}

func publish(crr *create_request.CreateRequestRepository) {
	cr := create_request.Request("some-requestId", 123, 234, "emag-ro")
	err := cr.Approve()
	if err != nil {
		panic(err)
	}

	err = crr.Save(context.Background(), cr)
	if err != nil {
		panic(err)
	}
}

func read(crr *create_request.CreateRequestRepository, id string) {
	fmt.Printf("%v", crr.Load(context.Background(), uuid.Must(uuid.FromString(id))))
}
