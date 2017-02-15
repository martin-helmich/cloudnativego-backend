package kafka

import (
	"github.com/Shopify/sarama"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"encoding/json"
	"log"
)

type kafkaEventEmitter struct {
	producer sarama.AsyncProducer
}

func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := kafkaEventEmitter{
		producer: producer,
	}

	return &emitter, nil
}

func (k *kafkaEventEmitter) Emit(evt msgqueue.Event) error {
	jsonBody, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: evt.EventName(),
		Value: sarama.ByteEncoder(jsonBody),
	}

	k.producer.Input() <- msg
	log.Printf("published message with topic %s: %v", evt.EventName(), jsonBody)

	go func() {
		for err := range k.producer.Errors() {
			log.Printf("error on emitter: %s", err)
		}
	}()

	success := <- k.producer.Successes()
	log.Printf("message successfully published: %v", success)
	return nil

	/*
	select {
	//case k.producer.Input() <- msg:
	//	log.Printf("published message with topic %s: %v", evt.EventName(), jsonBody)
	case success := <- k.producer.Successes():
		log.Printf("message successfully published: %v", success)
		return nil
	case err := <- k.producer.Errors():
		return err
	}

	return nil
	*/
}