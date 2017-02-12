package amqp

import (
	"github.com/streadway/amqp"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"encoding/json"
	"fmt"
)

type AMQPEventEmitter struct {
	channel  *amqp.Channel
	exchange string
	events   chan *emittedEvent
}

type emittedEvent struct {
	event     msgqueue.Event
	errorChan chan error
}

func NewAMQPEventEmitter(conn *amqp.Connection, exchange string) (*AMQPEventEmitter, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not create AMQP channel: %s", err)
	}

	emitter := AMQPEventEmitter{
		channel: channel,
		exchange: exchange,
		events: make(chan *emittedEvent),
	}

	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for item := range emitter.events {
			emitter.emitItem(item)
		}
	}()

	return &emitter, nil
}

func (e *AMQPEventEmitter) emitItem(item *emittedEvent) {
	defer close(item.errorChan)

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

func (e *AMQPEventEmitter) Emit(event msgqueue.Event) error {
	errChan := make(chan error)
	item := &emittedEvent{
		event: event,
		errorChan: errChan,
	}

	e.events <- item
	return <-errChan
}