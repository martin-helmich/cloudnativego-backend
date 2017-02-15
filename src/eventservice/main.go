package main

import (
	"time"
	evtamqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"os"
	"bitbucket.org/minamartinteam/myevents/src/contracts/events"
	"bitbucket.org/minamartinteam/myevents/src/lib/helper/amqp"
	"log"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"github.com/Shopify/sarama"
	"strings"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/kafka"
)

func main() {
	var emitter msgqueue.EventEmitter

	log.Println("emitting example event")

	exampleEvent := &events.EventCreatedEvent{
		ID: "asasd",
		Name: "Wacken Open Air",
		Start: time.Now(),
		End: time.Now().Add(3 * 24 * time.Hour),
	}

	if url := os.Getenv("AMQP_URL"); url != "" {
		log.Printf("connecting to AMQP broker at %s", url)

		conn := <- amqp.RetryConnect(url, 5 * time.Second)
		emitter, _ = evtamqp.NewAMQPEventEmitter(conn, "example")
	} else if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		log.Printf("connecting to Kafka brokers at %s", brokers)

		brokers := strings.Split(brokers, ",")
		client, err := sarama.NewClient(brokers, sarama.NewConfig())
		if err != nil {
			panic(err)
		}

		emitter, err = kafka.NewKafkaEventEmitter(client)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Neither AMQP_URL nor KAFKA_BROKERS specified")
	}

	log.Println("sleeping 10 seconds")
	time.Sleep(10 * time.Second)

	log.Println("emitting example event")
	err := emitter.Emit(exampleEvent)
	if err != nil {
		panic(err)
	}
}
