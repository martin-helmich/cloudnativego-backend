package main

import (
	"bitbucket.org/minamartinteam/myevents/src/contracts/events"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	evtamqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/kafka"
	"fmt"
	"log"
	"os"
	"reflect"
)

func main() {
	var listener msgqueue.EventListener
	var err error

	if url := os.Getenv("AMQP_URL"); url != "" {
		log.Printf("connecting to AMQP broker at %s", url)

		listener, err = evtamqp.NewAMQPEventListenerFromEnvironment()
		if err != nil {
			panic(err)
		}
	} else if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		log.Printf("connecting to Kafka brokers at %s", brokers)

		listener, err = kafka.NewKafkaEventListenerFromEnvironment()
		if err != nil {
			panic(err)
		}
	} else {
		panic("Neither AMQP_URL nor KAFKA_BROKERS specified")
	}

	log.Println("listening or events")

	listener.Map(reflect.TypeOf(events.EventCreatedEvent{}))
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
	case *events.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
	case *events.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
	default:
		log.Printf("unknown event type: %T", e)
	}
}
