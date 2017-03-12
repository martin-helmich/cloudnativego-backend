package main

import (
	"bitbucket.org/minamartinteam/myevents/src/contracts"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	evtamqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/kafka"
	"log"
	"os"
	"time"
)

func main() {
	var emitter msgqueue.EventEmitter
	var err error

	log.Println("emitting example event")

	exampleEvent := &contracts.EventCreatedEvent{
		ID:    "asasd",
		Name:  "Wacken Open Air",
		Start: time.Now(),
		End:   time.Now().Add(3 * 24 * time.Hour),
	}

	if url := os.Getenv("AMQP_URL"); url != "" {
		log.Printf("connecting to AMQP broker at %s", url)

		emitter, err = evtamqp.NewAMQPEventEmitterFromEnvironment()
		if err != nil {
			panic(err)
		}
	} else if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		log.Printf("connecting to Kafka brokers at %s", brokers)

		emitter, err = kafka.NewKafkaEventEmitterFromEnvironment()
		if err != nil {
			panic(err)
		}
	} else {
		panic("Neither AMQP_URL nor KAFKA_BROKERS specified")
	}

	log.Println("sleeping 10 seconds")
	time.Sleep(10 * time.Second)

	log.Println("emitting example event")
	err = emitter.Emit(exampleEvent)
	if err != nil {
		panic(err)
	}
}
