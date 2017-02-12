package main

import (
	"time"
	evtamqp "bitbucket.org/minamartinteam/myevents/src/lib/msgqueue/amqp"
	"os"
	"bitbucket.org/minamartinteam/myevents/src/contracts/events"
	"bitbucket.org/minamartinteam/myevents/src/lib/helper/amqp"
	"log"
)

func main() {
	conn := <- amqp.RetryConnect(os.Getenv("AMQP_URL"), 5 * time.Second)

	time.Sleep(10 * time.Second)

	log.Println("emitting example event")

	emitter, _ := evtamqp.NewAMQPEventEmitter(conn, "example")
	emitter.Emit(&events.EventCreatedEvent{
		ID: "asasd",
		Name: "Wacken Open Air",
		Start: time.Now(),
		End: time.Now().Add(3 * 24 * time.Hour),
	})
}
