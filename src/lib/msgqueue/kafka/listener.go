package kafka

import (
	"github.com/Shopify/sarama"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"log"
	"reflect"
	"fmt"
)

type kafkaEventListener struct {
	topic      string
	consumer   sarama.Consumer
	partitions []int32
	mapper     *msgqueue.EventMapper
}

func NewKafkaEventListener(client sarama.Client, topic string, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		topic: topic,
		consumer: consumer,
		partitions: partitions,
		mapper: msgqueue.NewEventMapper(),
	}

	return listener, nil
}

func (k *kafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error

	results := make(chan msgqueue.Event)
	errors := make(chan error)

	for _, topic := range events {
		partitions := k.partitions
		if len(partitions) == 0 {
			partitions, err = k.consumer.Partitions(topic)
			if err != nil {
				return nil, nil, err
			}
		}

		log.Printf("topic %s has partitions: %v", topic, partitions)

		for _, partition := range partitions {
			log.Printf("consuming partition %s:%d", topic, partition)

			pConsumer, err := k.consumer.ConsumePartition(topic, partition, 0)
			if err != nil {
				return nil, nil, err
			}

			go func() {
				for msg := range pConsumer.Messages() {
					log.Printf("received message %v", msg)

					event, err := k.mapper.MapEvent(msg.Topic, msg.Value)
					if err != nil {
						errors <- fmt.Errorf("could not map message: %v", err)
					}

					results <- event
				}
			}()

			go func() {
				for err := range pConsumer.Errors() {
					errors <- err
				}
			}()
		}
	}

	return results, errors, nil
}


// Map registers event names that should be mapped to certain types.
func (l *kafkaEventListener) Map(typ reflect.Type) {
	l.mapper.RegisterMapping(typ)
}