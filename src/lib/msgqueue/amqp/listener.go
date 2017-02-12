package amqp

import (
	"github.com/streadway/amqp"
	"fmt"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"reflect"
	"encoding/json"
)

const eventNameHeader = "x-event-name";

type AMQPEventListener struct {
	channel  *amqp.Channel
	exchange string
	queue    string
	typeMap  map[string]reflect.Type
}

func NewAMQPEventListener(conn *amqp.Connection, exchange string, queue string) (*AMQPEventListener, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not create AMQP channel: %s", err)
	}

	listener := AMQPEventListener{
		channel: channel,
		exchange: exchange,
		queue: queue,
		typeMap: make(map[string]reflect.Type),
	}

	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

func (l *AMQPEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	_, err := l.channel.QueueDeclare(l.queue, true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not declare queue %s: %s", l.queue, err)
	}

	// Create binding between queue and exchange for each listened event type
	for _, event := range eventNames {
		if err := l.channel.QueueBind(l.queue, event, l.exchange, false, nil); err != nil {
			return nil, nil, fmt.Errorf("could not bind event %s to queue %s: %s", event, l.queue, err)
		}
	}

	msgs, err := l.channel.Consume(l.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not consume queue: %s", err)
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers[eventNameHeader]
			if !ok {
				errors <- fmt.Errorf("message did not contain %s header", eventNameHeader)
				msg.Nack(false, false)
				return
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("header %s did not contain string", eventNameHeader)
				msg.Nack(false, false)
				return
			}

			typ, ok := l.typeMap[eventName]
			if !ok {
				errors <- fmt.Errorf("no mapping configured for event %s", eventName)
				msg.Nack(false, false)
				return
			}

			instance := reflect.New(typ)
			iface := instance.Interface()

			event, ok := iface.(msgqueue.Event)
			if !ok {
				errors <- fmt.Errorf("type %T does not implement the Event interface", iface)
				msg.Nack(false, false)
				return
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errors <- fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
				msg.Nack(false, false)
				return
			}

			events <- event
			msg.Ack(false)
		}
	}()

	return events, errors, nil
}

func (l *AMQPEventListener) Map(eventName string, typ reflect.Type) {
	l.typeMap[eventName] = typ
}