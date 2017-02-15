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
	"strings"
	"github.com/Shopify/sarama"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/kafka"
)

func main() {
	var listener msgqueue.EventListener

	if url := os.Getenv("AMQP_URL"); url != "" {
		log.Printf("connecting to AMQP broker at %s", url)

		conn := <- amqp.RetryConnect(url, 5 * time.Second)
		listener, _ = evtamqp.NewAMQPEventListener(conn, "example", "queue")
	} else if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		log.Printf("connecting to Kafka brokers at %s", brokers)

		brokers := strings.Split(brokers, ",")
		client, err := sarama.NewClient(brokers, sarama.NewConfig())
		if err != nil {
			panic(err)
		}
		listener, err = kafka.NewKafkaEventListener(client, "", []int32{})
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