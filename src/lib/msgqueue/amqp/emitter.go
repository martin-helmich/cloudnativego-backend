package amqp

import (
	"github.com/streadway/amqp"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"encoding/json"
	"fmt"
)

type amqpEventEmitter struct {
	channel  *amqp.Channel
	exchange string
	events   chan *emittedEvent
}

type emittedEvent struct {
	event     msgqueue.Event
	errorChan chan error
}

// NewAMQPEventEmitter creates a new event emitter.
// It will need an AMQP connection passed as parameter and use this connection
// to create its own channel (note: AMQP channels are not thread-safe, so just
// accepting the connection as a parameter and then creating our own private
// channel is the safest way to ensure this).
func NewAMQPEventEmitter(conn *amqp.Connection, exchange string) (msgqueue.EventEmitter, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not create AMQP channel: %s", err)
	}

	emitter := amqpEventEmitter{
		channel: channel,
		exchange: exchange,
		events: make(chan *emittedEvent),
	}

	// Normally, all(many) of these options should be configurable.
	// For our example, it'll probably do.
	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// AMQP channels are not thread-safe. Since this struct's `Emit` function
	// could be called from anywhere, we have `Emit` simply write into a channel
	// and have a single go-routine read from that channel and do the actual
	// publishing.
	go func() {
		for item := range emitter.events {
			emitter.emitItem(item)
		}
	}()

	return &emitter, nil
}

func (e *amqpEventEmitter) emitItem(item *emittedEvent) {
	defer close(item.errorChan)

	// TODO: Alternatives to JSON? Msgpack or Protobuf, maybe?
	jsonBody, err := json.Marshal(item.event)
	if err != nil {
		item.errorChan <- fmt.Errorf("could not JSON-serialize event: %s", err)
		return
	}

	pub := amqp.Publishing{
		Headers: amqp.Table{
			"x-event-name": item.event.EventName(),
		},
		ContentType: "application/json",
		Body: jsonBody,
	}

	err = e.channel.Publish(e.exchange, item.event.EventName(), false, false, pub)
	if err != nil {
		item.errorChan <- fmt.Errorf("could not publish message to AMQP: %s", err)
	}
}

func (e *amqpEventEmitter) Emit(event msgqueue.Event) error {
	errChan := make(chan error)
	item := &emittedEvent{
		event: event,
		errorChan: errChan,
	}

	e.events <- item
	return <-errChan
}