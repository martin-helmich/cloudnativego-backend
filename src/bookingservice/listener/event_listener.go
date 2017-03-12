package listener

import (
	"bitbucket.org/minamartinteam/myevents/src/contracts"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/builder"
	"fmt"
	"log"
	"sync"
)

func ProcessEvents(wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := builder.NewEventListenerFromEnvironment()
	if err != nil {
		panic(err)
	}

	log.Println("listening or events")

	received, errors, err := listener.Listen("eventCreated")

	if err != nil {
		panic(err)
	}

	for {
		select {
		case evt := <-received:
			fmt.Printf("got event %T: %s\n", evt, evt)
			handleEvent(evt)
		case err = <-errors:
			fmt.Printf("got error while receiving event: %s\n", err)
		}
	}
}

func handleEvent(event msgqueue.Event) {
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
	case *contracts.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
	default:
		log.Printf("unknown event type: %T", e)
	}
}
