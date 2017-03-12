package amqp

import (
	amqphelper "bitbucket.org/minamartinteam/myevents/src/lib/helper/amqp"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"time"
)

const eventNameHeader = "x-event-name"

type amqpEventListener struct {
	channel  *amqp.Channel
	exchange string
	queue    string
	mapper   msgqueue.EventMapper
}

// NewAMQPEventListenerFromEnvironment will create a new event listener from
// the configured environment variables. Important variables are:
//
//   - AMQP_URL; the URL of the AMQP broker to connect to
//   - AMQP_EXCHANGE; the name of the exchange to bind to
//   - AMQP_QUEUE; the name of the queue to bind and subscribe
//
// For missing environment variables, this function will assume sane defaults.
func NewAMQPEventListenerFromEnvironment() (msgqueue.EventListener, error) {
	var url string
	var exchange string
	var queue string

	if url = os.Getenv("AMQP_URL"); url == "" {
		url = "amqp://localhost:5672"
	}

	if exchange = os.Getenv("AMQP_EXCHANGE"); exchange == "" {
		exchange = "example"
	}

	if queue = os.Getenv("AMQP_QUEUE"); queue == "" {
		queue = "example"
	}

	conn := <-amqphelper.RetryConnect(url, 5*time.Second)
	return NewAMQPEventListener(conn, exchange, queue)
}

// NewAMQPEventListener creates a new event listener.
// It will need an AMQP connection passed as parameter and use this connection
// to create its own channel (note: AMQP channels are not thread-safe, so just
// accepting the connection as a parameter and then creating our own private
// channel is the safest way to ensure this).
func NewAMQPEventListener(conn *amqp.Connection, exchange string, queue string) (msgqueue.EventListener, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("could not create AMQP channel: %s", err)
	}

	listener := amqpEventListener{
		channel:  channel,
		exchange: exchange,
		queue:    queue,
		mapper:   msgqueue.NewEventMapper(),
	}

	err = channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

// Listen configures the event listener to listen for a set of events that are
// specified by name as parameter.
// This method will return two channels: One will contain successfully decoded
// events, the other will contain errors for messages that could not be
// successfully decoded.
func (l *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
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
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("header %s did not contain string", eventNameHeader)
				msg.Nack(false, false)
				continue
			}

			event, err := l.mapper.MapEvent(eventName, msg.Body)
			if err != nil {
				errors <- fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
				msg.Nack(false, false)
				continue
			}

			events <- event
			msg.Ack(false)
		}
	}()

	return events, errors, nil
}

func (l *amqpEventListener) Mapper() msgqueue.EventMapper {
	return l.mapper
}