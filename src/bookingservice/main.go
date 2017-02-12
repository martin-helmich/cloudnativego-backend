package main

import (
	"time"
	evtamqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"bitbucket.org/minamartinteam/myevents/src/lib/helper/amqp"
	"bitbucket.org/minamartinteam/myevents/src/contracts/events"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"os"
	"reflect"
	"fmt"
	"log"
)

func main() {
	conn := <- amqp.RetryConnect(os.Getenv("AMQP_URL"), 5 * time.Second)

	listener, err := evtamqp.NewAMQPEventListener(conn, "example", "queue")
	if err != nil {
		panic(err)
	}

	listener.Map("eventCreated", reflect.TypeOf(events.EventCreatedEvent{}))
	received, errors, err := listener.Listen("eventCreated")
	for {
		select {
		case evt := <-received:
			fmt.Printf("got event %T: %s\n", evt, evt)
			handleEvent(evt)
		case err = <- errors:
			fmt.Printf("got error while receiving event: %s\n", err)
		}
	}
}

func handleEvent(event msgqueue.Event) {
	switch e := event.(type) {
	case *events.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
	case *events.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
	default:
		log.Printf("unknown event type: %T", e)
	}
}